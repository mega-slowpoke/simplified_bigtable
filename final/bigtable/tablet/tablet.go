package tablet

import (
	"context"
	proto "final/proto/external-api"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"path/filepath"
	"sort"
)

const (
	dbPath = "mockGFS" // use local storage to mimic GFS
)

// TabletService for gRPC
type TabletServiceServer struct {
	proto.UnimplementedTabletServiceServer
	tables  map[string]*leveldb.DB // tableName -> levelDB
	Address string
}

type ValueWithKeyAndTimestamps struct {
	Timestamp int64
	Value     string
	Key       string
}

func NewTabletService(address string) *TabletServiceServer {
	return &TabletServiceServer{
		tables:  make(map[string]*leveldb.DB),
		Address: address,
	}
}

func (s *TabletServiceServer) CreateTable(ctx context.Context, req *proto.CreateTableRequest) (*proto.CreateTableResponse, error) {
	tableName := req.TableName
	db, err := leveldb.OpenFile(filepath.Join(dbPath, tableName), nil)
	if err != nil {
		return &proto.CreateTableResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Failed to create table: %v", err),
		}, nil
	}
	s.tables[tableName] = db
	return &proto.CreateTableResponse{
		Success:      true,
		ErrorMessage: "",
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
	tableFile, err := leveldb.OpenFile(req.TableName, nil)
	if err != nil {
		return nil, err
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
				log.Printf("Failed to delete old version: %v\n", err)
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
	tableFile, err := leveldb.OpenFile(req.TableName, nil)

	defer func(tableFile *leveldb.DB) {
		closeErr := tableFile.Close()
		if err != nil {
			log.Println(closeErr.Error())
		}
	}(tableFile)

	err = tableFile.Put([]byte(key), req.Value, nil)
	if err != nil {
		log.Printf("Failed to write to LevelDB: %v\n", err)
		return &proto.WriteResponse{Success: false}, err
	}

	return &proto.WriteResponse{Success: true}, nil
}

func (s *TabletServiceServer) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	//key := buildKey(req.RowKey, req.ColumnFamily, req.ColumnQualifier)
	//if err := s.db.Delete([]byte(key), nil); err != nil {
	//	log.Printf("Delete error: %v", err)
	//	return &proto.DeleteResponse{Success: false, ErrorMessage: err.Error()}, err
	//}
	//
	//return &proto.DeleteResponse{Success: true}, nil
	return nil, nil
}

// Scan
//func (s *TabletServiceServer) Scan(req *proto.ScanRequest, stream proto.TabletService_ScanServer) error {
// TODO: how to add lock to Scan

//startKey := buildKey(req.StartRowKey, req.ColumnFamily, req.ColumnQualifier, )
//endKey := buildKey(req.EndRowKey, req.ColumnFamily, req.ColumnQualifier)

//iter := s.db.NewIterator(nil, nil)
//defer iter.Release()
//
//for iter.Next() {
//	key := string(iter.Key())
//	if key >= startKey && key <= endKey {
//		value := iter.Value()
//		timestamp := extractTimestampFromKey(string(value))
//		if err := stream.Send(&proto.ScanResponse{
//			RowKey:    key,
//			Value:     value,
//			Timestamp: timestamp,
//		}); err != nil {
//			log.Printf("Scan send error: %v", err)
//			return err
//		}
//	}
//}
//
//if err := iter.Error(); err != nil {
//	log.Printf("Scan iterator error: %v", err)
//	return err
//}

//	return nil
//}

//func (s *TabletServiceServer) DeleteTable(ctx context.Context, req *proto.DeleteTableRequest) (*pb.Response, error) {
//	tableName := req.TableName
//	db, ok := s.dbs[tableName]
//	if !ok {
//		return &pb.Response{
//			Success: false,
//			Message: "Table not found",
//		}, nil
//	}
//	err := db.Close()
//	if err != nil {
//		return &pb.Response{
//			Success: false,
//			Message: fmt.Sprintf("Failed to close table: %v", err),
//		}, nil
//	}
//	delete(s.dbs, tableName)
//	return &pb.Response{
//		Success: true,
//		Message: "Table deleted successfully",
//	}, nil
//}
