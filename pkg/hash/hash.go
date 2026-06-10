// Package hash provides lightweight, stateless utility routines to calculate
// deterministic checksum signatures for password string verification.
package hash

import (
	"fmt"
	"hash/crc32"
)

// HashPassword computes an 8-character, zero-padded hexadecimal representation of the
// provided plaintext string using the CRC-32 IEEE 802.3 cyclic redundancy polynomial.
func HashPassword(expectedPassword string) string {
	return fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(expectedPassword)))
}
