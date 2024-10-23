package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}
