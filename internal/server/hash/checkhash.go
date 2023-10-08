package hash

import (
	"crypto/hmac"
	"errors"
	"fmt"
	"hash"
)

// Проверка хэша при получении запроса
func CheckHash(data []byte, secretKey, checksum string, hashNew func() hash.Hash) error {
	secretKeyToByte := []byte(secretKey)
	h := hmac.New(hashNew, secretKeyToByte)
	h.Write(data)
	// вычисляем хеш
	hash := h.Sum(nil)
	hashString := fmt.Sprintf("%x", hash)
	if hashString != checksum {
		return nil
	}
	return errors.New("discrepancy between received and calculated hash")
}
