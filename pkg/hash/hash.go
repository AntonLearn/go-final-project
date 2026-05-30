// Package hash
package hash

import (
	"fmt"
	"hash/crc32"
)

func HashPassword(expectedPassword string) string {
	return fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(expectedPassword)))
}
