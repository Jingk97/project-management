package common

import (
	"crypto/rand"
	"math/big"
)

const randomCodePool = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateCode(length int) (string, error) {
	var result = make([]byte, length)
	maxBit := big.NewInt(int64(len(randomCodePool)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, maxBit)
		if err != nil {
			return "", err
		}
		result[i] = randomCodePool[n.Int64()]
	}

	return string(result), nil
}
