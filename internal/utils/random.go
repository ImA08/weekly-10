package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomString returns a URL-safe, base64 encoded
// random string with the specified length.
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
