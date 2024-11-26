package test

import (
	"context"
	"final/bigtable"
	proto "final/proto/external-api"
	"github.com/syndtr/goleveldb/leveldb"
	"testing"
	"time"
)

var db *leveldb.DB

// setup before each test
//func TestMain(m *testing.M) {
//	var err error
//	// Clean up previous database files if any
//	os.RemoveAll("testdb")
//
//	// Open a fresh LevelDB instance
//	db, err = leveldb.OpenFile("testdb", nil)
//	if err != nil {
//		log.Fatalf("Failed to open LevelDB: %v", err)
//	}
//	// Ensure the database is properly closed after the tests
//	defer db.Close()
//
//	// Run the tests
//	code := m.Run()
//
//	// Clean up the database directory after the tests
//	os.RemoveAll("testdb")
//	os.Exit(code)
//}

func TestTabletSingleWriteAndRead(t *testing.T) {
	server := bigtable.NewTabletService()
	timeNow := time.Now().UnixNano()
	writeRequest := &proto.WriteRequest{
		TableName:       "testdb",
		RowKey:          "row1",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		Value:           []byte("write1"),
		Timestamp:       timeNow,
	}
	ctx := context.Background()
	_, writeErr := server.Write(ctx, writeRequest)
	if writeErr != nil {
		return
	}

	if writeErr != nil {
		t.Fatalf("Failed to write to LevelDB: %v", writeErr)
	}

	readRequest := &proto.ReadRequest{
		TableName:       "testdb",
		RowKey:          "row1",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		ReturnVersion:   1,
	}

	readResponse, readErr := server.Read(ctx, readRequest)
	if readErr != nil {
		t.Fatalf("Failed to read from LevelDB: %v", readErr)
	}

	if len(readResponse.Values) != 1 {
		t.Fatalf("Read returned wrong number of values, expected %d, got %d", 1, len(readResponse.Values))
	}

	expected := "write1"
	actual := readResponse.Values[0].Value
	if actual != expected {
		t.Fatalf("Expected value %s, but got %s", expected, actual)
	}

}

func TestTabletMultipleRead(t *testing.T) {
	server := bigtable.NewTabletService()
	ctx := context.Background()

	readRequest := &proto.ReadRequest{
		TableName:       "testdb",
		RowKey:          "row1",
		ColumnFamily:    "cf1",
		ColumnQualifier: "col1",
		ReturnVersion:   1,
	}

	readResponse, readErr := server.Read(ctx, readRequest)
	if readErr != nil {
		t.Fatalf("Failed to read from LevelDB: %v", readErr)
	}

	for i, r := range readResponse.Values {
		t.Logf("%d element: %v", i, r)
	}
}
