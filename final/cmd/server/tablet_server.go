package main

import (
	"final/bigtable"
	proto "final/proto/external-api"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc"
	"log"
	"net"
)

//var redisAddress =

func main() {
	// TODO: setup dbPath
	db, err := OpenDB("mydb")
	if err != nil {
		log.Fatalf("Failed to open LevelDB: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// TODO: setup chubbyLock
	grpcServer := grpc.NewServer()
	chubbyLock := bigtable.NewChubbyLock()
	tabletService := bigtable.NewTabletService(db)

	proto.RegisterTabletServiceServer(grpcServer, tabletService)

	log.Println("Server is running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func OpenDB(dbPath string) (*leveldb.DB, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		log.Fatalf("Failed to open LevelDB at %s: %v", dbPath, err)
		return nil, err
	}
	return db, nil
}
