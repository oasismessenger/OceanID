package utils

import (
	"crypto/rand"
	"math/big"
)

func Int64(max int64) int64 {
	v, _ := rand.Int(rand.Reader, big.NewInt(max-1))
	return v.Int64()
}
