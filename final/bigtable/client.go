package bigtable

// TODO: tablet_start_row, tablet_end_row are not used

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "final/proto/external-api"
)

type BigTableClient struct {
	metadataAddr   string
	metadataClient pb.MetadataServiceClient
	tabletCache    sync.Map
}

func NewBigTableClient(metadataAddr string) (*BigTableClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(metadataAddr, opts...)
	if err != nil {
		return nil, err
	}

	metadataClient := pb.NewMetadataServiceClient(conn)
	return &BigTableClient{
		metadataAddr:   metadataAddr,
		metadataClient: metadataClient,
	}, nil
}

func (c *BigTableClient) getTabletLocation(tableName, rowKey string) (*pb.TabletLocationResponse, error) {
	cacheKey := fmt.Sprintf("%s:%s", tableName, rowKey)
	if val, ok := c.tabletCache.Load(cacheKey); ok {
		return val.(*pb.TabletLocationResponse), nil
	}

	req := &pb.TabletLocationRequest{
		TableName: tableName,
		RowKey:    rowKey,
	}

	resp, err := c.metadataClient.GetTabletLocation(context.Background(), req)
	if err != nil {
		return nil, err
	}

	c.tabletCache.Store(cacheKey, resp)
	return resp, nil
}

func (c *BigTableClient) Read(tableName, rowKey, columnFamily, columnQualifier string) (*pb.ReadResponse, error) {
	tabletLocation, err := c.getTabletLocation(tableName, rowKey)
	if err != nil {
		return nil, err
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(tabletLocation.ServerAddress, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	tabletClient := pb.NewTabletServiceClient(conn)

	req := &pb.ReadRequest{
		TableName:       tableName,
		RowKey:          rowKey,
		ColumnFamily:    columnFamily,
		ColumnQualifier: columnQualifier,
	}

	resp, err := tabletClient.Read(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *BigTableClient) Write(tableName, rowKey, columnFamily, columnQualifier string, value []byte, timestamp int64) (*pb.WriteResponse, error) {
	tabletLocation, err := c.getTabletLocation(tableName, rowKey)
	if err != nil {
		return nil, err
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(tabletLocation.ServerAddress, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	tabletClient := pb.NewTabletServiceClient(conn)

	req := &pb.WriteRequest{
		TableName:       tableName,
		RowKey:          rowKey,
		ColumnFamily:    columnFamily,
		ColumnQualifier: columnQualifier,
		Value:           value,
		Timestamp:       timestamp,
	}

	resp, err := tabletClient.Write(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *BigTableClient) Delete(tableName, rowKey, columnFamily, columnQualifier string) (*pb.DeleteResponse, error) {
	tabletLocation, err := c.getTabletLocation(tableName, rowKey)
	if err != nil {
		return nil, err
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(tabletLocation.ServerAddress, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	tabletClient := pb.NewTabletServiceClient(conn)

	req := &pb.DeleteRequest{
		TableName:       tableName,
		RowKey:          rowKey,
		ColumnFamily:    columnFamily,
		ColumnQualifier: columnQualifier,
	}

	resp, err := tabletClient.Delete(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *BigTableClient) Scan(tableName, startRowKey, endRowKey, columnFamily, columnQualifier string) error {
	currentRowKey := startRowKey

	for {
		tabletLocation, err := c.getTabletLocation(tableName, currentRowKey)
		if err != nil {
			return err
		}

		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

		conn, err := grpc.NewClient(tabletLocation.ServerAddress, opts...)
		if err != nil {
			return err
		}
		defer conn.Close()

		tabletClient := pb.NewTabletServiceClient(conn)

		tabletEndRow := tabletLocation.TabletEndRow
		var scanEndKey string
		if tabletEndRow == "" || tabletEndRow > endRowKey {
			scanEndKey = endRowKey
		} else {
			scanEndKey = tabletEndRow
		}

		req := &pb.ScanRequest{
			TableName:       tableName,
			StartRowKey:     currentRowKey,
			EndRowKey:       scanEndKey,
			ColumnFamily:    columnFamily,
			ColumnQualifier: columnQualifier,
		}

		stream, err := tabletClient.Scan(context.Background(), req)
		if err != nil {
			return err
		}

		for {
			resp, err := stream.Recv()
			if err != nil {
				break
			}
			fmt.Printf("RowKey: %s, Value: %s\n", resp.RowKey, string(resp.Value))
		}

		if scanEndKey == endRowKey {
			break
		} else {
			currentRowKey = scanEndKey
		}
	}

	return nil
}
