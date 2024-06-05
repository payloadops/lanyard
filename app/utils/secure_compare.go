package utils

import (
	"crypto/subtle"
)

// SecureCompare compares two strings in a constant time to prevent timing attacks.
func SecureCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
