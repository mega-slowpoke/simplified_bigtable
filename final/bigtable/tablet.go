package bigtable

import (
	"context"
	proto "final/proto/external-api"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"sort"
	"strconv"
	"strings"
)

// TabletService for gRPC
type TabletServiceServer struct {
	proto.UnimplementedTabletServiceServer
}

type ValueWithKeyAndTimestamps struct {
	Timestamp int64
	Value     string
	Key       string
}

func NewTabletService() *TabletServiceServer {
	return &TabletServiceServer{}
}

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

		timestamp, err := extractTimestampFromKey(key)
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
	timestamp := ^reversedTimestamp
	return timestamp, nil
}
