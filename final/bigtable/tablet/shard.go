package tablet

//import (
//	"context"
//	ipb "final/proto/internal-api"
//	"fmt"
//	"google.golang.org/grpc"
//	"sort"
//)
//
//// 模拟全局变量
//var (
//	tablesRows    map[string][]string
//	tablesColumns map[string]map[string][]string
//	tablesInfo    map[string]string
//	tableList     struct {
//		Tables []string
//	}
//	tablesMaxSize int
//	shardMaxSize  int
//)
//
//// CheckShardRequest实现
//func (s *TabletServiceServer) checkAndNotifyMasterForShard() {
//	client := *(s.MasterClient)
//
//	for tableName := range tables_rows {
//		if len(tables_rows[tableName]) >= shard_max_size {
//			rowkeys := make([]string, len(tables_rows[tableName]))
//			copy(rowkeys, tables_rows[tableName])
//			sort.Strings(rowkeys)
//
//			shardRow := rowkeys[:shard_max_size/2+1]
//			originalRow := rowkeys[shard_max_size/2:]
//
//			req := &shard_service.ShardRequest{
//				TabletServerName: TABLET_SERVER_NAME,
//				TableName:        tableName,
//				ShardRow:         shardRow,
//				OriginalRow:      originalRow,
//			}
//
//			_, err := client.NotifyShardRequest(ctx, req)
//			if err!= nil {
//				return fmt.Errorf("failed to notify master for shard request: %v", err)
//			}
//		}
//	}
//	return nil
//}
//
//}
//
//func (s *TabletServiceServer) CheckShardRequest(ctx context.Context, table string) error {
//	client := *(s.MasterClient)
//	if len(tablesRows[table]) >= shardMaxSize {
//		rowkeys := make([]string, len(tablesRows[table]))
//		copy(rowkeys, tablesRows[table])
//		sort.Strings(rowkeys)
//		shardRow := rowkeys[:shardMaxSize/2+1]
//		originalRow := rowkeys[shardMaxSize/2:]
//
//		req := &ipb.ShardRequest{
//			TabletAddress: s.TabletAddress,
//			Table:         table,
//			ShardRow:      shardRow,
//			OriginalRow:   originalRow,
//		}
//		_, err := client.CheckShardRequest(ctx, req)
//		return err
//	}
//	return nil
//}
//
//// UpdateShardRequest实现
//func (s *TabletServiceServer) UpdateShardRequest(ctx context.Context, table string, shardRowKeySet []string) error {
//	client := *(s.MasterClient)
//	tableRows := map[string][]string{
//		table: shardRowKeySet,
//	}
//	tableColumns := tablesColumns[table]
//	tableInfo := tablesInfo[table]
//
//	req := &ipb.UpdateShardRequest{
//		Table:        table,
//		TableRows:    tableRows,
//		TableColumns: tableColumns,
//		TableInfo:    tableInfo,
//	}
//	_, err := client.UpdateShard(ctx, req)
//	return err
//}
//
//// UpdateShardServer实现
//func (s *TabletServiceServer) UpdateShardServer(ctx context.Context, req *pb.UpdateShardRequest) error {
//	table := req.Table
//	tableRows := req.TableRows
//	tableColumns := req.TableColumns
//	tableInfo := req.TableInfo
//
//	if !contains(tableList.Tables, table) {
//		tableList.Tables = append(tableList.Tables, table)
//	}
//
//	tablesRows[table] = []string{}
//	for rowKey := range tableRows[table] {
//		// metadataForRowIndex(table, rowKey) 这里需要你根据实际功能补充具体逻辑，以下为占位
//		tablesRows[table] = append(tablesRows[table], rowKey)
//	}
//
//	tablesColumns[table] = tableColumns[table]
//	// columnsJson相关逻辑，这里需要你根据实际Python代码中的逻辑补充完整，以下为示意
//	// 假设构建类似的结构并处理元数据等
//	// columnsJson := buildColumnsJson(table, tablesColumns[table])
//	// metadataForColIndex(columnsJson)
//
//	tablesInfo[table] = tableInfo[table]
//
//	return nil
//}
//
//func (s *TabletServiceServer) ShardFinishRequest(ctx context.Context, conn *grpc.ClientConn, table string, originRowKeySet []string, shardRow []string) error {
//	client := pb.NewShardServiceClient(conn)
//	data := []*pb.ShardServerInfo{
//		{
//			Hostname: tabletServerName,
//			Port:     "", // 这里假设你根据实际情况获取端口，目前先留空
//			RowFrom:  originRowKeySet[0],
//			RowTo:    originRowKeySet[len(originRowKeySet)-1],
//		},
//		{
//			Hostname: "", // 同样需按实际情况补充，先留空
//			Port:     "",
//			RowFrom:  shardRow[0],
//			RowTo:    shardRow[len(shardRow)-1],
//		},
//	}
//	req := &pb.ShardFinishRequest{
//		Table: table,
//		Data:  data,
//	}
//	_, err := client.ShardFinish(ctx, req)
//	return err
//}
//
//func contains(s []string, e string) bool {
//	for _, a := range s {
//		if a == e {
//			return true
//		}
//	}
//	return false
//}
