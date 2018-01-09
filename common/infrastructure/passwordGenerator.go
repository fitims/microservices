package infrastructure

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!%&[]&@$"

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// GeneratePassword generates a random password
func GeneratePassword(length int) string {

	l := len(charset)
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rnd.Intn(l)]
	}

	return string(b)
}
