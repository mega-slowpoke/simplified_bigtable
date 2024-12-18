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
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	tableName := "users"
	columnFamilies := map[string][]string{
		"profile":  {"name", "email", "phone", "age"},
		"activity": {"last_login", "last_post", "last_comment"},
	}

	err = client.CreateTable(tableName, columnFamilies)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	t.Logf("Try to create a duplicate table: %v", tableName)
	err = client.CreateTable(tableName, columnFamilies)
	if err == nil {
		t.Fatalf("Expected error when creating a duplicate table, but got none")
	}
	t.Logf("Expected error received: %v", err)
}

func TestClientDeleteTable(t *testing.T) {
	masterAddress := "localhost:9090"

	client, err := bigtable.NewClient(masterAddress)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	tableName := "users"

	err = client.DeleteTable(tableName)
	if err != nil {
		t.Fatalf("Failed to delete table: %v", err)
	}

	t.Logf("Table '%s' deleted successfully.", tableName)
}

func TestClientWrite(t *testing.T) {
	masterAddress := "localhost:9090"

	client, err := bigtable.NewClient(masterAddress)
	if err != nil {
		log.Fatal(fmt.Printf("Failed to create client: %v", err))
	}

	tableName := "users"
	columnFamilies := map[string][]string{
		"profile":  {"name", "email", "phone", "age"},
		"activity": {"last_login", "last_post", "last_comment"},
	}

	err = client.CreateTable(tableName, columnFamilies)
	if err != nil {
		log.Fatal(fmt.Printf("Failed to create table: %v", err))
	}

	columnFamily := "profile"
	columnName := "name"

	err = client.Write(tableName, "staff_01", columnFamily, columnName, []byte("Xiao"), time.Now().UnixNano())

	if err != nil {
		log.Fatal(fmt.Printf("Failed to create table: %v", err))
	}
}

func TestClientRead(t *testing.T) {
	masterAddress := "localhost:9090"

	client, err := bigtable.NewClient(masterAddress)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	tableName := "users"
	columnFamilies := map[string][]string{
		"profile":  {"name", "email", "phone", "age"},
		"activity": {"last_login", "last_post", "last_comment"},
	}

	err = client.DeleteTable(tableName)
	if err != nil {
		t.Fatalf("Failed to delete table: %v", err)
	}

	// Create table
	err = client.CreateTable(tableName, columnFamilies)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	defer func() {
		err := client.DeleteTable(tableName)
		if err != nil {
			t.Fatalf("Failed to delete table during cleanup: %v", err)
		}
	}()

	// Write data
	columnFamily := "profile"
	columnName := "name"
	rowKey := "staff_01"
	expectedValue := "Xiao"
	timestamp := time.Now().UnixNano()

	err = client.Write(tableName, rowKey, columnFamily, columnName, []byte(expectedValue), timestamp)
	if err != nil {
		t.Fatalf("Failed to write data: %v", err)
	}

	// Read data
	values, err := client.Read(tableName, rowKey, columnFamily, columnName, 1)
	if err != nil {
		t.Fatalf("Failed to read data: %v", err)
	}

	if len(values) == 0 {
		t.Fatalf("No values returned for row key '%s'", rowKey)
	}

	if string(values[0].Value) != expectedValue {
		t.Errorf("Expected value '%s', but got '%s'", expectedValue, values[0].Value)
	} else {
		t.Logf("Read value '%s' successfully.", expectedValue)
	}
}

func TestClientDeleteRow(t *testing.T) {
	masterAddress := "localhost:9090"

	client, err := bigtable.NewClient(masterAddress)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	tableName := "users"
	columnFamilies := map[string][]string{
		"profile":  {"name", "email", "phone", "age"},
		"activity": {"last_login", "last_post", "last_comment"},
	}

	// Create table
	err = client.CreateTable(tableName, columnFamilies)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	defer func() {
		err := client.DeleteTable(tableName)
		if err != nil {
			t.Fatalf("Failed to delete table during cleanup: %v", err)
		}
	}()

	// Write data
	columnFamily := "profile"
	columnName := "name"
	rowKey := "staff_02"
	expectedValue := "Li"
	timestamp := time.Now().UnixNano()

	err = client.Write(tableName, rowKey, columnFamily, columnName, []byte(expectedValue), timestamp)
	if err != nil {
		t.Fatalf("Failed to write data: %v", err)
	}

	// Delete row
	err = client.Delete(tableName, rowKey, columnFamily, columnName)
	if err != nil {
		t.Fatalf("Failed to delete row: %v", err)
	}

	// Try to read deleted row
	_, err = client.Read(tableName, rowKey, columnFamily, columnName, 1)
	if err == nil {
		t.Fatalf("Expected error when reading deleted row key '%s', but got none", rowKey)
	}
	t.Logf("Expected error received when reading deleted row key '%s': %v", rowKey, err)
}

func TestClientReadNonExistingRow(t *testing.T) {
	masterAddress := "localhost:9090"

	client, err := bigtable.NewClient(masterAddress)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	tableName := "users"
	columnFamilies := map[string][]string{
		"profile":  {"name", "email", "phone", "age"},
		"activity": {"last_login", "last_post", "last_comment"},
	}

	// Create table
	err = client.CreateTable(tableName, columnFamilies)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	defer func() {
		err := client.DeleteTable(tableName)
		if err != nil {
			t.Fatalf("Failed to delete table during cleanup: %v", err)
		}
	}()

	// Try to read non-existing row
	rowKey := "non_existing_row"
	columnFamily := "profile"
	columnName := "name"

	_, err = client.Read(tableName, rowKey, columnFamily, columnName, 1)
	if err == nil {
		t.Fatalf("Expected error when reading non-existing row key '%s', but got none", rowKey)
	}
	t.Logf("Expected error received when reading non-existing row key '%s': %v", rowKey, err)
}

func TestClientConcurrentWrites(t *testing.T) {
	masterAddress := "localhost:9090"

	client, err := bigtable.NewClient(masterAddress)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	tableName := "concurrent_users"
	columnFamilies := map[string][]string{
		"profile": {"name", "email", "phone", "age"},
	}

	// 创建表
	err = client.CreateTable(tableName, columnFamilies)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	defer func() {
		err := client.DeleteTable(tableName)
		if err != nil {
			t.Fatalf("Failed to delete table during cleanup: %v", err)
		}
	}()

	numGoroutines := 100
	numWritesPerGoroutine := 10

	done := make(chan bool)

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < numWritesPerGoroutine; j++ {
				rowKey := fmt.Sprintf("user_%d_%d", id, j)
				columnFamily := "profile"
				columnName := "name"
				value := fmt.Sprintf("User%d_%d", id, j)
				timestamp := time.Now().UnixNano()

				err := client.Write(tableName, rowKey, columnFamily, columnName, []byte(value), timestamp)
				if err != nil {
					t.Errorf("Failed to write data for row key '%s': %v", rowKey, err)
				}
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Validate data
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < numWritesPerGoroutine; j++ {
			rowKey := fmt.Sprintf("user_%d_%d", i, j)
			columnFamily := "profile"
			columnName := "name"
			expectedValue := fmt.Sprintf("User%d_%d", i, j)

			values, err := client.Read(tableName, rowKey, columnFamily, columnName, 1)
			if err != nil {
				t.Errorf("Failed to read data for row key '%s': %v", rowKey, err)
				continue
			}

			if len(values) == 0 {
				t.Errorf("No values returned for row key '%s'", rowKey)
				continue
			}

			if string(values[0].Value) != expectedValue {
				t.Errorf("Expected value '%s' for row key '%s', but got '%s'", expectedValue, rowKey, values[0].Value)
			}
		}
	}
	t.Logf("Concurrent writes and reads completed successfully.")
}
