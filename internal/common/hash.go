package common

import (
	"crypto/rand"
	"math/big"
)

const charset = "ABCDEFGHIJK0123456789LMNOPQRSTUVWXYZ0123456789"

func GenerateHash() string {
	result := make([]byte, 12)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}
