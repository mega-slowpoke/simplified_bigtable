package tablet

import (
	"context"
	ipb "final/proto/internal-api"
	"log"
	"sort"
)

var (
	tablesColumns map[string]map[string][]string
	tablesInfo    map[string]string
	tableList     struct {
		Tables []string
	}
	shardMaxSize int
)

//-------------------------- Tablet As a grpc service consumer ------------------------------

// checkAndNotifyMasterForShard if the size of any table rows is greater than max size, notify the master to issue a shard command
// Problem: what if this table is being modified during this time
func (s *TabletServiceServer) checkAndNotifyMasterForShard() {
	client := *(s.MasterClient)
	tablesRows := s.TablesRows
	for tableName := range s.TablesRows {
		if len(tablesRows[tableName]) >= s.MaxShardSize {
			rowKeys := make([]string, len(tablesRows[tableName]))
			copy(rowKeys, tablesRows[tableName])
			sort.Strings(rowKeys) // TODO: check the sort? sort by what

			shardRowKeys := rowKeys[:s.MaxShardSize/2+1]
			remainingRowKeys := rowKeys[s.MaxShardSize/2:]

			req := &ipb.ShardRequest{
				TabletAddress:    s.TabletAddress,
				TableName:        tableName,
				ShardRowKeys:     shardRowKeys,
				RemainingRowKeys: remainingRowKeys,
			}

			_, err := client.NotifyShardRequest(context.Background(), req)
			if err != nil {
				log.Printf("failed to notify master for shard request: %v", err)
			}
		}
	}
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

//-------------------------- Tablet As a grpc service provider ------------------------------
