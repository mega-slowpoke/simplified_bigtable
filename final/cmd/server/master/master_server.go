package main

import (
	"final/bigtable"
	epb "final/proto/external-api"
	ipb "final/proto/internal-api"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	masterAddress := flag.String("master_address", "", "Master service address")
	flag.Parse()

	if *masterAddress == "" {
		log.Fatal("Please provide master address to start master server")
	}

	masterServer := bigtable.NewMasterServer()

	lis, err := net.Listen("tcp", *masterAddress)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to listen on %s: %v", *masterAddress, err))
	}

	grpcServer := grpc.NewServer()
	epb.RegisterMasterExternalServiceServer(grpcServer, masterServer)
	ipb.RegisterMasterInternalServiceServer(grpcServer, masterServer)

	heartbeatInterval := 10 * time.Second
	heartbeatTimeout := 5 * time.Second
	go masterServer.MonitorHeartbeats(heartbeatInterval, heartbeatTimeout)

	log.Printf("Master server is running on %s...", *masterAddress)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
