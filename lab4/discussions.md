# A4
There are two TTL cleanup methods in my design (borrowed from Redis): lazy deletion and periodic active deletion.
Lazy Deletion: Get() deletes the key if it’s expired.
Active Deletion: ExpireCleanUpWatchdog() goroutine removes expired key-value pairs every second.

Overhead: O(n), where n is the number of expired key-value pairs.

In a write-intensive server, lock contention may arise, but per-shard locking in our design mitigates it.
To further improve, reducing lock granularity would allow concurrent setting and cleanup of expired keys across shards.

# B2
I'm using round robin load balancing. For this strategy, it treats all nodes as equally capable, which may not reflect the actual load distribution. If one node is under heavier load due to hosting more shards or handling more concurrent requests, round-robin could exacerbate that node's congestion by continuing to direct traffic to it. This strategy also doesn't account for differences in processing power or anything similar.

I could implement a weighted round robin load balancing to help this scenario. For this strategy, i need to assign weights to nodes based on their capacity or current load, directing more traffic to more capable nodes.

# B3
I think there are these flaws in this strategy:
1. Unnecessary Latency Increase: Each new request is made without any prior knowledge of the server's state. Therefore, every request wastes time in determining which node is operational.
2. Network Overhead: Multiple retries to different nodes generate additional network traffic, which can impact the performance of the cluster.
3. Resource Wastage: Even if a node is offline for an extended period, many requests will still attempt to send to it.

My improvement: Implement a node state tracking system and health check. Maintain a status cache of node health. If a node is known to be down or failing, skip it until its state changes. This reduces unnecessary retries and improves response times.

# B4
It will lead to inconsistency across the replicas of a shard. If some nodes successfully handle the Set call while others fail, the key-value data may be stored in an incomplete state where only a subset of nodes has the updated data.
Anomalies: Clients fetching the key through Get calls may observe different values depending on which node they connect to. Some nodes may return the new value, while others may return stale or outdated data due to the failed update.

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
Zhaoxu Wu -- Part A, Part C
Yuzhe Ruan -- Part B, Part D