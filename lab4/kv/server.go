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

	shards map[int]*Shard // keep track of actually sharding data

	cleanupChan chan struct{} // close the watch dog
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
	// TODO: Part C
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
	}

	logrus.WithFields(
		logrus.Fields{"node": server.nodeName},
	).Debugf("Create Server")
	numShards := shardMap.NumShards()
	server.shards = make(map[int]*Shard, numShards)
	for i := 0; i < numShards; i++ {
		server.shards[i] = &Shard{
			dataMap: make(map[string]Value),
			mu:      sync.RWMutex{},
		}
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
	panic("TODO: Part C")
}

func (server *KvServerImpl) GetCorrespondingShard(key string) *Shard {
	numShards := server.shardMap.NumShards()
	expectedShardId := GetShardForKey(key, numShards)
	hostShards := server.shardMap.ShardsForNode(server.nodeName)
	if !isHosted(hostShards, expectedShardId) {
		return nil
	}
	return server.shards[expectedShardId-1]
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
