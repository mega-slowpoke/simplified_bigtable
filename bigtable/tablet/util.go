package tablet

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type SetupOptions struct {
	TabletAddress string
	MasterAddress string
	MaxTableSize  int
	TestMode      bool
}

// localhost:9090, : is not allowed in the file path, replace it with _
func GetFilePath(address string, tableName string) string {
	address = strings.Replace(address, ":", "_", -1)
	return filepath.Join(address, tableName)
}

func BuildKey(rowKey, columnFamily, columnQualifier string, timestamp int64) string {
	reversedTimestamp := ^timestamp // descendingOrder timestamp so we can retrieve the most recent data first
	return fmt.Sprintf("%s:%s:%s:%d", rowKey, columnFamily, columnQualifier, reversedTimestamp)
}

func ExtractTimestampFromKey(key string) (int64, error) {
	parts := strings.Split(key, ":")
	if len(parts) < 4 {
		return 0, fmt.Errorf("invalid key format") // Return an error if the key is in an invalid format
	}

	// Extract the reversed timestamp (the 4th part)
	reversedTimestampStr := parts[3]
	reversedTimestamp, err := strconv.ParseInt(reversedTimestampStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse timestamp: %v", err) // Return an error if parsing fails
	}
	timestamp := ^reversedTimestamp
	return timestamp, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
