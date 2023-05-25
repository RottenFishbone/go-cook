package config

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

// Generates a random key of `length`. Assumes SHA-512 is used.
func generateHMACKey(length int) string {
	key := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		randVal, _ := rand.Int(rand.Reader, big.NewInt(256))
		randByte := byte(randVal.Int64())
		key = append(key, randByte)
	}

	return hex.EncodeToString(key)
}

// Returns the HMAC secret key as a byte slice
//
// Panics if config has not been loaded
func GetHMACKeyBytes() []byte {
	if !loaded {
		panic("Attempted to read HMACKey before config was loaded.")
	}

	data, _ := hex.DecodeString(conf.HMACKey)
	return data
}
