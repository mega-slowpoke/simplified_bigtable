package kv

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cs426.yale.edu/lab4/kv/proto"
	"github.com/sirupsen/logrus"
)

type ShardState struct {
	mu        sync.Mutex
	nextIndex int
}

type Kv struct {
	shardMap   *ShardMap
	clientPool ClientPool

	// Add any client-side state you want here
	shardStates map[int]*ShardState // Map from shardId to ShardState
}

func MakeKv(shardMap *ShardMap, clientPool ClientPool) *Kv {
	kv := &Kv{
		shardMap:    shardMap,
		clientPool:  clientPool,
		shardStates: make(map[int]*ShardState),
	}

	// Add any initialization logic
	// Initialize ShardState for each shard
	shardNum := kv.shardMap.NumShards()
	for shardId := 0; shardId < shardNum; shardId++ {
		kv.shardStates[shardId] = &ShardState{}
	}

	return kv
}

func (kv *Kv) Get(ctx context.Context, key string) (string, bool, error) {
	// Trace-level logging -- you can remove or use to help debug in your tests
	// with `-log-level=trace`. See the logging section of the spec.
	logrus.WithFields(
		logrus.Fields{"key": key},
	).Trace("client sending Get() request")

	// Call tryAllNodes() to try all nodes in the shard
	var value string
	var wasFound bool
	err := kv.tryAllNodes(key, func(client proto.KvClient) error {
		// Call Get() on the client
		resq, err := client.Get(ctx, &proto.GetRequest{Key: key})
		if err != nil {
			return err
		}
		value = resq.Value
		wasFound = resq.WasFound
		return nil
	})

	if err != nil {
		return "", false, err
	}

	return value, wasFound, nil
}

func (kv *Kv) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	logrus.WithFields(
		logrus.Fields{"key": key},
	).Trace("client sending Set() request")

	// Call operateOnAllNodes() to send Set() requests to all nodes in the shard
	err := kv.operateOnAllNodes(key, func(client proto.KvClient) error {
		// Call Set() on the client
		_, err := client.Set(ctx, &proto.SetRequest{Key: key, Value: value, TtlMs: int64(ttl / time.Millisecond)})
		return err
	})

	return err
}

func (kv *Kv) Delete(ctx context.Context, key string) error {
	logrus.WithFields(
		logrus.Fields{"key": key},
	).Trace("client sending Delete() request")

	// Call operateOnAllNodes() to send Delete() requests to all nodes in the shard
	err := kv.operateOnAllNodes(key, func(client proto.KvClient) error {
		// Call Delete() on the client
		_, err := client.Delete(ctx, &proto.DeleteRequest{Key: key})
		return err
	})

	return err
}

// If one node is down, the client should try another node in the shard until it finds a node that is up or until it has tried all nodes in the shard.
func (kv *Kv) tryAllNodes(key string, operation func(client proto.KvClient) error) error {
	// Get the corresponding shard for the key
	shardNum := kv.shardMap.NumShards()
	shardId := GetShardForKey(key, shardNum)
	shardIndex := shardId - 1

	// Get the nodes that host the shard
	nodes := kv.shardMap.NodesForShard(shardId)

	if len(nodes) == 0 {
		return fmt.Errorf("no nodes assigned to shard %d", shardId)
	}

	// Use round-robin to select a node
	shardState := kv.shardStates[shardIndex]
	shardState.mu.Lock()
	defer shardState.mu.Unlock()

	// First, try the one selected by round-robin
	// If it fails, try all other nodes
	var err error
	for i := 0; i < len(nodes); i++ {
		node := nodes[shardState.nextIndex%len(nodes)]
		shardState.nextIndex = (shardState.nextIndex + 1) % len(nodes)

		// Get the client for the node
		var client proto.KvClient
		client, err = kv.clientPool.GetClient(node)
		if err != nil {
			continue
		}

		// Call the operation
		err = operation(client)
		if err == nil {
			return nil
		}
	}

	return err
}

func (kv *Kv) operateOnAllNodes(key string, operation func(client proto.KvClient) error) error {
	// Get the corresponding shard for the key
	shardNum := kv.shardMap.NumShards()
	shardId := GetShardForKey(key, shardNum)

	// Get the nodes that host the shard
	nodes := kv.shardMap.NodesForShard(shardId)

	if len(nodes) == 0 {
		return fmt.Errorf("no nodes assigned to shard %d", shardId)
	}

	// Perform the operation on all nodes
	// Calls must be sent concurrently
	// If any node fails (in GetClient() or in the Set() RPC), return the error eventually, but still send the requests to all nodes
	// If multiple nodes fail, you may return any error heard (or an error which wraps them)
	var wg sync.WaitGroup
	var errs []error = make([]error, len(nodes))
	for i, node := range nodes {
		wg.Add(1)
		go func(idx int, node string) {
			defer wg.Done()

			// Get the client for the node
			var client proto.KvClient
			client, errs[i] = kv.clientPool.GetClient(node)
			if errs[i] != nil {
				return
			}

			// Call the operation
			errs[i] = operation(client)
		}(i, node)
	}

	wg.Wait()
	// Return the first error heard
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
