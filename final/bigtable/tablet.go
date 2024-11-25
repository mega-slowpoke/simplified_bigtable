package bigtable

import (
	"context"
	proto "final/proto/external-api"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"strconv"
	"strings"
)

// TabletService for gRPC
type TabletServiceServer struct {
	proto.UnimplementedTabletServiceServer
}

func NewTabletService() *TabletServiceServer {
	return &TabletServiceServer{}
}

func (s *TabletServiceServer) Read(ctx context.Context, req *proto.ReadRequest) (*proto.ReadResponse, error) {
	prefix := fmt.Sprintf("%s:%s:%s:", req.RowKey, req.ColumnFamily, req.ColumnQualifier)
	tableFile, err := leveldb.OpenFile(req.TableName, nil)

	defer func(tableFile *leveldb.DB) {
		closeErr := tableFile.Close()
		if err != nil {
			log.Fatalf(closeErr.Error())
		}
	}(tableFile)

	iter := tableFile.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	values := make([]*proto.ValueWithTimestamps, 0)
	var count int64 = 0

	for iter.Next() {
		if count >= req.MaxVersion {
			break
		}

		key := string(iter.Key())
		value := string(iter.Value())
		timestamp, err := extractTimestampFromKey(key)
		if err != nil {
			return nil, err
		}

		values = append(values, &proto.ValueWithTimestamps{
			Timestamp: timestamp,
			Value:     value,
		})
		count++
	}

	return &proto.ReadResponse{Values: values}, nil
}

func (s *TabletServiceServer) Write(ctx context.Context, req *proto.WriteRequest) (*proto.WriteResponse, error) {
	key := buildKey(req.RowKey, req.ColumnFamily, req.ColumnQualifier, req.Timestamp)
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
func (s *TabletServiceServer) Scan(req *proto.ScanRequest, stream proto.TabletService_ScanServer) error {
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

	return nil
}

func buildKey(rowKey, columnFamily, columnQualifier string, timestamp int64) string {
	reversedTimestamp := ^timestamp // descendingOrder timestamp so we can retrieve the most recent data first
	return fmt.Sprintf("%s:%s:%s:%d", rowKey, columnFamily, columnQualifier, reversedTimestamp)
}

func extractTimestampFromKey(key string) (int64, error) {
	parts := strings.Split(key, ":")
	if len(parts) < 4 {
		return 0, fmt.Errorf("invalid key format") // Return an error if the key is in an invalid format
	}

	// Extract the reversed timestamp (the 4th part)
	reversedTimestampStr := parts[3]
	reversedTimestamp, err := strconv.ParseInt(reversedTimestampStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse timestamp: %v", err) // Return an error if parsing fails
	}

	return reversedTimestamp, nil
}
