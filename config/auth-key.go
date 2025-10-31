package config

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

// GenerateRandomKey creates a cryptographically secure random key
// It generates a 32-byte random key and encodes it in base64 URL encoding
// This key is used for signing JWT tokens
// Returns:
//   - string: Base64-encoded random key (256 bits of entropy)
func GenerateRandomKey() string {
	// Create a 32-byte (256-bit) slice for the random key
	bytes := make([]byte, 32)

	// Fill the slice with cryptographically secure random bytes
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal("Failed to generate random key:", err)
	}

	// Encode the random bytes to base64 URL-safe string
	return base64.URLEncoding.EncodeToString(bytes)
}
