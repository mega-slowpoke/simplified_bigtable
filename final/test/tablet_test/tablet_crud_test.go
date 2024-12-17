package tablet_test

import (
	"context"
	"final/bigtable/tablet"
	proto "final/proto/external-api"
	ipb "final/proto/internal-api"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"path/filepath"
	"testing"
	"time"
)

// setup before each test
//
//	func TestMain(m *testing.M) {
//		var err error
//		// Clean up previous database files if any
//		os.RemoveAll("testdb")
//
//		// Open a fresh LevelDB instance
//		db, err = leveldb.OpenFile("testdb", nil)
//		if err != nil {
//			log.Fatalf("Failed to open LevelDB: %v", err)
//		}
//		// Ensure the database is properly closed after the tests
//		defer db.Close()
//
//		// Run the tests
//		code := m.Run()
//
//		// Clean up the database directory after the tests
//		os.RemoveAll("testdb")
//		os.Exit(code)
//	}

const (
	TABLET_ADDRESS  = "localhost:9001"
	MASTER_ADDRESS  = "localhost:10000"
	TEST_TABLE_NAME = "testdb"
)

func TestCreateTable(t *testing.T) {
	server := tablet.TabletServiceServer{
		TabletAddress: TABLET_ADDRESS,
		MasterAddress: MASTER_ADDRESS,
		Tables:        make(map[string]*leveldb.DB),
	}

	req := ipb.CreateTableInternalRequest{
		TableName: TEST_TABLE_NAME,
	}

	_, err := server.CreateTable(context.Background(), &req)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to create table: %v", err))
	}

	_, exist := server.Tables[TEST_TABLE_NAME]
	if !exist {
		log.Fatal(fmt.Sprintf("created table is not added to the map"))
	}
}

func TestCreateDuplicateTable(t *testing.T) {
	var db leveldb.DB
	server := tablet.TabletServiceServer{
		TabletAddress: TABLET_ADDRESS,
		MasterAddress: MASTER_ADDRESS,
		Tables: map[string]*leveldb.DB{
			TEST_TABLE_NAME: &db,
		},
	}

	req := ipb.CreateTableInternalRequest{
		TableName: TEST_TABLE_NAME,
	}

	_, err := server.CreateTable(context.Background(), &req)
	if err == nil {
		log.Fatal(fmt.Sprintf("Table exists, You shouldn't create duplicate table %v", err))
	}
}

func TestDeleteTable(t *testing.T) {
	server := tablet.TabletServiceServer{
		TabletAddress: TABLET_ADDRESS,
		MasterAddress: MASTER_ADDRESS,
		Tables:        make(map[string]*leveldb.DB),
	}

	req := ipb.CreateTableInternalRequest{
		TableName: TEST_TABLE_NAME,
	}

	_, err := server.CreateTable(context.Background(), &req)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to create table %v", err))
	}

	ctx := context.Background()
	timeNow := time.Now().UnixNano()
	writeRequest := &proto.WriteRequest{
		TableName:       TEST_TABLE_NAME,
		RowKey:          "row1",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		Value:           []byte("write1"),
		Timestamp:       timeNow,
	}
	_, writeErr := server.Write(ctx, writeRequest)

	if writeErr != nil {
		t.Fatalf("Failed to write to LevelDB: %v", writeErr)
	}

	server.Write(ctx, writeRequest)

	readRequest := &proto.ReadRequest{
		TableName:       TEST_TABLE_NAME,
		RowKey:          "row1",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		ReturnVersion:   1,
	}

	readResponse, readErr := server.Read(context.Background(), readRequest)
	if readErr != nil {
		t.Fatalf("Failed to read from LevelDB: %v", readErr)
	}

	if len(readResponse.Values) != 1 {
		t.Fatalf("Read returned wrong number of values, expected %d, got %d", 1, len(readResponse.Values))
	}

	deleteReq := &ipb.DeleteTableInternalRequest{
		TableName: TEST_TABLE_NAME,
	}

	_, err = server.DeleteTable(ctx, deleteReq)
	if err != nil {
		log.Fatal(fmt.Sprintf("delete failed"))
	}

	readResponse, readErr = server.Read(context.Background(), readRequest)
	if readErr == nil {
		t.Fatalf("Delete but still can read")
	}

	_, exist := server.Tables[TEST_TABLE_NAME]
	if exist {
		log.Fatal(fmt.Sprintf("delete didn't remove table from the map"))
	}
}

func TestTabletSingleWrite(t *testing.T) {
	db, _ := leveldb.OpenFile(filepath.Join(TABLET_ADDRESS, TEST_TABLE_NAME), nil)

	server := tablet.TabletServiceServer{
		TabletAddress: TABLET_ADDRESS,
		MasterAddress: MASTER_ADDRESS,
		Tables: map[string]*leveldb.DB{
			TEST_TABLE_NAME: db,
		},
	}

	ctx := context.Background()
	timeNow := time.Now().UnixNano()
	writeRequest := &proto.WriteRequest{
		TableName:       TEST_TABLE_NAME,
		RowKey:          "row1",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		Value:           []byte("write1"),
		Timestamp:       timeNow,
	}
	_, writeErr := server.Write(ctx, writeRequest)

	if writeErr != nil {
		t.Fatalf("Failed to write to LevelDB: %v", writeErr)
	}
}

func TestTabletSingleRead(t *testing.T) {
	db, _ := leveldb.OpenFile(filepath.Join(TABLET_ADDRESS, TEST_TABLE_NAME), nil)

	server := tablet.TabletServiceServer{
		TabletAddress: TABLET_ADDRESS,
		MasterAddress: MASTER_ADDRESS,
		Tables: map[string]*leveldb.DB{
			TEST_TABLE_NAME: db,
		},
	}

	readRequest := &proto.ReadRequest{
		TableName:       TEST_TABLE_NAME,
		RowKey:          "row1",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		ReturnVersion:   1,
	}

	readResponse, readErr := server.Read(context.Background(), readRequest)
	if readErr != nil {
		t.Fatalf("Failed to read from LevelDB: %v", readErr)
	}

	if len(readResponse.Values) != 1 {
		t.Fatalf("Read returned wrong number of values, expected %d, got %d", 1, len(readResponse.Values))
	}

	t.Logf("%v", readResponse.Values)
	expected := "write1"
	actual := readResponse.Values[0].Value
	if actual != expected {
		t.Fatalf("Expected value %s, but got %s", expected, actual)
	}
}

// TODO: SingleWriteAndRead might lead to problem
func TestTabletSingleWriteAndRead(t *testing.T) {
	db, _ := leveldb.OpenFile(filepath.Join(TABLET_ADDRESS, TEST_TABLE_NAME), nil)

	server := tablet.TabletServiceServer{
		TabletAddress: TABLET_ADDRESS,
		MasterAddress: MASTER_ADDRESS,
		Tables: map[string]*leveldb.DB{
			TEST_TABLE_NAME: db,
		},
	}

	ctx := context.Background()
	timeNow := time.Now().UnixNano()
	writeRequest := &proto.WriteRequest{
		TableName:       TEST_TABLE_NAME,
		RowKey:          "row1",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		Value:           []byte("write1"),
		Timestamp:       timeNow,
	}
	_, writeErr := server.Write(ctx, writeRequest)
	if writeErr != nil {
		return
	}

	if writeErr != nil {
		t.Fatalf("Failed to write to LevelDB: %v", writeErr)
	}

	readRequest := &proto.ReadRequest{
		TableName:       TEST_TABLE_NAME,
		RowKey:          "row1",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		ReturnVersion:   1,
	}

	readResponse, readErr := server.Read(context.Background(), readRequest)
	if readErr != nil {
		t.Fatalf("Failed to read from LevelDB: %v", readErr)
	}

	if len(readResponse.Values) != 1 {
		t.Fatalf("Read returned wrong number of values, expected %d, got %d", 1, len(readResponse.Values))
	}

	t.Logf("%v", readResponse.Values)
	expected := "write1"
	actual := readResponse.Values[0].Value
	if actual != expected {
		t.Fatalf("Expected value %s, but got %s", expected, actual)
	}
}

func TestTabletMultipleRead(t *testing.T) {
	db, _ := leveldb.OpenFile(filepath.Join(TABLET_ADDRESS, TEST_TABLE_NAME), nil)

	server := tablet.TabletServiceServer{
		TabletAddress: TABLET_ADDRESS,
		MasterAddress: MASTER_ADDRESS,
		Tables: map[string]*leveldb.DB{
			TEST_TABLE_NAME: db,
		},
	}
	ctx := context.Background()

	readRequest := &proto.ReadRequest{
		TableName:       TEST_TABLE_NAME,
		RowKey:          "row1",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		ReturnVersion:   5,
	}

	readResponse, readErr := server.Read(ctx, readRequest)
	if readErr != nil {
		t.Fatalf("Failed to read from LevelDB: %v", readErr)
	}

	t.Logf("element: %v", readResponse.Values)
}

func TestTabletDelete(t *testing.T) {
	db, _ := leveldb.OpenFile(filepath.Join(TABLET_ADDRESS, TEST_TABLE_NAME), nil)

	server := tablet.TabletServiceServer{
		TabletAddress: TABLET_ADDRESS,
		MasterAddress: MASTER_ADDRESS,
		Tables: map[string]*leveldb.DB{
			TEST_TABLE_NAME: db,
		},
	}
	ctx := context.Background()

	deleteReq := &proto.DeleteRequest{
		TableName:       TEST_TABLE_NAME,
		RowKey:          "row1",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
	}

	deleteResponse, deleteErr := server.Delete(ctx, deleteReq)
	if deleteErr != nil {
		t.Fatalf("Failed to delete from LevelDB: %v", deleteErr)
	}

	if deleteResponse.Success != true {
		t.Fatalf("Expected success to be true, but got %v", deleteResponse.Success)
	}

	readRequest := &proto.ReadRequest{
		TableName:       TEST_TABLE_NAME,
		RowKey:          "row1",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		ReturnVersion:   5,
	}

	readResponse, readErr := server.Read(ctx, readRequest)
	if readErr != nil {
		t.Fatalf("Failed to read from LevelDB: %v", readErr)
	}

	if len(readResponse.Values) != 0 {
		t.Fatalf("Deleted but still read value %v", readResponse.Values)
	}
}
