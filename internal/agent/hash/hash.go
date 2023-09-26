package hash

import (
	"crypto/hmac"
	"fmt"
	"hash"
)

// Создаем хеш для его передачи на сервер.
func CrateHash(secretKey string, data []byte, hashNew func() hash.Hash) string {
	secretKeyToByte := []byte(secretKey)
	h := hmac.New(hashNew, secretKeyToByte)
	h.Write(data)
	// вычисляем хеш
	dst := h.Sum(nil)
	hash := fmt.Sprintf("%x", dst)
	return hash
}
