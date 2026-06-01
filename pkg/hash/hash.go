// Package hash provides utility functions for password hashing.
package hash

import (
	"fmt"
	"hash/crc32"
)

// HashPassword creates a simple hash of the password using CRC32.
// Returns an 8-character hexadecimal string.
func HashPassword(expectedPassword string) string {
	return fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(expectedPassword)))
}
