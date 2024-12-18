package tablet

import (
	"final/bigtable/tablet"
	"github.com/syndtr/goleveldb/leveldb"
	"testing"
)

const (
	MASTER_ADDRESS = "localhost:10000"
)

func TestMigrate(t *testing.T) {
	TABLET_ADDRESS := "RECOVER_ADDRESS"
	TEST_TABLE_NAME := "testdb"

	server := tablet.TabletServiceServer{
		TabletAddress: TABLET_ADDRESS,
		MasterAddress: MASTER_ADDRESS,
		Tables:        make(map[string]*leveldb.DB),
	}

	SOURCE_TABLE_ADDRESS := "TEST_ADDRESS"
	server.MigrateTableToSelf(SOURCE_TABLE_ADDRESS, TEST_TABLE_NAME)

}
