package main

import (
	"context"
	tablet "final/bigtable/tablet"
	epb "final/proto/external-api"
	ipb "final/proto/internal-api"
	"flag"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// read input
	tabletAddress := flag.String("tablet_address", "", "Tablet service address")
	masterAddress := flag.String("master_address", "", "Master service address")
	maxTableCnt := flag.Int("max_table_cnt", 0, "Max table count")
	checkMaxCntPeriod := flag.Int("check_max_period", 100, "Microsecond, Modify this if you want the tablet to check if tables it takes care of are greater than max cnt more frequently")
	flag.Parse()

	// check if parameters are legal
	if *tabletAddress == "" || *masterAddress == "" || *maxTableCnt == 0 {
		log.Fatal("Please provide tablet address, master address and max shard size via command line")
	}

	setupOptions := tablet.SetupOptions{
		TabletAddress: *tabletAddress,
		MasterAddress: *masterAddress,
		MaxTableSize:  *maxTableCnt,
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

	// Graceful shutdown logic
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run the server in a goroutine
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	go tabletService.PeriodicallyCheckMaxSize(ctx, *checkMaxCntPeriod)

	// Wait for termination signal
	<-ctx.Done()
	log.Println("Shutting down server gracefully...")

	//  clean up
	server.GracefulStop()
	tabletService.UnRegisterMyself()
	tabletService.CloseMasterConnection()
	tabletService.CloseAllTables()

	logrus.Println("Server exited successfully.")
}
