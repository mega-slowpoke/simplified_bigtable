package integration_test

import (
	"final/bigtable"
	"fmt"
	"log"
	"testing"
)

func TestClientCreateTable(t *testing.T) {
	masterAddress := "localhost:9090"

	client, err := bigtable.NewClient(masterAddress)
	if err != nil {
		log.Fatal(fmt.Printf("Failed to create client: %v", err))
	}

	tableName := "users"
	columnFamily := map[string][]string{
		"profile":  {"name", "email", "phone", "age"},
		"activity": {"last_login", "last_post", "last_comment"},
	}

	err = client.CreateTable(tableName, columnFamily)
	if err != nil {
		log.Fatal(fmt.Printf("Failed to create table: %v", err))
	}
}
