package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/subtle"
	"hash"
)

import (
	"golang.org/x/crypto/pbkdf2"
)

const (
	keySize = 64
	keyIter = 65536
)

// hash generator
type HashFunc func() hash.Hash

// Digest algorithms
type Digest int

const (
	SHA1 Digest = iota
	SHA256
)

// hash functions
var hashFuncs = []HashFunc{
	sha1.New,
	sha256.New,
}

// Obtain the digest hash function
func (d Digest) Func() HashFunc {
	if d >= SHA1 && d <= SHA256 {
		return hashFuncs[int(d)]
	}
	return nil
}

// Pad an AES block
func PKCS5Pad(in []byte) []byte {
	pad := aes.BlockSize - (len(in) % aes.BlockSize)
	return append(in, bytes.Repeat([]byte{byte(pad)}, pad)...)
}

// Unpad an AES block
func PKCS5Unpad(in []byte) []byte {
	l := len(in)
	p := int(in[l-1])
	return in[:l-p]
}

// Generate a Rails-compatible key
func GenerateKey(secret, salt string, digest Digest) []byte {
	return GenerateKeyWithOptions(secret, salt, keyIter, keySize, digest.Func())
}

// Generate a key with options
func GenerateKeyWithOptions(secret, salt string, iter, size int, hgen HashFunc) []byte {
	return pbkdf2.Key([]byte(secret), []byte(salt), iter, size, hgen)
}

// Compare strings in constant time
func SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
