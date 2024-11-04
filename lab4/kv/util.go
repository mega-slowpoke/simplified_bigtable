package kv

import "hash/fnv"

/// This file can be used for any common code you want to define and separate
/// out from server.go or client.go

func GetShardForKey(key string, numShards int) int {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	return int(hasher.Sum32())%numShards + 1
}

func CheckKeyIsNull(key string) bool {
	return len(key) == 0
}

func isHosted(hostShardIds map[int]struct{}, targetShardId int) bool {
	_, ok := hostShardIds[targetShardId]
	return ok
}

func GetShardIdx(shardName int) int {
	return shardName - 1
}
