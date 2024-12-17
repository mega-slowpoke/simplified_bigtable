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
	"log"
)

//tables_rows ={
//	"customers": {
//		"customer_0001", "customer_0002", "customer_0003", // 客户表中的部分行键示例，假设行键是客户编号
//		},
//		"orders": {
//		"order_0001", "order_0002", "order_0003", "order_0004", // 订单表中的部分行键示例，假设行键是订单编号
//		},
//	}

//tables_columns = {
//		"customers": {
//			"basic_info": {
//				"customer_name", "customer_email", "customer_phone",
//			},
//			"address_info": {
//				"province", "city", "street", "zip_code",
//			},
//		},
//		"orders": {
//			"order_detail": {
//				"product_name", "product_quantity", "product_price",
//			},
//			"order_status": {
//				"status", "update_time",
//			},
//		},
//   }
//

// TabletService
type TabletServiceServer struct {
	epb.UnimplementedTabletExternalServiceServer
	ipb.UnimplementedTabletInternalServiceServer
	Tables        map[string]*leveldb.DB // tableName -> levelDB
	TabletAddress string
	MasterAddress string
	MasterConn    *grpc.ClientConn
	MasterClient  *ipb.MasterInternalServiceClient
	MaxTableCnt   int
	TablesRows    map[string]map[string]struct{}
	TablesColumns map[string]map[string][]string
}

func NewTabletService(opt SetupOptions) (*TabletServiceServer, error) {
	//var masterClient ipb.MasterInternalServiceClient
	//var conn *grpc.ClientConn

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
		MaxTableCnt:   opt.MaxTableSize,
		TablesRows:    make(map[string]map[string]struct{}),
		TablesColumns: make(map[string]map[string][]string),
	}, nil
}

func (s *TabletServiceServer) CloseMasterConnection() {
	err := s.MasterConn.Close()
	if err != nil {
		return
	}
}

func (s *TabletServiceServer) CloseAllTables() {
	for name, db := range s.Tables {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close table %s: %v\n", name, err)
		}
	}
}
