package tablet

import (
	"bytes"
	"context"
	"encoding/gob"
	ipb "final/proto/internal-api"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"path/filepath"
)

func (s *TabletServiceServer) RecoverCrashedTablet(ctx context.Context, req *ipb.RecoveryRequest) (*ipb.RecoveryResponse, error) {
	crashedTabletAddress := req.CrashedTabletAddress
	tableName := req.TableName
	err := s.MigrateTableToSelf(crashedTabletAddress, tableName)
	if err != nil {
		return &ipb.RecoveryResponse{
			Success: false,
		}, status.Error(codes.Internal, err.Error())
	}

	return &ipb.RecoveryResponse{
		Success: true,
	}, nil
}

func (s *TabletServiceServer) MigrateTableToSelf(sourceServerAddress string, tableName string) error {
	// Step1: create table if not exist
	dbPath := GetFilePath(s.TabletAddress, tableName)
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		logrus.Debugf("create table for recovery server failed: %v", err)
		return status.Errorf(codes.Internal, "create table for recovery server failed: %v", err)
	}
	s.Tables[tableName] = db
	db.Close() // close the db first so that the recoverData can access data

	// move data (table contents and metadata) from "sourceServerAddress/tableName" to "s.TabletAddress/tableName"
	err = s.recoverData(sourceServerAddress, tableName)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to recover actual data: %v", err)
	}

	// recover metadata first
	db, err = leveldb.OpenFile(dbPath, nil)
	if err != nil {
		logrus.Debugf("create table for recovery server failed: %v", err)
		return status.Errorf(codes.Internal, "create table for recovery server failed: %v", err)
	}
	s.Tables[tableName] = db

	err = s.rebuildColumnsMetadata(db)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to rebuild columns metadata: %v", err)
	}

	err = s.rebuildRowsMetadata(db)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to rebuild rows metadata: %v", err)
	}

	logrus.Info(fmt.Sprintf("Recovery completed successfully for: %s - %s, data is moved to %s", sourceServerAddress, tableName, s.TabletAddress))
	return nil
}

func (s *TabletServiceServer) rebuildColumnsMetadata(db *leveldb.DB) error {
	data, err := db.Get([]byte("meta_row"), nil)
	if err != nil {
		logrus.Fatalf("Read from LevelDB failed: %v", err)
	}

	// decode serialized data to rebuild rowSet
	var restoredRowSet map[string]struct{}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err = decoder.Decode(&restoredRowSet); err != nil {
		log.Fatal("Decode failed:", err)
	}

	// just to check restored row set data when debugging
	logrus.Debug("Restored Set: %v", restoredRowSet)
	return nil
}

func (s *TabletServiceServer) rebuildRowsMetadata(db *leveldb.DB) error {
	data, err := db.Get([]byte("meta_column"), nil)
	if err != nil {
		logrus.Fatalf("Read from LevelDB failed: %v", err)
	}

	// decode serialized data to rebuild rowSet
	var restoredColumnFamilies map[string][]string
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err = decoder.Decode(&restoredColumnFamilies); err != nil {
		return status.Errorf(codes.Internal, "decoder failed: %v", err)
	}

	// just to check restored row set data when debugging
	logrus.Debug("Restored ColumnFamilies: %v", restoredColumnFamilies)
	return nil
}

// recoverActualData - mimic moving data in crashedServerAddress to current tablet server
func (s *TabletServiceServer) recoverData(crashedServerAddress string, tableName string) error {
	sourceDataPath := GetFilePath(crashedServerAddress, tableName)
	destDataPath := GetFilePath(s.TabletAddress, tableName)

	// Ensure destination exists
	if err := os.MkdirAll(filepath.Dir(destDataPath), os.ModePerm); err != nil {
		return status.Errorf(codes.Internal, "failed to create dest data directory: %v", err)
	}

	// Move LevelDB files
	if err := moveDir(sourceDataPath, destDataPath); err != nil {
		return status.Errorf(codes.Internal, "failed to move data directory: %v", err)
	}

	logrus.Debug("move level db succeed")
	return nil
}

// moveDir - Utility function to move a directory (recursively)
func moveDir(sourceDir, destDir string) error {
	// Remove destination directory if it already exists
	if _, err := os.Stat(destDir); err == nil {
		if err = os.RemoveAll(destDir); err != nil {
			return fmt.Errorf("failed to remove existing destination directory: %v", err)
		}
	}

	// Move directory
	if err := os.Rename(sourceDir, destDir); err != nil {
		return err
	}
	return nil

	////
	//err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
	//	if err != nil {
	//		return err
	//	}
	//	relativePath, err := filepath.Rel(sourceDir, path)
	//	if err != nil {
	//		return err
	//	}
	//	destFilePath := filepath.Join(destDir, relativePath)
	//	if info.IsDir() {
	//		return os.MkdirAll(destFilePath, info.Mode())
	//	}
	//	data, err := os.ReadFile(path)
	//	if err != nil {
	//		return err
	//	}
	//	return os.WriteFile(destFilePath, data, info.Mode())
	//})
	//if err != nil {
	//	return fmt.Errorf("failed to copy LevelDB data: %v", err)
	//}

	//db, ok := s.Tables[tableName]
	//if ok {
	//	err := db.Close()
	//	if err != nil {
	//		return fmt.Errorf("failed to close LevelDB file: %v", err)
	//	}
	//}
	//
	//// 3.
	//if err := os.RemoveAll(sourceDataPath); err != nil {
	//	return fmt.Errorf("failed to remove source LevelDB data directory: %v", err)
	//}
	//return nil
	return nil
}
