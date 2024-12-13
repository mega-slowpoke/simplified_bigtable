package tablet

import (
	epb "final/proto/external-api"
	ipb "final/proto/internal-api"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// TabletService for gRPC
type TabletServiceServer struct {
	epb.UnimplementedTabletServiceServer
	Tables        map[string]*leveldb.DB // tableName -> levelDB
	TabletAddress string
	MasterAddress string
	MasterConn    *grpc.ClientConn
	MasterClient  *ipb.MasterInternalServiceClient
}

func NewTabletService(opt SetupOptions) (*TabletServiceServer, error) {
	conn, err := grpc.NewClient(opt.MasterAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, status.Error(codes.Unavailable, fmt.Sprint("Tablet %s Could not connect to master", opt.TabletAddress))
	}
	masterClient := ipb.NewMasterInternalServiceClient(conn)

	return &TabletServiceServer{
		Tables:        make(map[string]*leveldb.DB),
		TabletAddress: opt.TabletAddress,
		MasterAddress: opt.MasterAddress, // "localhost:12345"
		MasterConn:    conn,
		MasterClient:  &masterClient,
	}, nil
}

func (s *TabletServiceServer) CloseMasterConnection() {
	err := s.MasterConn.Close()
	if err != nil {
		return
	}
}
