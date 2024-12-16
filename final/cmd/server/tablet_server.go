package main

import (
	tablet "final/bigtable/tablet"
	epb "final/proto/external-api"
	ipb "final/proto/internal-api"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// read input
	tabletAddress := flag.String("tabletaddress", "", "Tablet service address")
	masterAddress := flag.String("masteraddress", "", "Master service address")
	maxShardSize := flag.Int("maxshardsize", 0, "Max shard size")
	flag.Parse()

	// check if parameters are legal
	if *tabletAddress == "" || *masterAddress == "" || *maxShardSize == 0 {
		log.Fatal("Please provide tablet address, master address and max shard size via command line")
	}

	// 创建TabletService实例
	setupOptions := tablet.SetupOptions{
		TabletAddress: *tabletAddress,
		MasterAddress: *masterAddress,
		MaxShardSize:  *maxShardSize,
	}
	tabletService, err := tablet.NewTabletService(setupOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 这里简单使用net/http包监听一个端口作为示例，实际中可能需要换成grpc等合适的服务框架来处理业务逻辑
	// 定义一个简单的HTTP处理函数，你可以替换为真实的服务处理逻辑
	host, port, err := net.SplitHostPort(tabletService.TabletAddress)
	if err != nil {
		log.Fatal("Invalid tablet address format:", err)
	}

	lis, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}
	log.Printf("TabletService is listening on %s...", tabletService.TabletAddress)

	// 创建gRPC服务器实例
	server := grpc.NewServer()

	// 注册相关的gRPC服务（这里需要根据实际的服务方法实现进行注册，以下只是示例占位，你要替换为真实的注册逻辑）
	epb.RegisterTabletExternalServiceServer(server, tabletService)
	ipb.RegisterTabletInternalServiceServer(server, tabletService)

	// 启动gRPC服务器，开始接收请求并处理
	if err := server.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
