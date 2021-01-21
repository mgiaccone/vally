package hashutil

import (
	"crypto/sha1"
	"encoding/hex"
)

// Sha1Hex computes the SHA-1 hash for the given string and returns it in hex format.
func Sha1Hex(in string) string {
	h := sha1.New()
	h.Write([]byte(in))
	return hex.EncodeToString(h.Sum(nil))
}
