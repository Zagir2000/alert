package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func CrateHash(secretKey string, data []byte) string {
	secretKeyToByte := []byte(secretKey)
	h := hmac.New(sha256.New, secretKeyToByte)
	h.Write(data)
	// вычисляем хеш
	dst := h.Sum(nil)
	hash := fmt.Sprintf("%x", dst)
	return hash
}
