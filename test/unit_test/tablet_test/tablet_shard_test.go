package tablet

import (
	"context"
	"final/bigtable/tablet"
	proto "final/proto/external-api"
	ipb "final/proto/internal-api"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"testing"
	"time"
)

const (
	TABLET_ADDRESS  = "TEST_ADDRESS"
	MASTER_ADDRESS  = "localhost:10000"
	TEST_TABLE_NAME = "testdb"
)

func TestMigrate(t *testing.T) {
	RECOVER_ADDRESS := "RECOVER_ADDRESS"
	TEST_TABLE_NAME := "testdb"

	server := tablet.TabletServiceServer{
		TabletAddress: RECOVER_ADDRESS,
		MasterAddress: MASTER_ADDRESS,
		Tables:        make(map[string]*leveldb.DB),
		TablesRows:    make(map[string]map[string]struct{}),
		TablesColumns: make(map[string]map[string][]string),
	}

	SOURCE_TABLE_ADDRESS := TABLET_ADDRESS
	server.MigrateTableToSelf(SOURCE_TABLE_ADDRESS, TEST_TABLE_NAME)

	log.Printf("restored rows : %v", server.TablesRows)
	log.Printf("restored columns : %v", server.TablesColumns)
}

func TestRecoverWrite(t *testing.T) {

	server := tablet.TabletServiceServer{
		TabletAddress: TABLET_ADDRESS,
		MasterAddress: MASTER_ADDRESS,
		Tables:        map[string]*leveldb.DB{},
		TablesRows:    make(map[string]map[string]struct{}),
		TablesColumns: make(map[string]map[string][]string),
	}

	ctx := context.Background()

	columnf := map[string][]string{
		"cf1": {"col1"},
	}

	var cfMsgs []*ipb.CreateTableInternalRequest_ColumnFamily
	for family, columns := range columnf {
		cfMsgs = append(cfMsgs, &ipb.CreateTableInternalRequest_ColumnFamily{
			FamilyName: family,
			Columns:    columns,
		})
	}

	req := ipb.CreateTableInternalRequest{
		TableName:      TEST_TABLE_NAME,
		ColumnFamilies: cfMsgs,
	}

	_, err := server.CreateTable(ctx, &req)
	if err != nil {
		return
	}

	writeRequest1 := &proto.WriteRequest{
		TableName:       TEST_TABLE_NAME,
		RowKey:          "row1",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		Value:           []byte("row1"),
		Timestamp:       time.Now().UnixNano(),
	}

	writeRequest2 := &proto.WriteRequest{
		TableName:       TEST_TABLE_NAME,
		RowKey:          "row2",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		Value:           []byte("row2"),
		Timestamp:       time.Now().UnixNano(),
	}

	writeRequest2_2 := &proto.WriteRequest{
		TableName:       TEST_TABLE_NAME,
		RowKey:          "row2",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		Value:           []byte("row2"),
		Timestamp:       time.Now().UnixNano(),
	}
	writeRequest3 := &proto.WriteRequest{
		TableName:       TEST_TABLE_NAME,
		RowKey:          "row3",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		Value:           []byte("row3"),
		Timestamp:       time.Now().UnixNano(),
	}

	_, writeErr := server.Write(ctx, writeRequest1)
	if writeErr != nil {
		t.Fatal(fmt.Sprintf("Failed to write to LevelDB: %v", writeErr))
	}

	_, writeErr = server.Write(ctx, writeRequest2)
	if writeErr != nil {
		t.Fatal(fmt.Sprintf("Failed to write to LevelDB: %v", writeErr))
	}

	_, writeErr = server.Write(ctx, writeRequest2_2)
	if writeErr != nil {
		t.Fatal(fmt.Sprintf("Failed to write to LevelDB: %v", writeErr))
	}

	_, writeErr = server.Write(ctx, writeRequest3)

	if writeErr != nil {
		t.Fatal(fmt.Sprintf("Failed to write to LevelDB: %v", writeErr))
	}
}
