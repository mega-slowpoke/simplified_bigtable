package bigtable

// TODO: Master Server
// 1. Tablet 生命周期管理及负载均衡：Master Server 负责决定哪些 tablet 分配给哪些服务器（即 Tablet Servers）。
//	  - 负责管理表的创建、删除和更新操作 （通过调用Metadata API ??), 当新表创建时，Master 负责决定如何分配 tablets
//    - 负载均衡：当某个服务器负载过高或某些 tablets 存在性能瓶颈时，Master 会决定将某些 tablets 从一个服务器迁移到另一个服务器，从而保证系统负载均衡
// 2. Tablet 的拆分与合并：
//    - Tablet 拆分：随着表的大小不断增大，Master 负责将大的 tablet 拆分成多个更小的 tablet，以便分布到不同的 Tablet Servers 上。这是为了避免单个 tablet 过大，导致性能瓶颈。
//	  - Tablet 合并：如果某些 tablets 过小或者服务器的负载过轻，Master 可以决定将多个 tablets 合并成一个，优化系统的性能和存储。
// 3. Fault Tolerance:
//	  - 故障检测：Master 定期检查每个 Tablet Server 的健康状况。当发现某个服务器宕机或出现故障时，Master 会执行必要的操作，如将该服务器上的 tablets 转移到其他healthy的服务器上，
//	  		如果一个tablets由多个servers负责，也可以把所有的请求route到healthy节点而暂停负载均衡
//    - 恢复机制：
//   		- 当 Tablet Server 出现故障时，Master 需要能够快速识别并恢复数据，保证系统的可靠性和可用性。
// 	  		- 当 Metadata 出现 crash, 需要能够恢复metadata的内容，通过API（暂定，也可以假设metadata server永远不crash）

// 实现细节？Metadata server 怎么知道tablet的大小？
