package kv

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	shardsLock  sync.RWMutex
	cleanupChan chan struct{} // close the watch dog
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

	newHostShards := server.shardMap.ShardsForNode(server.nodeName)

	toRemove := make(map[int]struct{})
	toAdd := make(map[int]struct{})

	// Lock for reading `server.shards`
	server.shardsLock.RLock()
	for shardIdx := range server.shards {
		if !ContainShard(newHostShards, shardIdx) {
			toRemove[shardIdx] = struct{}{}
		}
	}
	server.shardsLock.RUnlock()

	for _, shardIdx := range newHostShards {
		server.shardsLock.RLock()
		_, exists := server.shards[shardIdx]
		server.shardsLock.RUnlock()
		if !exists {
			toAdd[shardIdx] = struct{}{}
		}
	}

	// Process shard additions first
	for shardIdx := range toAdd {
		server.CopyShardData(shardIdx)
	}

	// Process shard removals
	for shardIdx := range toRemove {
		server.removeShard(shardIdx)
	}
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
		expiryTime: time.Now().Add(time.Duration(TtlMs) * time.Millisecond),
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

	// Lock for reading `server.shards`
	server.shardsLock.RLock()
	shardPtr, exists := server.shards[shardIdx]
	server.shardsLock.RUnlock()

	if !exists {
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

	server.shardsLock.RLock()
	shard, ok := server.shards[expectedShardId]
	server.shardsLock.RUnlock()

	if !ok {
		return nil
	}
	return shard
}

func (server *KvServerImpl) ExpireCleanUpWatchdog() {
	for {
		select {
		case <-server.cleanupChan:
			return
		default:
			time.Sleep(time.Microsecond * 10)

			server.shardsLock.RLock()
			shardsCopy := make(map[int]*Shard)
			for shardIdx, shard := range server.shards {
				shardsCopy[shardIdx] = shard
			}
			server.shardsLock.RUnlock()

			// Iterate over a copy to avoid holding the lock for long
			for _, shard := range shardsCopy {
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

		// FIX1: time.Duration 默认是 nanosecond，所以要乘以 time.Millisecond
		for _, pairs := range response.Values {
			k := pairs.Key
			v := pairs.Value
			ttlRemainingMs := pairs.TtlMsRemaining
			if ttlRemainingMs <= 0 {
				continue
			}
			newShardData[k] = Value{
				content:    v,
				expiryTime: time.Now().Add(time.Duration(ttlRemainingMs) * time.Millisecond),
			}
		}

		server.shardsLock.Lock()
		if _, exists := server.shards[shardIdx]; !exists {
			server.shards[shardIdx] = &Shard{
				dataMap: newShardData,
				mu:      sync.RWMutex{},
			}
		} else {
			shardPtr := server.shards[shardIdx]
			shardPtr.mu.Lock()
			shardPtr.dataMap = newShardData
			shardPtr.mu.Unlock()
		}
		server.shardsLock.Unlock()

		logrus.WithFields(
			logrus.Fields{"node": server.nodeName},
		).Debugf("Successfully copied shard %d data from peer", shardIdx)
		return
	}

	// If all attempts to copy data failed, initialize as an empty shard
	server.shardsLock.Lock()
	if _, exists := server.shards[shardIdx]; !exists {
		server.shards[shardIdx] = &Shard{
			dataMap: newShardData,
			mu:      sync.RWMutex{},
		}
	}
	server.shardsLock.Unlock()

	logrus.WithFields(
		logrus.Fields{"node": server.nodeName},
	).Debugf("Failed to copy shard %d data from all peers, initializing as empty", shardIdx)
}

func (server *KvServerImpl) removeShard(shardIdx int) {
	// Lock for writing to `server.shards`
	server.shardsLock.Lock()
	defer server.shardsLock.Unlock()

	shardPtr, exists := server.shards[shardIdx]
	if !exists {
		return
	}
	shardPtr.mu.Lock()
	delete(server.shards, shardIdx)
	shardPtr.mu.Unlock()
}
