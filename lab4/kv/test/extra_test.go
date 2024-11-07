package kvtest

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"cs426.yale.edu/lab4/kv"
	"github.com/stretchr/testify/assert"
)

// Like previous labs, you must write some tests on your own.
// Add your test cases in this file and submit them as extra_test.go.
// You must add at least 5 test cases, though you can add as many as you like.
//
// You can use any type of test already used in this lab: server
// tests, client tests, or integration tests.
//
// You can also write unit tests of any utility functions you have in utils.go
//
// Tests are run from an external package, so you are testing the public API
// only. You can make methods public (e.g. utils) by making them Capitalized.

func TestShardMapUpdate(t *testing.T) {
	// Initial shard map: shard 1 assigned to nodes n1 and n2
	setup := MakeTestSetupWithoutServers(kv.ShardMapState{
		NumShards: 1,
		Nodes:     makeNodeInfos(2),
		ShardsToNodes: map[int][]string{
			1: {"n1", "n2"},
		},
	})

	setup.clientPool.OverrideSetResponse("n1")
	setup.clientPool.OverrideSetResponse("n2")
	setup.clientPool.OverrideGetResponse("n1", "value from n1", true)
	setup.clientPool.OverrideGetResponse("n2", "value from n2", true)

	// Perform Set and Get operations on key "key1"
	err := setup.Set("key1", "value1", 10*time.Second)
	assert.Nil(t, err)
	val, wasFound, err := setup.Get("key1")
	assert.Nil(t, err)
	assert.True(t, wasFound)
	assert.Contains(t, []string{"value from n1", "value from n2"}, val)

	// Update shard map: remove node n1 from shard 1
	setup.UpdateShardMapping(map[int][]string{
		1: {"n2"},
	})

	// Clear request counters
	setup.clientPool.ClearRequestsSent("n1")
	setup.clientPool.ClearRequestsSent("n2")

	// After shard map update, perform Get operations
	for i := 0; i < 10; i++ {
		val, wasFound, err = setup.Get("key1")
		assert.Nil(t, err)
		assert.True(t, wasFound)
		assert.Equal(t, "value from n2", val)
	}

	// Ensure no requests are sent to n1
	assert.Equal(t, 0, setup.clientPool.GetRequestsSent("n1"))
	// All requests should be sent to n2
	assert.Equal(t, 10, setup.clientPool.GetRequestsSent("n2"))
}

func TestServerTTLExpirationConcurrent(t *testing.T) {
	setup := MakeTestSetup(
		kv.ShardMapState{
			NumShards: 1,
			Nodes:     makeNodeInfos(1),
			ShardsToNodes: map[int][]string{
				1: {"n1"},
			},
		},
	)

	key := "ttlKey"
	value := "initialValue"

	err := setup.Set(key, value, 2*time.Second)
	assert.Nil(t, err)

	var wg sync.WaitGroup
	numGoroutines := 50

	// Start concurrent Gets before TTL expiration
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, wasFound, err := setup.Get(key)
			assert.Nil(t, err)
			assert.True(t, wasFound)
			assert.Equal(t, value, val)
		}()
	}

	wg.Wait()

	// Sleep to allow TTL to expire
	time.Sleep(3 * time.Second)

	// Verify the key has expired
	val, wasFound, err := setup.Get(key)
	assert.Nil(t, err)
	assert.False(t, wasFound)
	assert.Equal(t, "", val)

	setup.Shutdown()
}

func TestServerLargeNumberOfKeysWithVaryingTTLs(t *testing.T) {
	// Setup a server with a single node hosting multiple shards
	shardMapState := kv.ShardMapState{
		NumShards: 5,
		Nodes:     makeNodeInfos(1), // Single node "n1"
		ShardsToNodes: map[int][]string{
			1: {"n1"},
			2: {"n1"},
			3: {"n1"},
			4: {"n1"},
			5: {"n1"},
		},
	}

	setup := MakeTestSetup(shardMapState)

	// Insert a large number of keys with TTLs varying from 1 to 10 seconds
	numKeys := 1000
	keys := RandomKeys(numKeys, 20)
	for i, key := range keys {
		ttl := time.Duration((i%10)+1) * time.Second // TTLs from 1s to 10s
		err := setup.NodeSet("n1", key, fmt.Sprintf("value-%d", i), ttl)
		assert.Nil(t, err)
	}

	// Verify that keys are accessible before expiration
	for i, key := range keys {
		val, wasFound, err := setup.NodeGet("n1", key)
		assert.Nil(t, err)
		assert.True(t, wasFound)
		assert.Equal(t, fmt.Sprintf("value-%d", i), val)
	}

	// Wait for the longest TTL to expire plus buffer time
	time.Sleep(11 * time.Second)

	// Verify that all keys have expired
	for _, key := range keys {
		_, wasFound, err := setup.NodeGet("n1", key)
		assert.Nil(t, err)
		assert.False(t, wasFound)
	}

	setup.Shutdown()
}

func TestClientLoadBalancing(t *testing.T) {
	// Ensure that the test setup is initialized properly and that it includes all nodes
	setup := MakeTestSetupWithoutServers(MakeTwoNodeBothAssignedSingleShard())

	// Override client pool responses to return mock data
	setup.clientPool.OverrideGetResponse("n1", "value-from-n1", true)
	setup.clientPool.OverrideGetResponse("n2", "value-from-n2", true)

	numRequests := 100
	for i := 0; i < numRequests; i++ {
		val, wasFound, err := setup.Get("key-for-load-balancing")
		assert.Nil(t, err)
		assert.True(t, wasFound)
		assert.Contains(t, []string{"value-from-n1", "value-from-n2"}, val)
	}

	// Check distribution of requests
	n1Requests := setup.clientPool.GetRequestsSent("n1")
	n2Requests := setup.clientPool.GetRequestsSent("n2")
	assert.Greater(t, n1Requests, numRequests/4, "Expected n1 to handle part of the load")
	assert.Greater(t, n2Requests, numRequests/4, "Expected n2 to handle part of the load")
	assert.LessOrEqual(t, n1Requests, 3*numRequests/4, "Load should be balanced")
	assert.LessOrEqual(t, n2Requests, 3*numRequests/4, "Load should be balanced")
}

func TestIntegrationConcurrentGetSet(t *testing.T) {
	setup := MakeTestSetup(MakeFourNodesWithFiveShards())

	numKeys := 50
	numGoroutines := 10
	numIterations := 20

	keys := RandomKeys(numKeys, 15)
	vals := RandomKeys(numKeys, 30)

	for i := 0; i < numKeys; i++ {
		err := setup.Set(keys[i], vals[i], 100*time.Second)
		assert.Nil(t, err)
	}

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < numIterations; j++ {
				for k := 0; k < numKeys; k++ {
					_, _, err := setup.Get(keys[k])
					assert.Nil(t, err)
				}
			}
			wg.Done()
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < numIterations; j++ {
				for k := 0; k < numKeys; k++ {
					newVal := fmt.Sprintf("%s-updated-%d", vals[k], j)
					err := setup.Set(keys[k], newVal, 100*time.Second)
					assert.Nil(t, err)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	setup.Shutdown()
}
