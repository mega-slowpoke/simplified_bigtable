package bigtable

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	epb "final/proto/external-api"
	ipb "final/proto/internal-api"
)

type Tablet struct {
	TableName    string
	StartRow     string
	EndRow       string
	TabletServer string
	Sharded      bool // Indicates if the tablet has been sharded
}

type TabletServerInfo struct {
	Address        string
	LastHeartbeat  time.Time
	RegisteredTime time.Time
	TabletCount    int // Number of tablets assigned to this server
	MaxTableCount  int
}

type Table struct {
	Name           string
	ColumnFamilies map[string][]string // map of family name to columns
	Tablets        []*Tablet
}

type MasterState struct {
	mu            sync.RWMutex
	Tables        map[string]*Table            // map of table name to Table
	TabletServers map[string]*TabletServerInfo // map of tablet server address to info
}

type MasterServer struct {
	epb.UnimplementedMasterExternalServiceServer
	ipb.UnimplementedMasterInternalServiceServer

	State *MasterState
}

func NewMasterServer() *MasterServer {
	return &MasterServer{
		State: &MasterState{
			Tables:        make(map[string]*Table),
			TabletServers: make(map[string]*TabletServerInfo),
		},
	}
}

// Make starts the master server on the specified address.
// It returns a function to stop the server.
func Make(address string) (func(), error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %v", address, err)
	}

	grpcServer := grpc.NewServer()

	masterServer := NewMasterServer()

	// Register external and internal services
	epb.RegisterMasterExternalServiceServer(grpcServer, masterServer)
	ipb.RegisterMasterInternalServiceServer(grpcServer, masterServer)

	// Start heartbeat monitoring in a separate goroutine
	go masterServer.MonitorHeartbeats(10*time.Second, 5*time.Second)

	log.Printf("Master server is running on %s...", address)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// Return a function to stop the server gracefully
	stopFunc := func() {
		grpcServer.GracefulStop()
		log.Println("Master server has been stopped.")
	}

	return stopFunc, nil
}

// ExternalMasterService Implementation
func (ms *MasterServer) CreateTable(ctx context.Context, req *epb.CreateTableRequest) (*epb.CreateTableResponse, error) {
	ms.State.mu.Lock()
	defer ms.State.mu.Unlock()

	if _, exists := ms.State.Tables[req.TableName]; exists {
		msg := fmt.Sprintf("Table '%s' already exists.", req.TableName)
		log.Println(msg)
		return &epb.CreateTableResponse{
			Success: false,
			Message: msg,
		}, nil
	}

	// Initialize the table with column families
	columnFamilies := make(map[string][]string)
	for _, cf := range req.ColumnFamilies {
		columnFamilies[cf.FamilyName] = cf.Columns
	}

	table := &Table{
		Name:           req.TableName,
		ColumnFamilies: columnFamilies,
		Tablets:        []*Tablet{},
	}

	ms.State.Tables[req.TableName] = table

	// Assign the entire table to the least loaded tablet server if any are registered
	if len(ms.State.TabletServers) > 0 {
		assignedServer, err := ms.getLeastLoadedTabletServer()
		if err != nil {
			// Rollback table creation
			delete(ms.State.Tables, req.TableName)
			msg := fmt.Sprintf("Failed to find a tablet server to assign table '%s': %v", req.TableName, err)
			log.Println(msg)
			return &epb.CreateTableResponse{
				Success: false,
				Message: msg,
			}, nil
		}

		tablet := &Tablet{
			TableName:    req.TableName,
			StartRow:     "",
			EndRow:       "",
			TabletServer: assignedServer,
			Sharded:      false,
		}

		table.Tablets = append(table.Tablets, tablet)

		ms.State.TabletServers[assignedServer].TabletCount++

		// Notify the tablet server to create the table
		err = ms.notifyTabletCreateTable(assignedServer, table)
		if err != nil {
			// Rollback table creation and decrement TabletCount
			delete(ms.State.Tables, req.TableName)
			ms.State.TabletServers[assignedServer].TabletCount--
			msg := fmt.Sprintf("Failed to notify tablet server '%s' to create table '%s': %v", assignedServer, req.TableName, err)
			log.Println(msg)
			return &epb.CreateTableResponse{
				Success: false,
				Message: msg,
			}, nil
		}
	}
	log.Printf("Table '%s' created successfully.", req.TableName)
	logrus.Debugf("Create Table Master: %v", ms.State.Tables)
	return &epb.CreateTableResponse{
		Success: true,
		Message: "Table created successfully.",
	}, nil
}

func (ms *MasterServer) DeleteTable(ctx context.Context, req *epb.DeleteTableRequest) (*epb.DeleteTableResponse, error) {
	ms.State.mu.Lock()
	defer ms.State.mu.Unlock()

	table, exists := ms.State.Tables[req.TableName]
	if !exists {
		msg := fmt.Sprintf("Table '%s' does not exist.", req.TableName)
		log.Println(msg)
		return &epb.DeleteTableResponse{
			Success: false,
			Message: msg,
		}, nil
	}

	// Notify all tablet servers managing this table to delete it
	for _, tablet := range table.Tablets {
		err := ms.notifyTabletDeleteTable(tablet.TabletServer, table.Name)
		if err != nil {
			msg := fmt.Sprintf("Failed to notify tablet server '%s' to delete table '%s': %v", tablet.TabletServer, table.Name, err)
			log.Println(msg)
			return &epb.DeleteTableResponse{
				Success: false,
				Message: msg,
			}, nil
		}

		ms.State.TabletServers[tablet.TabletServer].TabletCount--
	}

	// Remove the table from master State
	delete(ms.State.Tables, req.TableName)

	log.Printf("Table '%s' deleted successfully.", req.TableName)
	return &epb.DeleteTableResponse{
		Success: true,
		Message: "Table deleted successfully.",
	}, nil
}

func (ms *MasterServer) GetTabletLocation(ctx context.Context, req *epb.GetTabletLocationRequest) (*epb.GetTabletLocationResponse, error) {
	ms.State.mu.RLock()
	defer ms.State.mu.RUnlock()

	logrus.Debugf("Get Tablet Location: %v", ms.State.Tables)

	table, exists := ms.State.Tables[req.TableName]
	if !exists {
		msg := fmt.Sprintf("Table '%s' does not exist.", req.TableName)
		log.Println(msg)
		return &epb.GetTabletLocationResponse{
			TabletServerAddress: "",
			TabletStartRow:      "",
			TabletEndRow:        "",
		}, nil
	}

	// Find the tablet responsible for the given key
	for _, tablet := range table.Tablets {
		if isKeyInRange(req.Key, tablet.StartRow, tablet.EndRow) {
			return &epb.GetTabletLocationResponse{
				TabletServerAddress: tablet.TabletServer,
				TabletStartRow:      tablet.StartRow,
				TabletEndRow:        tablet.EndRow,
			}, nil
		}
	}

	// If not found, it might indicate an inconsistency
	msg := fmt.Sprintf("No tablet found for key '%s' in table '%s'.", req.Key, req.TableName)
	log.Println(msg)
	return &epb.GetTabletLocationResponse{
		TabletServerAddress: "",
		TabletStartRow:      "",
		TabletEndRow:        "",
	}, nil
}

// InternalMasterService Implementation
func (ms *MasterServer) RegisterTablet(ctx context.Context, req *ipb.RegisterTabletRequest) (*ipb.RegisterTabletResponse, error) {
	ms.State.mu.Lock()
	defer ms.State.mu.Unlock()

	if _, exists := ms.State.TabletServers[req.TabletAddress]; exists {
		msg := fmt.Sprintf("Tablet server '%s' is already registered.", req.TabletAddress)
		log.Println(msg)
		// Even if already registered, respond with success
		return &ipb.RegisterTabletResponse{}, nil
	}

	ms.State.TabletServers[req.TabletAddress] = &TabletServerInfo{
		Address:        req.TabletAddress,
		LastHeartbeat:  time.Now(),
		RegisteredTime: time.Now(),
		TabletCount:    0, // Initialize TabletCount to zero
		MaxTableCount:  int(req.MaxTableCnt),
	}

	log.Printf("Tablet server '%s' registered successfully.", req.TabletAddress)
	return &ipb.RegisterTabletResponse{}, nil
}

func (ms *MasterServer) UnregisterTablet(ctx context.Context, req *ipb.UnregisterTabletRequest) (*ipb.UnregisterTabletResponse, error) {
	ms.State.mu.Lock()
	defer ms.State.mu.Unlock()

	_, exists := ms.State.TabletServers[req.TabletAddress]
	if !exists {
		msg := fmt.Sprintf("Tablet server '%s' goes offline.", req.TabletAddress)
		log.Println(msg)
		return &ipb.UnregisterTabletResponse{}, nil
	}

	//ms.removeTabletServer(req.TabletAddress, true)

	//Remove the tablet server from the registry
	delete(ms.State.TabletServers, req.TabletAddress)
	log.Printf("Tablet server '%s' unregistered successfully.", req.TabletAddress)

	return &ipb.UnregisterTabletResponse{}, nil
}

func (ms *MasterServer) NotifyShardRequest(ctx context.Context, req *ipb.ShardRequest) (*ipb.ShardResponse, error) {
	ms.State.mu.Lock()
	defer ms.State.mu.Unlock()

	_, exists := ms.State.Tables[req.TableName]
	if !exists {
		msg := fmt.Sprintf("ShardRequest: Table '%s' does not exist.", req.TableName)
		log.Println(msg)
		return nil, errors.New(msg)
	}

	// // Find the tablet to shard
	// var tabletToShard *Tablet
	// for _, tablet := range table.Tablets {
	// 	if tablet.TabletServer == req.TabletAddress && !tablet.Sharded {
	// 		tabletToShard = tablet
	// 		break
	// 	}
	// }

	// if tabletToShard == nil {
	// 	msg := fmt.Sprintf("ShardRequest: No suitable tablet found on server '%s' for table '%s'.", req.TabletAddress, req.TableName)
	// 	log.Println(msg)
	// 	return nil, errors.New(msg)
	// }

	// Select the least loaded tablet server to host the shard
	targetServer, err := ms.getLeastLoadedTabletServerExcluding(req.TabletAddress)
	if err != nil {
		msg := fmt.Sprintf("ShardRequest: Failed to find a target tablet server for sharding table '%s': %v", req.TableName, err)
		//logrus.Debugf(msg)
		return nil, errors.New(msg)
	}

	// Respond with the target server address
	log.Printf("ShardRequest: Assigning sharding of tablet '%s' for table '%s' to server '%s'.",
		req.TabletAddress, req.TableName, targetServer)

	return &ipb.ShardResponse{
		TargetTabletAddress: targetServer,
	}, nil
}

// NotifyShardFinish handles the completion of a shard (tablet transfer)
func (ms *MasterServer) NotifyShardFinish(ctx context.Context, req *ipb.ShardFinishNotificationRequest) (*ipb.ShardFinishNotificationResponse, error) {
	ms.State.mu.Lock()
	defer ms.State.mu.Unlock()

	table, exists := ms.State.Tables[req.TableName]
	if !exists {
		msg := fmt.Sprintf("ShardFinishNotification: Table '%s' does not exist.", req.TableName)
		log.Println(msg)
		return nil, errors.New(msg)
	}

	if len(table.Tablets) != 1 {
		msg := fmt.Sprintf("ShardFinishNotification: Table '%s' has %d tablets, expected exactly 1.", req.TableName, len(table.Tablets))
		log.Println(msg)
		return nil, errors.New(msg)
	}

	tablet := table.Tablets[0]

	if tablet.TabletServer != req.Source {
		msg := fmt.Sprintf("ShardFinishNotification: Tablet server mismatch for table '%s'. Expected source '%s', but got '%s'.",
			req.TableName, req.Source, tablet.TabletServer)
		log.Println(msg)
		return nil, errors.New(msg)
	}

	tablet.TabletServer = req.Target

	if sourceServerInfo, exists := ms.State.TabletServers[req.Source]; exists {
		if sourceServerInfo.TabletCount > 0 {
			sourceServerInfo.TabletCount--
		} else {
			log.Printf("ShardFinishNotification: Source tablet server '%s' has inconsistent TabletCount.", req.Source)
		}
	} else {
		msg := fmt.Sprintf("ShardFinishNotification: Source tablet server '%s' is not registered.", req.Source)
		log.Println(msg)
		return nil, errors.New(msg)
	}

	if targetServerInfo, exists := ms.State.TabletServers[req.Target]; exists {
		targetServerInfo.TabletCount++
	} else {
		msg := fmt.Sprintf("ShardFinishNotification: Target tablet server '%s' is not registered.", req.Target)
		log.Println(msg)
		return nil, errors.New(msg)
	}

	log.Printf("ShardFinishNotification: Successfully transferred tablet of table '%s' from server '%s' to server '%s'.",
		req.TableName, req.Source, req.Target)

	return &ipb.ShardFinishNotificationResponse{}, nil
}

// func (ms *MasterServer) NotifyShardFinish(ctx context.Context, req *ipb.ShardFinishNotificationRequest) (*ipb.ShardFinishNotificationResponse, error) {
// 	ms.State.mu.Lock()
// 	defer ms.State.mu.Unlock()

// 	table, exists := ms.State.Tables[req.TableName]
// 	if !exists {
// 		msg := fmt.Sprintf("ShardFinishNotification: Table '%s' does not exist.", req.TableName)
// 		log.Println(msg)
// 		return nil, errors.New(msg)
// 	}

// 	// Update the tablets with the new shard information
// 	originalTablet := &Tablet{
// 		StartRow:     "",
// 		EndRow:       "",
// 		TabletServer: req.Source,
// 		Sharded:      true,
// 	}

// 	newTablet := &Tablet{
// 		TableName:    req.TableName,
// 		StartRow:     "",
// 		EndRow:       "",
// 		TabletServer: req.Target,
// 		Sharded:      false,
// 	}

// 	// Verify that the original tablet exists
// 	found := false
// 	for _, tablet := range table.Tablets {
// 		if tablet.TabletServer == originalTablet.TabletServer &&
// 			tablet.StartRow == originalTablet.StartRow &&
// 			tablet.EndRow == originalTablet.EndRow {
// 			// Update the original tablet's Sharded status
// 			tablet.Sharded = originalTablet.Sharded
// 			found = true
// 			break
// 		}
// 	}

// 	if !found {
// 		msg := fmt.Sprintf("ShardFinishNotification: Original tablet not found for table '%s'.", req.TableName)
// 		log.Println(msg)
// 		return nil, errors.New(msg)
// 	}

// 	// Add the new shard tablet
// 	table.Tablets = append(table.Tablets, newTablet)

// 	// Increment TabletCount for the new shard's tablet server
// 	if serverInfo, exists := ms.State.TabletServers[newTablet.TabletServer]; exists {
// 		serverInfo.TabletCount++
// 	} else {
// 		msg := fmt.Sprintf("ShardFinishNotification: New shard tablet server '%s' is not registered.", newTablet.TabletServer)
// 		log.Println(msg)
// 		return nil, errors.New(msg)
// 	}

// 	log.Printf("ShardFinishNotification: Sharding completed for table '%s'. New shard on server '%s' with range [%s, %s).",
// 		req.TableName, newTablet.TabletServer, newTablet.StartRow, newTablet.EndRow)

// 	return &ipb.ShardFinishNotificationResponse{}, nil
// }

// MonitorHeartbeats periodically sends Heartbeat RPCs to tablet servers to check their status.
func (ms *MasterServer) MonitorHeartbeats(interval time.Duration, timeout time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		ms.State.mu.RLock()
		servers := make([]*TabletServerInfo, 0, len(ms.State.TabletServers))
		for _, server := range ms.State.TabletServers {
			servers = append(servers, server)
		}
		ms.State.mu.RUnlock()

		var wg sync.WaitGroup
		for _, server := range servers {
			wg.Add(1)
			go func(s *TabletServerInfo) {
				defer wg.Done()
				ms.sendHeartbeat(s, timeout)
			}(server)
		}
		wg.Wait()
	}
}

// sendHeartbeat sends a Heartbeat RPC to a tablet server and updates its status.
func (ms *MasterServer) sendHeartbeat(server *TabletServerInfo, timeout time.Duration) {
	conn, err := grpc.NewClient(server.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Heartbeat failed: Unable to connect to tablet server '%s': %v. Removing server.", server.Address, err)
		ms.removeTabletServer(server.Address, false)
		return
	}
	defer conn.Close()

	client := ipb.NewTabletInternalServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := client.Heartbeat(ctx, &ipb.HeartbeatRequest{})
	if err != nil || !resp.Success {
		log.Printf("Heartbeat failed: Tablet server '%s' did not respond successfully: %v. Removing server.", server.Address, err)
		ms.removeTabletServer(server.Address, false)
		return
	}

	// Update the LastHeartbeat timestamp
	ms.State.mu.Lock()
	if s, exists := ms.State.TabletServers[server.Address]; exists {
		s.LastHeartbeat = time.Now()
	}
	ms.State.mu.Unlock()

	//log.Printf("Heartbeat successful: Tablet server '%s' is online.", server.Address)
}

// removeTabletServer removes a tablet server from the registry.
// removeTabletServer removes a tablet server from the registry and handles tablet reassignment
func (ms *MasterServer) removeTabletServer(address string, crash bool) {
	ms.State.mu.Lock()
	defer ms.State.mu.Unlock()

	serverInfo, exists := ms.State.TabletServers[address]
	if !exists {
		log.Printf("ShardFinishNotification: Tablet server '%s' is not registered.", address)
		return
	}

	for _, table := range ms.State.Tables {
		for _, tablet := range table.Tablets {
			if tablet.TabletServer == address {
				// Find a new tablet server
				var newServer string
				var err error
				if crash {
					newServer, err = ms.getLeastLoadedTabletServer()
					if err != nil {
						log.Printf("Failed to find a new tablet to recover '%s': %v", table.Name, err)
						continue
					}
				} else {
					newServer, err = ms.getRecoverTabletExcluding(address)
					if err != nil {
						log.Printf("Failed to find a new tablet server for table '%s': %v", table.Name, err)
						continue
					}
				}

				// Calling the RecoverCrashedTablet method for data recovery
				err = ms.callRecoverCrashedTablet(newServer, table.Name, address)
				if err != nil {
					log.Printf("Failed to recover table '%s' on server '%s': %v", table.Name, newServer, err)
					continue
				}

				tablet.TabletServer = newServer

				// Update TabletCount
				ms.State.TabletServers[newServer].TabletCount++
				serverInfo.TabletCount--

				log.Printf("Transferred tablet of table '%s' from server '%s' to server '%s'.",
					table.Name, address, newServer)
			}
		}
	}

	delete(ms.State.TabletServers, address)
	log.Printf("Tablet server '%s' has been removed from the registry.", address)
}

// Helper Functions
// getLeastLoadedTabletServer selects the tablet server with the least number of assigned tablets.
func (ms *MasterServer) getLeastLoadedTabletServer() (string, error) {
	minCount := -1
	var selectedServer string
	for addr, serverInfo := range ms.State.TabletServers {
		if minCount == -1 || serverInfo.TabletCount < minCount {
			minCount = serverInfo.TabletCount
			selectedServer = addr
		}
	}
	if selectedServer == "" {
		return "", errors.New("no tablet servers available")
	}
	return selectedServer, nil
}

func (ms *MasterServer) getRecoverTabletExcluding(excludeAddr string) (string, error) {
	minCount := -1
	var selectedServer string
	for addr, serverInfo := range ms.State.TabletServers {
		if addr == excludeAddr {
			continue
		}

		if minCount == -1 || serverInfo.TabletCount < minCount {
			minCount = serverInfo.TabletCount
			selectedServer = addr
		}
	}
	if selectedServer == "" {
		return "", errors.New("no tablet servers available")
	}
	return selectedServer, nil
}

// getLeastLoadedTabletServerExcluding selects the least loaded tablet server excluding the specified address.
func (ms *MasterServer) getLeastLoadedTabletServerExcluding(excludeAddr string) (string, error) {
	minCount := -1
	var selectedServer string
	for addr, serverInfo := range ms.State.TabletServers {
		if addr == excludeAddr {
			continue
		}
		if minCount == -1 || serverInfo.TabletCount < minCount {
			if serverInfo.TabletCount < serverInfo.MaxTableCount {
				minCount = serverInfo.TabletCount
				selectedServer = addr
			}
		}
	}
	if selectedServer == "" {
		return "", errors.New("no tablet servers available excluding the specified server")
	}
	return selectedServer, nil
}

// isKeyInRange checks if a given key falls within the start and end row of a tablet.
// Rules:
// - If both start and end are empty, the tablet covers the entire table.
// - If start is empty, the key must be less than end.
// - If end is empty, the key must be greater than or equal to start.
// - Otherwise, the key must be in [start, end).
func isKeyInRange(key, start, end string) bool {
	if start == "" && end == "" {
		// The tablet covers the entire table
		return true
	}
	if start == "" {
		return key < end
	}
	if end == "" {
		return key >= start
	}
	return key >= start && key < end
}

// notifyTabletCreateTable sends a CreateTableInternal RPC to the specified tablet server.
// Uses grpc.Dial directly without a separate NewClient function.
func (ms *MasterServer) notifyTabletCreateTable(serverAddress string, table *Table) error {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to tablet server '%s': %v", serverAddress, err)
	}
	defer conn.Close()

	client := ipb.NewTabletInternalServiceClient(conn)

	// Prepare the request
	var cfMsgs []*ipb.CreateTableInternalRequest_ColumnFamily
	for family, columns := range table.ColumnFamilies {
		cfMsgs = append(cfMsgs, &ipb.CreateTableInternalRequest_ColumnFamily{
			FamilyName: family,
			Columns:    columns,
		})
	}

	req := &ipb.CreateTableInternalRequest{
		TableName:      table.Name,
		ColumnFamilies: cfMsgs,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.CreateTable(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateTable RPC to tablet server '%s' failed: %v", serverAddress, err)
	}

	if !resp.Success {
		return fmt.Errorf("CreateTable on tablet server '%s' failed: %s", serverAddress, resp.Message)
	}

	log.Printf("Successfully notified tablet server '%s' to create table '%s'.", serverAddress, table.Name)
	return nil
}

// notifyTabletDeleteTable sends a DeleteTableInternal RPC to the specified tablet server.
func (ms *MasterServer) notifyTabletDeleteTable(serverAddress, tableName string) error {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to tablet server '%s': %v", serverAddress, err)
	}
	defer conn.Close()

	client := ipb.NewTabletInternalServiceClient(conn)

	// Prepare the request
	req := &ipb.DeleteTableInternalRequest{
		TableName: tableName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.DeleteTable(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteTable RPC to tablet server '%s' failed: %v", serverAddress, err)
	}

	if !resp.Success {
		return fmt.Errorf("DeleteTable on tablet server '%s' failed: %s", serverAddress, resp.Message)
	}

	log.Printf("Successfully notified tablet server '%s' to delete table '%s'.", serverAddress, tableName)
	return nil
}

// callRecoverCrashedTablet connects to the target tablet server and calls RecoverCrashedTablet
func (ms *MasterServer) callRecoverCrashedTablet(targetServerAddress, tableName, crashedServerAddress string) error {
	conn, err := grpc.NewClient(targetServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to target tablet server '%s': %v", targetServerAddress, err)
	}
	defer conn.Close()

	client := ipb.NewTabletInternalServiceClient(conn)

	req := &ipb.RecoveryRequest{
		CrashedTabletAddress: crashedServerAddress,
		TableName:            tableName,
	}

	resp, err := client.RecoverCrashedTablet(context.Background(), req)
	if err != nil {
		return fmt.Errorf("RecoverCrashedTablet RPC failed on server '%s' for table '%s': %v", targetServerAddress, tableName, err)
	}

	if !resp.Success {
		return fmt.Errorf("RecoverCrashedTablet failed on server '%s' for table '%s'", targetServerAddress, tableName)
	}

	log.Printf("Successfully recovered table '%s' on server '%s' after source server '%s' crashed.",
		tableName, targetServerAddress, crashedServerAddress)

	return nil
}
