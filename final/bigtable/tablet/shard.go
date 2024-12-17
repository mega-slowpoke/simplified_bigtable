package tablet

import (
	"context"
	ipb "final/proto/internal-api"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

//-------------------------- Tablet As a grpc service consumer ------------------------------

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

		err = s.notifyTabletServerForShardUpdate(tableName, []string{}, response.TargetTabletAddress)
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
func (s *TabletServiceServer) notifyTabletServerForShardUpdate(tableName string, shardRowKeySet []string, targetTabletAddress string) error {
	tableColumns := make(map[string]*ipb.Columns)
	for columnFamily, columns := range s.TablesColumns[tableName] {
		tableColumns[columnFamily] = &ipb.Columns{ColumnNames: columns}
	}

	req := &ipb.UpdateShardRequest{
		TableName: tableName,
		//TableRows:    shardRowKeySet,
		TableColumns: tableColumns,
		//TableInfo:    s.TablesInfo,
	}

	//
	conn, err := grpc.NewClient(targetTabletAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	client := ipb.NewTabletInternalServiceClient(conn)
	_, err = client.UpdateShard(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to notify tablet server for shard update: %v", err)
	}
	return nil
}

// ------------------------ -- Tablet As a grpc service provider ------------------------------
func (s *TabletServiceServer) UpdateShard(ctx context.Context, req *ipb.UpdateShardRequest) (*ipb.UpdateShardResponse, error) {
	tableName := req.TableName
	tableRows := req.TableRows
	tableColumns := req.TableColumns
	tableInfo := req.TableInfo

	// move metadata
	// add tableName to table list
	if !contains(s.TableList, tableName) {
		s.TableList = append(s.TableList, tableName)
	}

	// add row Keys to the tableRows
	s.TablesRows[tableName] = []string{}
	for _, rowKey := range tableRows {
		s.TablesRows[tableName] = append(s.TablesRows[tableName], rowKey)
	}

	for columnFamily, columns := range tableColumns {
		s.TablesColumns[tableName][columnFamily] = columns.ColumnNames
	}

	s.TablesInfo[tableName] = tableInfo

	// TODO:  metadata_for_col_index(columnsJSON)

	// move storage data

	return &ipb.UpdateShardResponse{}, nil
}
