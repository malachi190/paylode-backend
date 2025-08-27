package config

import (
	"crypto/rand"
	"math/big"
)

func GenerateOtp() int {
	n, _ := rand.Int(rand.Reader, big.NewInt(90000))
	return int(n.Int64() + 10000)
}
