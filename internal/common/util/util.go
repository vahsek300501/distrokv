package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHash(plaintext string) string {
	hasher := sha256.New()
	hasher.Write([]byte(plaintext))
	return hex.EncodeToString(hasher.Sum(nil))
}
