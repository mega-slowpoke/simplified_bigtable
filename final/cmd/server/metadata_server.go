package main

import (
	"final/bigtable"
	proto "final/proto/external-api"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 8080, "Port to listen on")
)

func main() {

	// 解析命令行参数
	flag.Parse()

	// 创建 gRPC 服务实例
	server := grpc.NewServer()

	// 初始化 MetadataService 实例
	metadataService := bigtable.NewMetadataService()

	// 注册 MetadataService 到 gRPC 服务
	proto.RegisterMetadataServiceServer(server, metadataService)

	// 启动监听
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Metadata server listening on %v", lis.Addr())

	// 启动 gRPC 服务
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
