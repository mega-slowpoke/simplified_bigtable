package tablet

import (
	"bytes"
	"context"
	"encoding/gob"
	proto "final/proto/external-api"
	ipb "final/proto/internal-api"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ValueWithKeyAndTimestamps struct {
	Timestamp int64
	Value     string
	Key       string
}

func (s *TabletServiceServer) CreateTable(ctx context.Context, req *ipb.CreateTableInternalRequest) (*ipb.CreateTableInternalResponse, error) {
	tableName := req.TableName
	_, exist := s.Tables[tableName]
	if exist {
		return nil, status.Error(codes.AlreadyExists, "table already exists")
	}
	// : is not allowed in the file path
	dbPath := GetFilePath(s.TabletAddress, tableName)
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return &ipb.CreateTableInternalResponse{
			Success: false,
		}, nil
	}

	// update table Map
	s.Tables[tableName] = db

	// Initialize TablesRows for the new table
	s.TablesRows[tableName] = make(map[string]struct{})

	// update table columns info
	s.TablesColumns[tableName] = make(map[string][]string)
	columnFamilyMap := s.TablesColumns[tableName]
	for _, columnFamily := range req.ColumnFamilies {
		columnFamilyMap[columnFamily.FamilyName] = columnFamily.Columns
	}

	// persist metadata row so they can be recovered when the server crashes
	err = WriteMetaDataToPersistent("column", columnFamilyMap, db)
	if err != nil {
		return nil, err
	}

	return &ipb.CreateTableInternalResponse{
		Success: true,
	}, nil
}

func (s *TabletServiceServer) DeleteTable(ctx context.Context, req *ipb.DeleteTableInternalRequest) (*ipb.DeleteTableInternalResponse, error) {
	tableName := req.TableName
	// Check if the table exists in the map of tables
	db, ok := s.Tables[tableName]
	if !ok {
		return &ipb.DeleteTableInternalResponse{
			Success: false,
			Message: fmt.Sprintf("Table %s not found", tableName),
		}, nil
	}

	// Close the LevelDB instance first
	err := db.Close()
	if err != nil {
		return &ipb.DeleteTableInternalResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to close LevelDB for table %s: %v", tableName, err),
		}, nil
	}

	dbPath := GetFilePath(s.TabletAddress, tableName)
	// Delete all files under the table
	err = os.RemoveAll(dbPath)
	if err != nil {
		return &ipb.DeleteTableInternalResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to delete LevelDB database at %s: %v", dbPath, err),
		}, nil
	}

	// Remove the reference from the Tables map
	delete(s.Tables, tableName)

	return &ipb.DeleteTableInternalResponse{
		Success: true,
	}, nil
}

// Paper Section2 -> timestamp subsection
// Each cell in a Bigtable can contain multiple versions of
// the same data; these versions are indexed by timestamp.  Different
// versions of a cell are stored in decreasing timestamp order, so that the
// most recent versions can be read first. In our implementation, only the most
// recent 3 version will be kept, older version will be deleted
func (s *TabletServiceServer) Read(ctx context.Context, req *proto.ReadRequest) (*proto.ReadResponse, error) {
	prefix := fmt.Sprintf("%s:%s:%s:", req.RowKey, req.ColumnFamily, req.ColumnQualifier)
	tableName := req.TableName
	tableFile, ok := s.Tables[tableName]
	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Table %s not found", tableName))
	}

	iter := tableFile.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	var allValues []ValueWithKeyAndTimestamps

	// collect all versions
	for iter.Next() {
		key := string(iter.Key())
		value := string(iter.Value())

		timestamp, err := ExtractTimestampFromKey(key)
		if err != nil {
			continue
		}

		allValues = append(allValues, ValueWithKeyAndTimestamps{
			Timestamp: timestamp,
			Value:     value,
			Key:       key, // store key for later delete
		})
	}

	// sort array
	sort.Slice(allValues, func(i, j int) bool {
		return allValues[i].Timestamp > allValues[j].Timestamp
	})

	// TODO: make it configurable
	maxVersion := 3
	if len(allValues) > maxVersion {
		// delete the old versions
		for i := maxVersion; i < len(allValues); i++ {
			deleteErr := tableFile.Delete([]byte(allValues[i].Key), nil)
			if deleteErr != nil {
				log.Printf("Failed to delete old version: %v\n", deleteErr)
			}
		}

		allValues = allValues[:maxVersion]
	}

	var returnValues []*proto.ValueWithTimestamps
	for i := 0; i < len(allValues); i++ {
		if i == int(req.ReturnVersion) {
			break
		}

		returnValues = append(returnValues, &proto.ValueWithTimestamps{
			Timestamp: allValues[i].Timestamp,
			Value:     allValues[i].Value,
		})
	}

	return &proto.ReadResponse{Values: returnValues}, nil
}

func (s *TabletServiceServer) Write(ctx context.Context, req *proto.WriteRequest) (*proto.WriteResponse, error) {
	key := BuildKey(req.RowKey, req.ColumnFamily, req.ColumnQualifier, req.Timestamp)
	tableName := req.TableName
	tableFile, ok := s.Tables[tableName]
	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Table %s not found", tableName))
	}

	err := tableFile.Put([]byte(key), req.Value, nil)
	if err != nil {
		log.Printf("Failed to write to LevelDB: %v\n", err)
		return &proto.WriteResponse{Success: false}, err
	}

	// update row info: persist metadata row so they can be recovered when the server crashes
	rowSet := s.TablesRows[tableName]
	_, exist := rowSet[req.RowKey]
	if !exist {
		// if this row doesn't exist before, add to the set and update persistent
		rowSet[req.RowKey] = struct{}{}
		err = WriteMetaDataToPersistent("row", rowSet, tableFile)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &proto.WriteResponse{Success: true}, nil
}

func (s *TabletServiceServer) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	prefix := fmt.Sprintf("%s:%s:%s:", req.RowKey, req.ColumnFamily, req.ColumnQualifier)

	// check if table exists
	tableName := req.TableName
	tableFile, ok := s.Tables[tableName]
	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Table %s not found", tableName))
	}

	// iterate over all keys
	iter := tableFile.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()
	var keysToDelete []string
	for iter.Next() {
		key := string(iter.Key())
		keysToDelete = append(keysToDelete, key)
	}

	if len(keysToDelete) == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("No value found for key %s", req.RowKey))
	}

	// delete all keys
	for _, key := range keysToDelete {
		deleteErr := tableFile.Delete([]byte(key), nil)
		if deleteErr != nil {
			log.Printf("Failed to delete key %s: %v\n", key, deleteErr)
		}
	}

	// update row info: persist metadata row so they can be recovered when the server crashes
	rowSet := s.TablesRows[tableName]
	delete(rowSet, req.RowKey)
	err := WriteMetaDataToPersistent("row", rowSet, tableFile)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.DeleteResponse{Success: true}, nil
}

func WriteMetaDataToPersistent(metadataType string, data interface{}, db *leveldb.DB) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		log.Fatal(fmt.Sprintf("Encode metadata %v failed: %v", metadataType, err))
	}

	err := db.Put([]byte(fmt.Sprintf("meta_%s", metadataType)), buf.Bytes(), nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Persist %s to LevelDB failed: %v", metadataType, err))
	}

	return nil
}
