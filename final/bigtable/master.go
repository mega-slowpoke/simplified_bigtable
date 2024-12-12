package bigtable

import (
	pb "final/proto/external-api"
	"sync"
)

// TODO: Master Server
// 以下只是github里面的实现，貌似它是把master作为中介来，client不直接和tablet server交互（不确定，我没有完全读懂代码），我们可以考虑client直接和tablet server交互
// 1.服务器管理与心跳检测
// 维护 tablet_list 来记录所有可用的 tablet 服务器（通过 hostname:port 标识），通过定时（每 10 秒）向列表中的各 tablet 服务器发送心跳检测请求（通过 http://<tablet>/api/heart 接口），若检测到某个 tablet 服务器连接异常（requests.ConnectionError），则执行 recover 操作来进行故障恢复，例如将故障 tablet 服务器上的表迁移到其他可用的 tablet 服务器等。
// 2. 表相关操作处理
// 表创建：接收客户端创建表的请求，验证表名是否重复（通过检查 table_list），若不重复则为新表选择一个 tablet 服务器（借助 tablet_index 和 tablet_list），更新 table_list、tablet_dict 和 tables_info 等相关数据结构来记录新表的信息，并且向对应的 tablet 服务器发起创建表请求。
// 表查询：响应客户端获取表列表，以及获取特定表详细信息的请求，从 table_list 和 tables_info 数据结构中提取相应数据并返回给客户端。
// 表删除：处理客户端删除表的请求，先检查该表是否被客户端锁定，若未锁定则更新 table_list、tables_info 和 tablet_dict 等数据结构来移除该表相关信息，并向对应的 tablet 服务器发起删除表的请求。
// 3. 数据操作协调
// 更新行键信息：接收客户端发送的更新表行键相关信息的请求，根据请求中的数据更新 tables_info 数据结构，来反映表数据范围等信息的变化。
// 处理数据检索：对于客户端检索表中某个单元格数据的请求根据 tables_info 中记录的表数据所在 tablet 服务器以及对应的数据范围，查找对应的 tablet 服务器并向其发起获取数据的请求，将获取到的数据返回给客户端。
// 4. 并发控制
// 通过 lock_tables 数据结构来记录每个客户端（通过 client_id 标识）锁定的表信息，处理客户端锁定表（POST 请求到 /api/lock/<table_name> 接口）和释放表锁（DELETE 请求到 /api/lock/<table_name> 接口）的请求，在操作表之前会先检查锁的状态，确保同一时间只有获得锁的客户端能对相应表进行操作，保证数据操作的一致性。
// 5. 数据分片管理
// 处理分片请求：接收客户端发送的表分片请求（POST 请求到 /api/shard_request 接口），根据请求中的表名等信息将表分配到合适的 tablet 服务器上，并更新相关数据结构（如 tablet_dict 等），后续可能涉及数据迁移等操作。
// 更新分片行键信息：当分片操作完成后（接收到 POST 请求到 /api/shard_finish 接口），根据请求中的数据更新 tables_info 以及 tablet_dict 等数据结构，来准确记录分片后表数据在各 tablet 服务器上的新范围等信息。
// 6. register新的 tablet 服务器
// 接收新的 tablet 服务器注册请求RegisterTablet，将新 tablet 服务器的信息（通过 hostname:port 形式）添加到 tablet_dict 和 tablet_list 中，使其纳入 master server 的管理范围，后续可参与表的存储等相关任务。

// 核心的数据结构 （以下只是github链接上的实现，我们可以考虑其他实现方式）
// 	tablet_dict:   其结构为map，键是 tablet_address（格式类似 hostname:port），值是一个列表 [table_name]，其中 table_name 代表数据表名称。
// 		用途：用于记录每个 tablet 服务器（通过 hostname:port 标识）所关联的表（table）的名称列表。比如哪个 tablet 服务器负责管理哪些具体的表，方便后续进行相关表操作时快速查找对应的 tablet 服务器，例如在进行表的创建、数据迁移等涉及表与 tablet 关联的操作时会用到该数据结构来维护这种映射关系。
// tablet_list:（TODO:按理来说可以取消掉，直接遍历tablet_dict的键就行）
// 		是一个列表结构，元素是 tablet_address（实际格式为 hostname:port）这样的字符串。
//		用途：存储所有可用的 tablet 服务器的标识信息（以 hostname:port 形式）。它可以看作是对所有参与系统的 tablet 服务器的一个简单罗列，便于遍历所有的 tablet 服务器来执行一些全局的操作，比如心跳检测等，通过遍历这个列表能对每个 tablet 服务器进行相关的状态检查或者任务分发等操作。
// table_list: table_list = {"tables": []}，是一个字典结构，其中有一个键 tables，对应的值是一个列表，列表元素是表名称（如 CS、EE 等）。
//		用途：主要用于维护系统中所有已创建的表的名称信息，方便快速查看当前系统中存在哪些表，比如在创建新表时可以通过检查这个列表来避免表名重复，或者在删除表等操作时确认要操作的表是否存在于系统中。
// tables_info {"table1": {"name": "table1", "tablets": [{'hostname': 'localhost', 'port': '8081', 'row_from': '', 'row_to': ''}]}}，整体是一个嵌套的字典结构。外层字典的键是表名称（如 table1），对应的值又是一个字典，这个内层字典包含表的名称（和外层键值相同，有一定的冗余但方便查询使用）以及 tablets 键，其对应的值是一个列表，列表元素是包含 hostname、port、row_from、row_to 等信息的字典，用于描述表的数据在哪些 tablet 服务器上存储以及对应的数据范围。
//用途：全面记录了每个表的详细信息，包括表名、存储该表数据的 tablet 服务器相关信息（通过 hostname 和 port 确定具体服务器）以及表数据在各 tablet 服务器上所负责的数据行范围（通过 row_from 和 row_to 确定），在进行数据查询（比如根据行值找对应 tablet 服务器获取数据）、数据更新以及表的分片等操作时，都需要参考这个数据结构来确定准确的操作对象和相关范围。
//tablet_index
//定义形式：global tablet_index 且初始值 tablet_index = 0，是一个全局的整数变量。
//用途：在创建表等操作中，用于按照顺序选择一个 tablet 服务器来关联新创建的表，通常会结合 tablet_list 通过取余等操作来轮询式地为新表分配对应的 tablet 服务器，起到一个索引定位的作用，确保 tablet 服务器分配的有序性和均衡性。
//lock_tables {"client_id":[table1, table2']}，是一个字典结构，键是 client_id（代表客户端的唯一标识），值是一个列表，列表元素是表名称（如 table1、table2 等）。
//用途：用于实现表的锁定机制，记录每个客户端（通过 client_id 区分）当前锁定了哪些表，在进行涉及表的并发操作（比如多个客户端同时操作同一个表可能导致数据不一致等问题）时，通过这个数据结构来判断某个客户端是否有权限操作某个表（是否已经锁定），从而保证数据操作的一致性和正确性。

type TabletInfo struct {
	Hostname string
	Port     string
	RowFrom  string
	RowTo    string
}

type TableInfo struct {
	Name           string
	Tablets        []TabletInfo
	ColumnFamilies []string
}

type MasterServer struct {
	pb.UnimplementedMasterServiceServer

	mu          sync.RWMutex
	tabletDict  map[string][]string // key: tablet_address, value: tables hosted on that tablet
	tabletList  []string            // list of tablet servers (host:port)
	tableList   []string            // list of all tables
	tablesInfo  map[string]*TableInfo
	lockTables  map[string][]string // key: client_id, value: locked tables
	tabletIndex int                 // for round-robin assignment of tablets
}

func NewMasterServer() *MasterServer {
	return &MasterServer{
		tabletDict: make(map[string][]string),
		tablesInfo: make(map[string]*TableInfo),
		lockTables: make(map[string][]string),
	}
}

// type MasterServer struct {
// 	pb.UnimplementedMasterServiceServer
// 	TabletAssignments map[string][]*pb.TabletAssignment
// 	mu                sync.RWMutex
// }

// func NewMasterServer() *MasterServer {
// 	return &MasterServer{
// 		TabletAssignments: make(map[string][]*pb.TabletAssignment),
// 	}
// }
