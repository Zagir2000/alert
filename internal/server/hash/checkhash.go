package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
)

func CheckHash(data []byte, secretKey, checksum string) error {
	secretKeyToByte := []byte(secretKey)
	h := hmac.New(sha256.New, secretKeyToByte)
	h.Write(data)
	// вычисляем хеш
	hash := h.Sum(nil)
	hashString := fmt.Sprintf("%x", hash)
	fmt.Println(fmt.Sprintf("%x", checksum))
	if hashString == checksum {
		return nil
	}
	return errors.New("discrepancy between received and calculated hash")
}
