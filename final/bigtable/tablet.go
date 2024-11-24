package bigtable

import (
	proto "final/proto/external-api"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
)

type TabletServiceServer struct {
	db *leveldb.DB
}

func NewTabletServiceServer(dbPath string) (*TabletServiceServer, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open LevelDB: %v", err)
	}

	return &TabletServiceServer{db: db}, nil
}

func (server *TabletServiceServer) Read(req *proto.ReadRequest) proto.ReadResponse {

}

func (server *TabletServiceServer) Write(req *proto.WriteRequest) proto.WriteResponse {

}

func (server *TabletServiceServer) Delete(req *proto.DeleteRequest) proto.DeleteResponse {

}

func (server *TabletServiceServer) Scan(req *proto.ScanRequest) proto.ScanResponse {

}
