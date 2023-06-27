package snowflake

import (
	"fmt"
	"testing"
)

// TestGetUUID
func TestGetUUID(t *testing.T) {
	fmt.Println(GetUUID())
}

// TestGetSnowflakeID
func TestGetSnowflakeID(t *testing.T) {
	for i := int(0); i < 100; i++ {
		fmt.Println(GetSnowflakeID())
	}
}
