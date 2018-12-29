package common

import (
	"math/rand"
	"time"
)

var r *rand.Rand // Rand for this package.

// check if applicable or should be placed inside random string
func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomString generates a random string of strlen length
func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := ""
	for i := 0; i < strlen; i++ {
		index := r.Intn(len(chars))
		result += chars[index : index+1]
	}
	return result
}
