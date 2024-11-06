package kvtest

import (
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
