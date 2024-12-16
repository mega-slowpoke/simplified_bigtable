package tablet

import (
	"context"
	proto "final/proto/external-api"
	ipb "final/proto/internal-api"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type ValueWithKeyAndTimestamps struct {
	Timestamp int64
	Value     string
	Key       string
}

func (s *TabletServiceServer) CreateTable(ctx context.Context, req *ipb.CreateTableInternalRequest) (*ipb.CreateTableInternalResponse, error) {
	tableName := req.TableName
	db, err := leveldb.OpenFile(filepath.Join(s.TabletAddress, tableName), nil)
	if err != nil {
		return &ipb.CreateTableInternalResponse{
			Success: false,
		}, nil
	}
	s.Tables[tableName] = db
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

	// Get the file path of the LevelDB database
	dbPath := filepath.Join(s.TabletAddress, tableName)
	// Walk through the directory and delete files
	err = filepath.Walk(dbPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return os.Remove(path)
		}
		return os.Remove(path)
	})
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

	defer tableFile.Close()

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
	defer tableFile.Close()

	err := tableFile.Put([]byte(key), req.Value, nil)
	if err != nil {
		log.Printf("Failed to write to LevelDB: %v\n", err)
		return &proto.WriteResponse{Success: false}, err
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
	defer tableFile.Close()

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

	return &proto.DeleteResponse{Success: true}, nil
}
