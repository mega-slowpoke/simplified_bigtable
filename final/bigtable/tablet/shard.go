package tablet

import (
	"context"
	ipb "final/proto/internal-api"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"strings"
	"time"
)

// -------------------------- Tablet As a grpc service consumer ------------------------------
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

	logrus.Debugf("tables are chosen to be moved %v", toMoveTables)

	for _, tableName := range toMoveTables {
		req := &ipb.ShardRequest{
			TabletAddress: s.TabletAddress,
			TableName:     tableName,
		}

		response, err := client.NotifyShardRequest(context.Background(), req)
		if err != nil {
			// it means no available tablet to move the shard, stop sharding
			if strings.Contains(err.Error(), "no tablet servers available excluding the specified server") {
				//logrus.Debugf("failed to notify master for shard request because no other tablet servers available")
				return nil
			}
			logrus.Infof("failed to notify master for shard request because of master/tablet failure: %v", err)
		}

		logrus.Debugf("get target tablet, try to ask it to move table %v", response.TargetTabletAddress)
		err = s.notifyTabletServerForShardUpdate(tableName, response.TargetTabletAddress)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to notify tablet server for shard update: %v", err)
		}

		err = s.notifyMasterShardFinished(context.Background(), tableName, response.TargetTabletAddress)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to notify master that shard is done: %v", err)
		}

	}
	return nil
}

// ShardFinishRequest notify master server after finishing shard
func (s *TabletServiceServer) notifyMasterShardFinished(ctx context.Context, tableName string, targetTabletAddress string) error {
	req := &ipb.ShardFinishNotificationRequest{
		TableName: tableName,
		Source:    s.TabletAddress,
		Target:    targetTabletAddress,
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
