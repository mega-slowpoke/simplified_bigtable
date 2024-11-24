package bigtable

// TODO: Master Server
// 1. Tablet 生命周期管理及负载均衡：Master Server 负责决定哪些 tablet 分配给哪些服务器（即 Tablet Servers）。
//	  - 负责管理表的创建、删除和更新操作 （通过调用Metadata API ??), 当新表创建时，Master 负责决定如何分配 tablets
//    - 负载均衡：当某个服务器负载过高或某些 tablets 存在性能瓶颈时，Master 会决定将某些 tablets 从一个服务器迁移到另一个服务器，从而保证系统负载均衡
// 2. 故障检测： (paper 5.2) Master 定期检查每个 Tablet Server 的健康状况。当发现某个服务器宕机或出现故障时，Master 会执行必要的操作，如将该服务器上的 tablets 转移到其他healthy的服务器上，
//		  		如果一个tablets由多个servers负责，也可以把所有的请求route到healthy节点而暂停负载均衡
//				实现细节：Chubby的替代方案(Redis分布式锁? 多个全局锁? single node etcd?)， 各自有trade-off
// 3. Fault Tolerance:
//   		- 当 Tablet Server 出现故障时，Master 需要能够快速识别并确保恢复后能恢复数据(?)
// 	  		- 当 Metadata 出现 crash, 需要能够恢复metadata的内容，通过API
