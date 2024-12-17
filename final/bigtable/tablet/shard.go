package tablet

import (
	"context"
	ipb "final/proto/internal-api"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

//-------------------------- Tablet As a grpc service consumer ------------------------------

func (s *TabletServiceServer) PeriodicallyCheckMaxSize(ctx context.Context, period int) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping PeriodicallyCheckMaxSize...")
			return
		default:
			time.Sleep(time.Microsecond * time.Duration(period))
			s.checkAndNotifyMasterForShard()
		}
	}
}

// checkAndNotifyMasterForShard if the size of any table rows is greater than max size, notify the master to issue a shard command
// Problem: what if this table is being modified during this time
func (s *TabletServiceServer) checkAndNotifyMasterForShard() error {
	if len(s.Tables) < s.MaxTableCnt {
		return nil
	}

	client := *(s.MasterClient)

	idx := 0
	toMoveTables := make([]string, 0)
	// collect tables that will be moved to another tablet
	for tableName, _ := range s.Tables {
		if idx >= s.MaxTableCnt {
			toMoveTables = append(toMoveTables, tableName)
		}
		idx++
	}

	for _, tableName := range toMoveTables {
		req := &ipb.ShardRequest{
			TabletAddress: s.TabletAddress,
			TableName:     tableName,
		}

		response, err := client.NotifyShardRequest(context.Background(), req)
		if err != nil {
			log.Printf("failed to notify master for shard request: %v", err)
		}

		err = s.notifyTabletServerForShardUpdate(tableName, response.TargetTabletAddress)
		if err != nil {
			return err
		}

		//for tableName := range s.TablesRows {
		//	if len(tablesRows[tableName]) >= s.MaxTableCnt {
		//		rowKeys := make([]string, len(tablesRows[tableName]))
		//		copy(rowKeys, tablesRows[tableName])
		//		sort.Strings(rowKeys) // TODO: check the sort? sort by what
		//
		//		shardRowKeys := rowKeys[:s.MaxTableCnt/2+1]
		//		remainingRowKeys := rowKeys[s.MaxTableCnt/2:]
		//
		//		req := &ipb.ShardRequest{
		//			TabletAddress: s.TabletAddress,
		//			TableName:     tableName,
		//		}
		//
		//		_, err := client.NotifyShardRequest(context.Background(), req)
		//		if err != nil {
		//			log.Printf("failed to notify master for shard request: %v", err)
		//		}
		//	}
	}
	return nil
}

// ShardFinishRequest notify master server after finishing shard
func (s *TabletServiceServer) notifyMasterShardFinished(ctx context.Context, tableName string, targetTabletAddress string, rowFromPre string, rowToPre string, rowFromPost string, rowToPost string) error {
	source := &ipb.TabletRowRange{
		TabletAddress: s.TabletAddress,
		RowFrom:       rowFromPre,
		RowTo:         rowToPre,
	}
	target := &ipb.TabletRowRange{
		TabletAddress: targetTabletAddress,
		RowFrom:       rowFromPost,
		RowTo:         rowToPost,
	}

	req := &ipb.ShardFinishNotificationRequest{
		TableName: tableName,
		Source:    source,
		Target:    target,
	}

	client := *s.MasterClient
	_, err := client.NotifyShardFinish(ctx, req)
	if err != nil {
		log.Printf("failed to notify master for shard finish: %v", err)
		return err
	}
	return nil
}

// UpdateShardRequest
func (s *TabletServiceServer) notifyTabletServerForShardUpdate(tableName string, targetTabletAddress string) error {
	req := &ipb.UpdateShardRequest{
		TableName:           tableName,
		SourceTabletAddress: s.TabletAddress,
	}

	conn, err := grpc.NewClient(targetTabletAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	client := ipb.NewTabletInternalServiceClient(conn)
	response, err := client.UpdateShard(context.Background(), req)
	if err != nil {
		return status.Errorf(codes.Internal, "%v", err)
	}

	if response.Success {
		// TODO: (Test This) shard is done, we need to delete the table path in the original server address
		// IMPORTANT: it has to be deleted after the migration is done, otherwise the migration cannot read the file
		deleteReq := &ipb.DeleteTableInternalRequest{
			TableName: tableName,
		}
		_, err = s.DeleteTable(context.Background(), deleteReq)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to delete table: %v", err)
		}

	}

	return nil
}

// ------------------------ -- Tablet As a grpc service provider ------------------------------
func (s *TabletServiceServer) UpdateShard(ctx context.Context, req *ipb.UpdateShardRequest) (*ipb.UpdateShardResponse, error) {
	tableName := req.TableName
	sourceTabletAddress := req.SourceTabletAddress

	// migrate
	err := s.MigrateTableToSelf(sourceTabletAddress, tableName)
	if err != nil {
		return &ipb.UpdateShardResponse{
			Success: false,
		}, err
	}

	return &ipb.UpdateShardResponse{
		Success: true,
	}, nil
}
