package crypto

import (
	"github.com/bww/go-util/rand"
)

// Generate a Devise-compatible general-purpose secure token
func GenerateToken(key []byte, digest Digest) (string, string) {
	t := rand.RandomString(32)
	return t, Sign(key, digest, []byte(t))
}
