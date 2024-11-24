package bigtable

import (
	"context"
	proto "final/proto/external-api"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"strings"
	"time"
)

func NewTabletServiceServer(dbPath string) (*TabletServiceServer, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open LevelDB: %v", err)
	}

	return &TabletServiceServer{db: db}, nil
}

// TabletService for gRPC
type TabletServiceServer struct {
	db      *leveldb.DB
	locker  *ChubbyLock
	lockKey string
	proto.UnimplementedTabletServiceServer
}

func NewTabletService(db *leveldb.DB, locker *ChubbyLock, lockKey string) *TabletServiceServer {
	return &TabletServiceServer{
		db:      db,
		locker:  locker,
		lockKey: lockKey,
	}
}

func (s *TabletServiceServer) Read(ctx context.Context, req *proto.ReadRequest) (*proto.ReadResponse, error) {
	// Try to acquire lock within 5 sec, otherwise exit
	if !s.locker.Lock(5 * time.Second) {
		log.Printf("Failed to acquire lock for Read operation")
		return nil, fmt.Errorf("failed to acquire lock")
	}
	defer s.locker.Unlock()

	key := buildKey(req.TableName, req.RowKey, req.ColumnFamily, req.ColumnQualifier)
	value, err := s.db.Get([]byte(key), nil)
	if err != nil {
		log.Printf("Read error: %v", err)
		return nil, err
	}

	return &proto.ReadResponse{
		Value:     value,
		Timestamp: extractTimestamp(value),
	}, nil
}

func (s *TabletServiceServer) Write(ctx context.Context, req *proto.WriteRequest) (*proto.WriteResponse, error) {
	if !s.locker.Lock(5 * time.Second) {
		log.Printf("Failed to acquire lock for Write operation")
		return nil, fmt.Errorf("failed to acquire lock")
	}
	defer s.locker.Unlock()

	key := buildKey(req.TableName, req.RowKey, req.ColumnFamily, req.ColumnQualifier)
	value := append(req.Value, []byte(fmt.Sprintf(":%d", req.Timestamp))...)

	if err := s.db.Put([]byte(key), value, nil); err != nil {
		log.Printf("Write error: %v", err)
		return &proto.WriteResponse{Success: false, ErrorMessage: err.Error()}, err
	}

	return &proto.WriteResponse{Success: true}, nil
}

func (s *TabletServiceServer) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	if !s.locker.Lock(5 * time.Second) {
		log.Printf("Failed to acquire lock for Write operation")
		return nil, fmt.Errorf("failed to acquire lock")
	}
	defer s.locker.Unlock()

	key := buildKey(req.TableName, req.RowKey, req.ColumnFamily, req.ColumnQualifier)
	if err := s.db.Delete([]byte(key), nil); err != nil {
		log.Printf("Delete error: %v", err)
		return &proto.DeleteResponse{Success: false, ErrorMessage: err.Error()}, err
	}

	return &proto.DeleteResponse{Success: true}, nil
}

// Scan
func (s *TabletServiceServer) Scan(req *proto.ScanRequest, stream proto.TabletService_ScanServer) error {
	// TODO: how to add lock to Scan

	startKey := buildKey(req.TableName, req.StartRowKey, req.ColumnFamily, req.ColumnQualifier)
	endKey := buildKey(req.TableName, req.EndRowKey, req.ColumnFamily, req.ColumnQualifier)

	iter := s.db.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		key := string(iter.Key())
		if key >= startKey && key <= endKey {
			value := iter.Value()
			timestamp := extractTimestamp(value)
			if err := stream.Send(&proto.ScanResponse{
				RowKey:    key,
				Value:     value,
				Timestamp: timestamp,
			}); err != nil {
				log.Printf("Scan send error: %v", err)
				return err
			}
		}
	}

	if err := iter.Error(); err != nil {
		log.Printf("Scan iterator error: %v", err)
		return err
	}

	return nil
}

func buildKey(tableName, rowKey, columnFamily, columnQualifier string) string {
	return tableName + ":" + rowKey + ":" + columnFamily + ":" + columnQualifier
}

func extractTimestamp(value []byte) int64 {
	parts := strings.Split(string(value), ":")
	if len(parts) > 1 {
		var timestamp int64
		if _, err := fmt.Sscanf(parts[len(parts)-1], "%d", &timestamp); err == nil {
			return timestamp
		}
	}
	return 0
}
