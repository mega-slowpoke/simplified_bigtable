package kv

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"

	"cs426.yale.edu/lab4/kv/proto"
	"github.com/sirupsen/logrus"
)

type KvServerImpl struct {
	proto.UnimplementedKvServer
	nodeName string

	shardMap   *ShardMap // keep track of sharding state
	listener   *ShardMapListener
	clientPool ClientPool
	shutdown   chan struct{}

	shards      map[int]*Shard // keep track of actually sharding data
	cleanupChan chan struct{}  // close the watch dog
	updateLock  sync.Mutex
}

type Shard struct {
	dataMap map[string]Value
	mu      sync.RWMutex
}

type Value struct {
	content    string
	expiryTime time.Time
}

func (server *KvServerImpl) handleShardMapUpdate() {
	server.updateLock.Lock()
	defer server.updateLock.Unlock()

	// move to a map so the lookup can be O(1)
	newHostShards := server.shardMap.ShardsForNode(server.nodeName)

	toRemove := make(map[int]struct{})
	toAdd := make(map[int]struct{})

	// in cur shards but not in newHostShards
	for shardIdx := range server.shards {
		if !ContainShard(newHostShards, shardIdx) {
			toRemove[shardIdx] = struct{}{}
		}
	}

	// in newHostShards, not in cur shards
	for _, shardIdx := range newHostShards {
		if _, exists := server.shards[shardIdx]; !exists {
			toAdd[shardIdx] = struct{}{}
		}
	}

	// remove toRemove
	var wgRemove sync.WaitGroup
	for shardIdx := range toRemove {
		wgRemove.Add(1)
		go func(idx int) {
			lock := &server.shards[idx].mu
			lock.Lock()
			delete(server.shards, idx)
			lock.Unlock()
		}(shardIdx)
	}

	// add toAdd
	var wgAdd sync.WaitGroup
	for shardIdx := range toAdd {
		wgAdd.Add(1)
		go func(idx int) {
			defer wgAdd.Done()
			server.CopyShardData(idx)
		}(shardIdx)
	}

	wgAdd.Wait()
	wgRemove.Wait()
}

func (server *KvServerImpl) shardMapListenLoop() {
	listener := server.listener.UpdateChannel()
	for {
		select {
		case <-server.shutdown:
			return
		case <-listener:
			server.handleShardMapUpdate()
		}
	}
}

func MakeKvServer(nodeName string, shardMap *ShardMap, clientPool ClientPool) *KvServerImpl {
	listener := shardMap.MakeListener()
	server := KvServerImpl{
		nodeName:    nodeName,
		shardMap:    shardMap,
		listener:    &listener,
		clientPool:  clientPool,
		shutdown:    make(chan struct{}),
		cleanupChan: make(chan struct{}),
		shards:      make(map[int]*Shard),
	}

	logrus.WithFields(
		logrus.Fields{"node": server.nodeName},
	).Debugf("Create Server")

	initialShards := shardMap.ShardsForNode(nodeName)
	for _, shardIdx := range initialShards {
		// copy shard data when initialize
		server.CopyShardData(shardIdx)
	}

	// clean expiry
	go server.ExpireCleanUpWatchdog()

	go server.shardMapListenLoop()
	server.handleShardMapUpdate()
	return &server
}

func (server *KvServerImpl) Shutdown() {
	server.shutdown <- struct{}{}
	server.cleanupChan <- struct{}{}
	close(server.cleanupChan)
	server.listener.Close()
}

func (server *KvServerImpl) Get(
	ctx context.Context,
	request *proto.GetRequest,
) (*proto.GetResponse, error) {
	// Trace-level logging for node receiving this request (enable by running with -log-level=trace),
	// feel free to use Trace() or Debug() logging in your code to help debug tests later without
	// cluttering logs by default. See the logging section of the spec.
	logrus.WithFields(
		logrus.Fields{"node": server.nodeName, "key": request.Key},
	).Trace("node received Get() request")

	// check key is not null
	key := request.Key
	if CheckKeyIsNull(key) {
		return &proto.GetResponse{}, status.Error(codes.InvalidArgument, "Empty keys are invalid")
	}

	//  get corresponding shard
	shardPtr := server.GetCorrespondingShard(key)
	if shardPtr == nil {
		return &proto.GetResponse{}, status.Error(codes.NotFound, "this request is not supposed to be directed to this shard")
	}

	shardPtr.mu.Lock()
	defer shardPtr.mu.Unlock()
	val, ok := shardPtr.dataMap[key]

	// 2. check exist
	if !ok {
		response := &proto.GetResponse{
			Value:    "",
			WasFound: false,
		}
		return response, nil
	}

	// 3. check expire
	if val.expiryTime.Before(time.Now()) {
		// delete key and return none
		delete(shardPtr.dataMap, key)
		response := &proto.GetResponse{
			Value:    "",
			WasFound: false,
		}
		return response, nil
	}

	// 4. not expire
	response := &proto.GetResponse{
		Value:    val.content,
		WasFound: true,
	}
	return response, nil
}

func (server *KvServerImpl) Set(
	ctx context.Context,
	request *proto.SetRequest,
) (*proto.SetResponse, error) {
	logrus.WithFields(
		logrus.Fields{"node": server.nodeName, "key": request.Key},
	).Trace("node received Set() request")

	// check key is not empty
	key := request.Key
	if CheckKeyIsNull(key) {
		return &proto.SetResponse{}, status.Error(codes.InvalidArgument, "Empty keys are invalid")
	}

	//  get corresponding shard
	shardPtr := server.GetCorrespondingShard(key)
	if shardPtr == nil {
		return &proto.SetResponse{}, status.Error(codes.NotFound, "this request is not supposed to be directed to this shard")
	}

	//  reset value
	valueContent := request.Value
	TtlMs := request.TtlMs
	shardPtr.mu.Lock()
	defer shardPtr.mu.Unlock()
	shardPtr.dataMap[key] = Value{
		content:    valueContent,
		expiryTime: time.Now().Add(time.Duration(TtlMs)),
	}

	return &proto.SetResponse{}, nil
}

func (server *KvServerImpl) Delete(
	ctx context.Context,
	request *proto.DeleteRequest,
) (*proto.DeleteResponse, error) {
	logrus.WithFields(
		logrus.Fields{"node": server.nodeName, "key": request.Key},
	).Trace("node received Delete() request")

	// check key is not empty
	key := request.Key
	if CheckKeyIsNull(key) {
		return &proto.DeleteResponse{}, status.Error(codes.InvalidArgument, "Empty keys are invalid")
	}

	//  get corresponding shard
	shardPtr := server.GetCorrespondingShard(key)
	if shardPtr == nil {
		return &proto.DeleteResponse{}, status.Error(codes.NotFound, "this request is not supposed to be directed to this shard")
	}

	shardPtr.mu.Lock()
	defer shardPtr.mu.Unlock()
	delete(shardPtr.dataMap, key)
	return &proto.DeleteResponse{}, nil
}

func (server *KvServerImpl) GetShardContents(
	ctx context.Context,
	request *proto.GetShardContentsRequest,
) (*proto.GetShardContentsResponse, error) {
	shardIdx := int(request.Shard)
	shardPtr := server.shards[shardIdx]
	if shardPtr == nil {
		logrus.WithFields(
			logrus.Fields{"node": server.nodeName},
		).Debugf("get shard content call fails")
		return nil, status.Error(codes.NotFound, "this request is not supposed to be directed to this shard")
	}

	shardPtr.mu.Lock() // Lock for reading to avoid write conflicts
	defer shardPtr.mu.Unlock()

	dataList := make([]*proto.GetShardValue, 0)
	for k, v := range shardPtr.dataMap {
		if v.expiryTime.After(time.Now()) {
			shardValue := &proto.GetShardValue{
				Key:            k,
				Value:          v.content,
				TtlMsRemaining: v.expiryTime.UnixMilli() - time.Now().UnixMilli(),
			}
			dataList = append(dataList, shardValue)
		}
	}

	response := &proto.GetShardContentsResponse{
		Values: dataList,
	}

	return response, nil
}

func (server *KvServerImpl) GetCorrespondingShard(key string) *Shard {
	numShards := server.shardMap.NumShards()
	expectedShardId := GetShardForKey(key, numShards)
	if !isHosted(server.shards, expectedShardId) {
		return nil
	}
	return server.shards[expectedShardId]
}

func (server *KvServerImpl) ExpireCleanUpWatchdog() {
	for {
		select {
		case <-server.cleanupChan:
			return
		default:
			time.Sleep(time.Microsecond * 10)
			for _, shard := range server.shards {
				shard.mu.Lock()
				for key, val := range shard.dataMap {
					if val.expiryTime.Before(time.Now()) {
						delete(shard.dataMap, key)
					}
				}
				shard.mu.Unlock()
			}
		}

	}
}

func (server *KvServerImpl) CopyShardData(shardIdx int) {
	logrus.WithFields(
		logrus.Fields{"node": server.nodeName},
	).Debugf("start to copy shardData %d", shardIdx)

	peers := server.shardMap.NodesForShard(shardIdx)
	newShardData := make(map[string]Value)

	for _, peer := range peers {
		// don't call the current node
		if peer == server.nodeName {
			continue
		}

		client, err := server.clientPool.GetClient(peer)
		if err != nil {
			logrus.WithFields(
				logrus.Fields{"node": server.nodeName},
			).Debugf("Failed to get client for peer %s: %v", peer, err)
			continue
		}

		response, err := client.GetShardContents(context.Background(), &proto.GetShardContentsRequest{
			Shard: int32(shardIdx),
		})
		if err != nil {
			logrus.WithFields(
				logrus.Fields{"node": server.nodeName},
			).Debugf("fetch shard %d contents from peer %s: %v", shardIdx, peer, err)
			continue
		}

		for _, pairs := range response.Values {
			k := pairs.Key
			v := pairs.Value
			ttlRemaining := pairs.TtlMsRemaining
			newShardData[k] = Value{
				content:    v,
				expiryTime: time.Now().Add(time.Duration(ttlRemaining)),
			}
		}
		logrus.WithFields(
			logrus.Fields{"node": server.nodeName},
		).Debugf("get all shard %v contents %v", shardIdx, newShardData)

		shardPtr, ok := server.shards[shardIdx]
		if !ok {
			// if such shard not exist
			server.shards[shardIdx] = &Shard{
				dataMap: newShardData,
				mu:      sync.RWMutex{},
			}
		} else {
			// if shardIdx exists, lock first to avoid race condition
			shardPtr.mu.Lock()
			server.shards[shardIdx].dataMap = newShardData
			shardPtr.mu.Unlock()
		}
		server.shards[shardIdx].mu.Lock()
		logrus.WithFields(
			logrus.Fields{"node": server.nodeName},
		).Debugf("Successfully copied shard %d data from peer content: %v", shardIdx, server.shards[shardIdx].dataMap)
		server.shards[shardIdx].mu.Unlock()
		return
	}

	// fail to get shard data from peers
	shardPtr, ok := server.shards[shardIdx]
	if !ok {
		server.shards[shardIdx] = &Shard{
			dataMap: newShardData,
			mu:      sync.RWMutex{},
		}
	} else {
		shardPtr.mu.Lock()
		server.shards[shardIdx].dataMap = newShardData
		shardPtr.mu.Unlock()
	}

	server.shards[shardIdx].mu.Lock()
	logrus.WithFields(
		logrus.Fields{"node": server.nodeName},
	).Debugf("Failed to copy shard %d data from all peers, initializing as empty", shardIdx)
	server.shards[shardIdx].mu.Unlock()
}
