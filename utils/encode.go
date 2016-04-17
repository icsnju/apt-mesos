package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func Encode(size int) (string, error) {
	id := make([]byte, size)
	n, err := rand.Read(id)
	if n != len(id) || err != nil {
		return "", err
	}
	return hex.EncodeToString(id), nil
}
