// package token provides functions to create random tokens
package token

import (
	// With the length constraint given (10) it doesn't really make sense to use e CSPRNG
	// so we're going with this for simplicity
	"math/rand"
	"time"
)

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// New returns a random token of the specified length, panics on negative length
func New(n int) string {
	token := make([]byte, 0, n)

	for i := 0; i < n; i++ {
		token = append(
			token,
			charSet[rand.Intn(len(charSet))],
		)
	}

	return string(token)
}
