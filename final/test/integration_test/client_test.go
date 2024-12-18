package integration_test

import (
	"final/bigtable"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestClientCreateTableAndDuplicateTable(t *testing.T) {
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

	log.Printf("Try to create a duplicate table : %v", tableName)
	err = client.CreateTable(tableName, columnFamily)
	if err == nil {
		log.Fatal(fmt.Printf("Created a duplicate table: %v", err))
	}
	log.Printf(err.Error())
}

func TestClientDeleteTable(t *testing.T) {
	masterAddress := "localhost:9090"

	client, err := bigtable.NewClient(masterAddress)
	if err != nil {
		log.Fatal(fmt.Printf("Failed to create client: %v", err))
	}

	tableName := "users"
	err = client.DeleteTable(tableName)
	if err != nil {
		log.Fatal(fmt.Printf("Failed to delete table: %v", err))
	}
}

func TestClientWrite(t *testing.T) {
	masterAddress := "localhost:9090"

	client, err := bigtable.NewClient(masterAddress)
	if err != nil {
		log.Fatal(fmt.Printf("Failed to create client: %v", err))
	}

	tableName := "users"
	columnFamily := "profile"
	columnName := "name"

	err = client.Write(tableName, "staff_01", columnFamily, columnName, []byte("Xiao"), time.Now().UnixNano())
	//err = client.Write(tabletAddress, "staff_02", columnFamily, columnName, []byte("Bradley_01"), time.Now().UnixNano())
	//time.Sleep(time.Microsecond)
	//err = client.Write(tabletAddress, "staff_02", columnFamily, columnName, []byte("Bradley_02"), time.Now().UnixNano())
	//time.Sleep(time.Microsecond)
	//err = client.Write(tabletAddress, "staff_02", columnFamily, columnName, []byte("Bradley_03"), time.Now().UnixNano())
	//time.Sleep(time.Microsecond)
	//err = client.Write(tabletAddress, "staff_02", columnFamily, columnName, []byte("Bradley_04"), time.Now().UnixNano())

	if err != nil {
		log.Fatal(fmt.Printf("Failed to create table: %v", err))
	}
}

func TestClientRead(t *testing.T) {

}

func TestClient(t *testing.T) {

}
