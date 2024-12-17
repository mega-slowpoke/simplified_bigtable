package main

import (
	tablet "final/bigtable/tablet"
	epb "final/proto/external-api"
	ipb "final/proto/internal-api"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// read input
	tabletAddress := flag.String("tablet_address", "", "Tablet service address")
	masterAddress := flag.String("master_address", "", "Master service address")
	maxTableCnt := flag.Int("max_table_cnt", 0, "Max table count")
	testMode := flag.Bool("test_mode", false, "Test mode")
	flag.Parse()

	// check if parameters are legal
	if *tabletAddress == "" || *masterAddress == "" || *maxTableCnt == 0 {
		log.Fatal("Please provide tablet address, master address and max shard size via command line")
	}

	setupOptions := tablet.SetupOptions{
		TabletAddress: *tabletAddress,
		MasterAddress: *masterAddress,
		MaxTableSize:  *maxTableCnt,
		TestMode:      *testMode,
	}

	tabletService, err := tablet.NewTabletService(setupOptions)
	if err != nil {
		log.Fatal(err)
	}

	host, port, err := net.SplitHostPort(tabletService.TabletAddress)
	if err != nil {
		log.Fatal("Invalid tablet address format:", err)
	}

	lis, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}
	log.Printf("TabletService is listening on %s...", tabletService.TabletAddress)

	server := grpc.NewServer()
	epb.RegisterTabletExternalServiceServer(server, tabletService)
	ipb.RegisterTabletInternalServiceServer(server, tabletService)

	// tell master a new tablet is online
	err = tabletService.RegisterMyself()
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to register %v, %v", tabletAddress, err))
	}

	if err = server.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}

	// check out of size periodically
	//go

	//
}
