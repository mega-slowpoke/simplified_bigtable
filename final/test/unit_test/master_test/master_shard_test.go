// File: test/unit_test/master_test/master_shard_test.go
package master_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"final/bigtable" // Adjust to your project path
	epb "final/proto/external-api"
	ipb "final/proto/internal-api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockTabletInternalServiceServer mocks the TabletInternalServiceServer and TabletExternalServiceServer interfaces for testing.
type MockTabletInternalServiceServer struct {
	ipb.UnimplementedTabletInternalServiceServer
	epb.UnimplementedTabletExternalServiceServer
	MaxTables     int
	CurrentTables int
	ServerAddress string
	MasterClient  ipb.MasterInternalServiceClient
	mu            sync.Mutex
}

// CreateTable handles table creation requests.
func (m *MockTabletInternalServiceServer) CreateTable(ctx context.Context, req *ipb.CreateTableInternalRequest) (*ipb.CreateTableInternalResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.CurrentTables >= m.MaxTables {
		shardReq := &ipb.ShardRequest{
			TabletAddress: m.ServerAddress,
			TableName:     req.TableName,
		}
		_, err := m.MasterClient.NotifyShardRequest(ctx, shardReq)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to notify master for sharding: %v", err)
		}
		return nil, status.Errorf(codes.ResourceExhausted, "tablet server %s has reached maximum table capacity", m.ServerAddress)
	}

	m.CurrentTables++
	return &ipb.CreateTableInternalResponse{Success: true}, nil
}

// DeleteTable handles table deletion requests.
func (m *MockTabletInternalServiceServer) DeleteTable(ctx context.Context, req *ipb.DeleteTableInternalRequest) (*ipb.DeleteTableInternalResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.CurrentTables > 0 {
		m.CurrentTables--
	}
	return &ipb.DeleteTableInternalResponse{Success: true}, nil
}

// Read simulates a read operation.
func (m *MockTabletInternalServiceServer) Read(ctx context.Context, req *epb.ReadRequest) (*epb.ReadResponse, error) {
	return &epb.ReadResponse{
		Values: []*epb.ValueWithTimestamps{
			{
				Timestamp: time.Now().Unix(),
				Value:     "mock_value",
			},
		},
	}, nil
}

// Write simulates a write operation.
func (m *MockTabletInternalServiceServer) Write(ctx context.Context, req *epb.WriteRequest) (*epb.WriteResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.CurrentTables >= m.MaxTables {
		return &epb.WriteResponse{
			Success:      false,
			ErrorMessage: "Tablet server at capacity",
		}, nil
	}

	m.CurrentTables++
	return &epb.WriteResponse{
		Success:      true,
		ErrorMessage: "",
	}, nil
}

// Delete handles delete operations.
func (m *MockTabletInternalServiceServer) Delete(ctx context.Context, req *epb.DeleteRequest) (*epb.DeleteResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.CurrentTables > 0 {
		m.CurrentTables--
	}
	return &epb.DeleteResponse{
		Success:      true,
		ErrorMessage: "",
	}, nil
}

// startMockTabletServer starts a mock Tablet Server and returns the gRPC server.
func startMockTabletServer(address string, server *MockTabletInternalServiceServer) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	ipb.RegisterTabletInternalServiceServer(grpcServer, server)
	epb.RegisterTabletExternalServiceServer(grpcServer, server)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Mock tablet server at %s stopped: %v", address, err)
		}
	}()

	time.Sleep(100 * time.Millisecond) // Ensure server starts
	return grpcServer, nil
}

// startMasterServer starts the MasterServer's gRPC server and returns the server.
func startMasterServer(master *bigtable.MasterServer, address string) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	epb.RegisterMasterExternalServiceServer(grpcServer, master)
	ipb.RegisterMasterInternalServiceServer(grpcServer, master)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Master server stopped: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond) // Ensure server starts
	return grpcServer, nil
}

// TestRegisterTablet verifies that a Tablet Server can register successfully.
func TestRegisterTablet(t *testing.T) {
	master := bigtable.NewMasterServer()
	masterAddress := "localhost:7000"
	grpcMasterServer, err := startMasterServer(master, masterAddress)
	if err != nil {
		t.Fatalf("Failed to start MasterServer: %v", err)
	}
	defer grpcMasterServer.Stop()

	mockServer := &MockTabletInternalServiceServer{
		MaxTables:     2,
		ServerAddress: "localhost:6001",
	}
	grpcMockServer, err := startMockTabletServer(mockServer.ServerAddress, mockServer)
	if err != nil {
		t.Fatalf("Failed to start mock tablet server: %v", err)
	}
	defer grpcMockServer.Stop()

	connMasterInternal, err := grpc.Dial(masterAddress, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		t.Fatalf("Failed to dial MasterInternalService: %v", err)
	}
	defer connMasterInternal.Close()
	mockServer.MasterClient = ipb.NewMasterInternalServiceClient(connMasterInternal)

	req := &ipb.RegisterTabletRequest{TabletAddress: mockServer.ServerAddress}
	_, err = master.RegisterTablet(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to register tablet server: %v", err)
	}

	if _, exists := master.State.TabletServers[mockServer.ServerAddress]; !exists {
		t.Fatalf("Tablet server '%s' not found in Master state", mockServer.ServerAddress)
	}
}

// TestCreateTable verifies table creation and assignment to the least loaded Tablet Server.
func TestCreateTable(t *testing.T) {
	master := bigtable.NewMasterServer()
	masterAddress := "localhost:7001"
	grpcMasterServer, err := startMasterServer(master, masterAddress)
	if err != nil {
		t.Fatalf("Failed to start MasterServer: %v", err)
	}
	defer grpcMasterServer.Stop()

	mockServer1 := &MockTabletInternalServiceServer{
		MaxTables:     2,
		ServerAddress: "localhost:6002",
	}
	grpcServer1, err := startMockTabletServer(mockServer1.ServerAddress, mockServer1)
	if err != nil {
		t.Fatalf("Failed to start mock tablet server1: %v", err)
	}
	defer grpcServer1.Stop()

	mockServer2 := &MockTabletInternalServiceServer{
		MaxTables:     2,
		ServerAddress: "localhost:6003",
	}
	grpcServer2, err := startMockTabletServer(mockServer2.ServerAddress, mockServer2)
	if err != nil {
		t.Fatalf("Failed to start mock tablet server2: %v", err)
	}
	defer grpcServer2.Stop()

	connMasterInternal, err := grpc.Dial(masterAddress, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		t.Fatalf("Failed to dial MasterInternalService: %v", err)
	}
	defer connMasterInternal.Close()
	mockServer1.MasterClient = ipb.NewMasterInternalServiceClient(connMasterInternal)
	mockServer2.MasterClient = ipb.NewMasterInternalServiceClient(connMasterInternal)

	_, err = master.RegisterTablet(context.Background(), &ipb.RegisterTabletRequest{TabletAddress: mockServer1.ServerAddress})
	if err != nil {
		t.Fatalf("Failed to register tablet server1: %v", err)
	}
	_, err = master.RegisterTablet(context.Background(), &ipb.RegisterTabletRequest{TabletAddress: mockServer2.ServerAddress})
	if err != nil {
		t.Fatalf("Failed to register tablet server2: %v", err)
	}

	createReq := &epb.CreateTableRequest{
		TableName: "test_table",
		ColumnFamilies: []*epb.CreateTableRequest_ColumnFamily{
			{
				FamilyName: "cf1",
				Columns:    []string{"col1", "col2"},
			},
		},
	}
	resp, err := master.CreateTable(context.Background(), createReq)
	if err != nil {
		t.Fatalf("CreateTable RPC failed: %v", err)
	}
	if !resp.Success {
		t.Fatalf("CreateTable failed: %s", resp.Message)
	}

	table, exists := master.State.Tables["test_table"]
	if !exists {
		t.Fatalf("Table 'test_table' not found in Master state")
	}
	if len(table.Tablets) != 1 {
		t.Fatalf("Expected 1 tablet for 'test_table', got %d", len(table.Tablets))
	}
	assignedServer := table.Tablets[0].TabletServer
	if assignedServer != mockServer1.ServerAddress && assignedServer != mockServer2.ServerAddress {
		t.Fatalf("Table assigned to unknown server: %s", assignedServer)
	}
}

// TestShardRequest_Success tests the sharding process when a Tablet Server exceeds its capacity.
func TestShardRequest_Success(t *testing.T) {
	master := bigtable.NewMasterServer()
	masterAddress := "localhost:7002"
	grpcMasterServer, err := startMasterServer(master, masterAddress)
	if err != nil {
		t.Fatalf("Failed to start MasterServer: %v", err)
	}
	defer grpcMasterServer.Stop()

	mockServer1 := &MockTabletInternalServiceServer{
		MaxTables:     1,
		ServerAddress: "localhost:6004",
	}
	grpcServer1, err := startMockTabletServer(mockServer1.ServerAddress, mockServer1)
	if err != nil {
		t.Fatalf("Failed to start mock tablet server1: %v", err)
	}
	defer grpcServer1.Stop()

	mockServer2 := &MockTabletInternalServiceServer{
		MaxTables:     1,
		ServerAddress: "localhost:6005",
	}
	grpcServer2, err := startMockTabletServer(mockServer2.ServerAddress, mockServer2)
	if err != nil {
		t.Fatalf("Failed to start mock tablet server2: %v", err)
	}
	defer grpcServer2.Stop()

	connMasterInternal, err := grpc.Dial(masterAddress, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		t.Fatalf("Failed to dial MasterInternalService: %v", err)
	}
	defer connMasterInternal.Close()
	mockServer1.MasterClient = ipb.NewMasterInternalServiceClient(connMasterInternal)
	mockServer2.MasterClient = ipb.NewMasterInternalServiceClient(connMasterInternal)

	_, err = master.RegisterTablet(context.Background(), &ipb.RegisterTabletRequest{TabletAddress: mockServer1.ServerAddress})
	if err != nil {
		t.Fatalf("Failed to register tablet server1: %v", err)
	}
	_, err = master.RegisterTablet(context.Background(), &ipb.RegisterTabletRequest{TabletAddress: mockServer2.ServerAddress})
	if err != nil {
		t.Fatalf("Failed to register tablet server2: %v", err)
	}

	createReq1 := &epb.CreateTableRequest{
		TableName: "table_shard1",
		ColumnFamilies: []*epb.CreateTableRequest_ColumnFamily{
			{
				FamilyName: "cf1",
				Columns:    []string{"col1"},
			},
		},
	}
	resp1, err := master.CreateTable(context.Background(), createReq1)
	if err != nil {
		t.Fatalf("CreateTable RPC1 failed: %v", err)
	}
	if !resp1.Success {
		t.Fatalf("CreateTable1 failed: %s", resp1.Message)
	}

	createReq2 := &epb.CreateTableRequest{
		TableName: "table_shard2",
		ColumnFamilies: []*epb.CreateTableRequest_ColumnFamily{
			{
				FamilyName: "cf2",
				Columns:    []string{"col2"},
			},
		},
	}
	resp2, err := master.CreateTable(context.Background(), createReq2)
	if err != nil && status.Code(err) != codes.ResourceExhausted {
		t.Fatalf("CreateTable RPC2 failed: %v", err)
	}
	if err == nil && !resp2.Success {
		t.Fatalf("CreateTable2 failed: %s", resp2.Message)
	}

	fmt.Println("Table creation response:", resp2)

	time.Sleep(100 * time.Millisecond) // Wait for sharding

	table2, exists := master.State.Tables["table_shard2"]
	if !exists {
		t.Fatalf("Table 'table_shard2' not found in Master state")
	}
	if len(table2.Tablets) != 1 {
		t.Fatalf("Expected 1 tablet for 'table_shard2', got %d", len(table2.Tablets))
	}
	if table2.Tablets[0].TabletServer != mockServer2.ServerAddress {
		t.Fatalf("Expected 'table_shard2' to be assigned to '%s', got '%s'", mockServer2.ServerAddress, table2.Tablets[0].TabletServer)
	}
}

// TestCreateTable_NoAvailableServers tests table creation when no Tablet Servers are available.
func TestCreateTable_NoAvailableServers(t *testing.T) {
	master := bigtable.NewMasterServer()
	masterAddress := "localhost:7004"
	grpcMasterServer, err := startMasterServer(master, masterAddress)
	if err != nil {
		t.Fatalf("Failed to start MasterServer: %v", err)
	}
	defer grpcMasterServer.Stop()

	mockServer := &MockTabletInternalServiceServer{
		MaxTables:     1,
		ServerAddress: "localhost:6008",
	}
	grpcMockServer, err := startMockTabletServer(mockServer.ServerAddress, mockServer)
	if err != nil {
		t.Fatalf("Failed to start mock tablet server: %v", err)
	}
	defer grpcMockServer.Stop()

	connMasterInternal, err := grpc.Dial(masterAddress, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		t.Fatalf("Failed to dial MasterInternalService: %v", err)
	}
	defer connMasterInternal.Close()
	mockServer.MasterClient = ipb.NewMasterInternalServiceClient(connMasterInternal)

	_, err = master.RegisterTablet(context.Background(), &ipb.RegisterTabletRequest{TabletAddress: mockServer.ServerAddress})
	if err != nil {
		t.Fatalf("Failed to register tablet server: %v", err)
	}

	createReq1 := &epb.CreateTableRequest{
		TableName: "table_no_shard1",
		ColumnFamilies: []*epb.CreateTableRequest_ColumnFamily{
			{
				FamilyName: "cf1",
				Columns:    []string{"col1"},
			},
		},
	}
	resp1, err := master.CreateTable(context.Background(), createReq1)
	if err != nil {
		t.Fatalf("CreateTable RPC1 failed: %v", err)
	}
	if !resp1.Success {
		t.Fatalf("CreateTable1 failed: %s", resp1.Message)
	}

	createReq2 := &epb.CreateTableRequest{
		TableName: "table_no_shard2",
		ColumnFamilies: []*epb.CreateTableRequest_ColumnFamily{
			{
				FamilyName: "cf2",
				Columns:    []string{"col2"},
			},
		},
	}
	resp2, err := master.CreateTable(context.Background(), createReq2)
	if err != nil {
		t.Fatalf("CreateTable RPC2 failed: %v", err)
	}
	if resp2.Success {
		t.Fatalf("Expected CreateTable2 to fail due to no available servers, but it succeeded")
	}
}
