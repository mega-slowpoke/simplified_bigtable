package tablet

import (
	"fmt"
	"strconv"
	"strings"
)

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
