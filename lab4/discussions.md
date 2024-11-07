# A4
There are two TTL cleanup methods in my design (borrowed from Redis): lazy deletion and periodic active deletion.
Lazy Deletion: Get() deletes the key if it’s expired.
Active Deletion: ExpireCleanUpWatchdog() goroutine removes expired key-value pairs every second.

Overhead: O(n), where n is the number of expired key-value pairs.

In a write-intensive server, lock contention may arise, but per-shard locking in our design mitigates it.
To further improve, reducing lock granularity would allow concurrent setting and cleanup of expired keys across shards.

# D2
## Experiment 1
- Command: go run cmd/stress/tester.go --shardmap shardmaps/test-3-node-100-shard.json --get-qps 200 --set-qps 50 --num-keys 100
- What I'm testing: I am testing whether my program can handle regular stress tests without errors, while also modifying the shard map to test update behavior.
- Result: 
```
Stress test completed!
Get requests: 12020/12020 succeeded = 100.000000% success rate
Set requests: 3020/3020 succeeded = 100.000000% success rate
Correct responses: 11970/11970 = 100.000000%
Total requests: 15040 = 250.657924 QPS
```
## Experiment 2
- Command: go run cmd/stress/tester.go --shardmap shardmaps/test-5-node.json --get-qps 200 --set-qps 100 --num-keys 1
- What I'm testing: I'm testing a hot-key workload
- Result:
```
Stress test completed!
Get requests: 12017/12017 succeeded = 100.000000% success rate
Set requests: 6020/6020 succeeded = 100.000000% success rate
Correct responses: 8661/8661 = 100.000000%
Total requests: 18037 = 300.592655 QPS
```

# Group work