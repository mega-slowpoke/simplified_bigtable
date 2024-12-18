package bigtable

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "final/proto/external-api"
)

// TabletRange represents a range of row keys handled by a specific tablet server.
type TabletRange struct {
	StartRow      string
	EndRow        string
	ServerAddress string
}

// TableCache maintains the mapping of row key ranges to tablet servers for a specific table.
type TableCache struct {
	ranges []TabletRange
	mu     sync.RWMutex
}

// AddRange adds a new TabletRange to the TableCache.
func (tc *TableCache) AddRange(tr TabletRange) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.ranges = append(tc.ranges, tr)
	// Sort the ranges by StartRow for efficient searching
	sort.Slice(tc.ranges, func(i, j int) bool {
		return strings.Compare(tc.ranges[i].StartRow, tc.ranges[j].StartRow) < 0
	})
}

// FindServerByRowKey searches for the tablet server responsible for the given rowKey.
// It returns the server address and a boolean indicating whether the rowKey was found in the cache.
func (tc *TableCache) FindServerByRowKey(rowKey string) (string, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	// Binary search since ranges are sorted
	index := sort.Search(len(tc.ranges), func(i int) bool {
		return strings.Compare(tc.ranges[i].StartRow, rowKey) > 0
	})

	if index > 0 {
		tr := tc.ranges[index-1]
		// Check if rowKey falls within this range
		if (tr.StartRow == "" || strings.Compare(rowKey, tr.StartRow) >= 0) &&
			(tr.EndRow == "" || strings.Compare(rowKey, tr.EndRow) < 0) {
			return tr.ServerAddress, true
		}
	}

	return "", false
}

type Client struct {
	masterConn    *grpc.ClientConn
	masterClient  pb.MasterExternalServiceClient
	tabletConns   map[string]*grpc.ClientConn
	tabletClients map[string]pb.TabletExternalServiceClient
	caches        map[string]*TableCache // map of table name to TableCache
	mu            sync.RWMutex           // Protects tabletConns, tabletClients, and caches
}

func NewClient(masterAddress string) (*Client, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(masterAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master server: %v", err)
	}

	masterClient := pb.NewMasterExternalServiceClient(conn)

	return &Client{
		masterConn:    conn,
		masterClient:  masterClient,
		tabletConns:   make(map[string]*grpc.ClientConn),
		tabletClients: make(map[string]pb.TabletExternalServiceClient),
		caches:        make(map[string]*TableCache),
	}, nil
}

func (c *Client) Close() {
	if c.masterConn != nil {
		c.masterConn.Close()
		log.Println("Client closed connection to master server.")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for addr, conn := range c.tabletConns {
		if conn != nil {
			conn.Close()
			log.Printf("Client closed connection to tablet server: %s", addr)
		}
	}
}

// CreateTable creates a new table in the master server with specified column families and columns.
func (c *Client) CreateTable(tableName string, columnFamilies map[string][]string) error {
	// Convert columnFamilies map to repeated ColumnFamily message
	var cfMsgs []*pb.CreateTableRequest_ColumnFamily
	for family, columns := range columnFamilies {
		cfMsgs = append(cfMsgs, &pb.CreateTableRequest_ColumnFamily{
			FamilyName: family,
			Columns:    columns,
		})
	}

	req := &pb.CreateTableRequest{
		TableName:      tableName,
		ColumnFamilies: cfMsgs,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.masterClient.CreateTable(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateTable RPC failed: %v", err)
	}

	if !resp.Success {
		return fmt.Errorf("CreateTable failed: %s", resp.Message)
	}

	log.Printf("Table '%s' created successfully.", tableName)
	return nil
}

// DeleteTable deletes an existing table from the master server.
func (c *Client) DeleteTable(tableName string) error {
	req := &pb.DeleteTableRequest{
		TableName: tableName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.masterClient.DeleteTable(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteTable RPC failed: %v", err)
	}

	if !resp.Success {
		return fmt.Errorf("DeleteTable failed: %s", resp.Message)
	}

	log.Printf("Table '%s' deleted successfully.", tableName)
	return nil
}

// GetTabletLocation retrieves the tablet server address and row range for a given key.
func (c *Client) GetTabletLocation(tableName, rowKey string) (string, string, string, error) {
	// Check if the table cache exists
	c.mu.RLock()
	tableCache, exists := c.caches[tableName]
	c.mu.RUnlock()

	if exists {
		// Attempt to find the server in the cache
		serverAddress, found := tableCache.FindServerByRowKey(rowKey)
		if found {
			return serverAddress, "", "", nil
		}
	} else {
		// Initialize a new cache for the table
		c.mu.Lock()
		tableCache, exists = c.caches[tableName]
		if !exists {
			tableCache = &TableCache{}
			c.caches[tableName] = tableCache
		}
		c.mu.Unlock()
	}

	// If not found in cache, perform GetTabletLocation RPC
	req := &pb.GetTabletLocationRequest{
		TableName: tableName,
		Key:       rowKey,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.masterClient.GetTabletLocation(ctx, req)
	if err != nil {
		return "", "", "", fmt.Errorf("GetTabletLocation RPC failed: %v", err)
	}

	// Handle the response based on the presence of tablet_server_address
	if resp.TabletServerAddress == "" {
		// Table does not exist
		return "", "", "", errors.New("table does not exist")
	}

	// Check if the tablet is not sharded (complete table on one server)
	if resp.TabletStartRow == "" && resp.TabletEndRow == "" {
		// Cache the entire table range
		tr := TabletRange{
			StartRow:      "",
			EndRow:        "",
			ServerAddress: resp.TabletServerAddress,
		}
		tableCache.AddRange(tr)
		return resp.TabletServerAddress, "", "", nil
	}

	// Cache the specific range returned
	tr := TabletRange{
		StartRow:      resp.TabletStartRow,
		EndRow:        resp.TabletEndRow,
		ServerAddress: resp.TabletServerAddress,
	}
	tableCache.AddRange(tr)

	// Verify if the rowKey falls within the returned range
	if (tr.StartRow == "" || strings.Compare(rowKey, tr.StartRow) >= 0) &&
		(tr.EndRow == "" || strings.Compare(rowKey, tr.EndRow) < 0) {
		return tr.ServerAddress, tr.StartRow, tr.EndRow, nil
	}

	// If the rowKey does not fall within the range, it does not exist
	return "", "", "", fmt.Errorf("rowKey '%s' does not exist in table '%s'", rowKey, tableName)
}

// getTabletClient retrieves or establishes a connection to a tablet server.
func (c *Client) getTabletClient(tabletAddress string) (pb.TabletExternalServiceClient, error) {
	c.mu.RLock()
	client, exists := c.tabletClients[tabletAddress]
	c.mu.RUnlock()

	if exists {
		return client, nil
	}

	// Establish a new connection
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(tabletAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tablet server %s: %v", tabletAddress, err)
	}

	client = pb.NewTabletExternalServiceClient(conn)

	c.mu.Lock()
	c.tabletConns[tabletAddress] = conn
	c.tabletClients[tabletAddress] = client
	c.mu.Unlock()

	log.Printf("Connected to tablet server: %s", tabletAddress)
	return client, nil
}

// Read retrieves data from a tablet server for a specific row key and column.
func (c *Client) Read(tableName, rowKey, columnFamily, columnQualifier string, returnVersion int64) ([]*pb.ValueWithTimestamps, error) {
	// Attempt to get the server address from cache
	serverAddress, _, _, err := c.GetTabletLocation(tableName, rowKey)
	if err != nil {
		return nil, err
	}

	client, err := c.getTabletClient(serverAddress)
	if err != nil {
		return nil, err
	}

	req := &pb.ReadRequest{
		TableName:       tableName,
		RowKey:          rowKey,
		ColumnFamily:    columnFamily,
		ColumnQualifier: columnQualifier,
		ReturnVersion:   returnVersion,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Read(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Read RPC failed: %v", err)
	}

	// If no values are returned, the rowKey does not exist
	if len(resp.Values) == 0 {
		return nil, fmt.Errorf("rowKey '%s' does not exist in table '%s'", rowKey, tableName)
	}

	return resp.Values, nil
}

// Write inserts or updates data in a tablet server for a specific row key and column.
func (c *Client) Write(tableName, rowKey, columnFamily, columnQualifier string, value []byte, timestamp int64) error {
	// Attempt to get the server address from cache
	serverAddress, _, _, err := c.GetTabletLocation(tableName, rowKey)
	if err != nil {
		return err
	}

	client, err := c.getTabletClient(serverAddress)
	if err != nil {
		return err
	}

	req := &pb.WriteRequest{
		TableName:       tableName,
		RowKey:          rowKey,
		ColumnFamily:    columnFamily,
		ColumnQualifier: columnQualifier,
		Value:           value,
		Timestamp:       timestamp,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Write(ctx, req)
	if err != nil {
		return fmt.Errorf("Write RPC failed: %v", err)
	}

	if !resp.Success {
		return fmt.Errorf("Write failed: %s", resp.ErrorMessage)
	}

	log.Printf("Write successful for row key '%s'.", rowKey)
	return nil
}

// Delete removes data from a tablet server for a specific row key and column.
func (c *Client) Delete(tableName, rowKey, columnFamily, columnQualifier string) error {
	// Attempt to get the server address from cache
	serverAddress, _, _, err := c.GetTabletLocation(tableName, rowKey)
	if err != nil {
		return err
	}

	client, err := c.getTabletClient(serverAddress)
	if err != nil {
		return err
	}

	req := &pb.DeleteRequest{
		TableName:       tableName,
		RowKey:          rowKey,
		ColumnFamily:    columnFamily,
		ColumnQualifier: columnQualifier,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("Delete RPC failed: %v", err)
	}

	if !resp.Success {
		return fmt.Errorf("Delete failed: %s", resp.ErrorMessage)
	}

	log.Printf("Delete successful for row key '%s'.", rowKey)
	return nil
}
