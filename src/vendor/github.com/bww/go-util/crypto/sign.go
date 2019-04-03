package crypto

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

var ErrSignatureNotValid = fmt.Errorf("Signature not valid")

// Sign data using HMAC
func Sign(key []byte, digest Digest, data []byte) string {
	mac := hmac.New(digest.Func(), key)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify a signature
func Verify(key []byte, digest Digest, sig string, data []byte) bool {
	b, err := hex.DecodeString(sig)
	if err != nil {
		return false
	}
	mac := hmac.New(digest.Func(), key)
	mac.Write(data)
	return hmac.Equal(b, mac.Sum(nil))
}

// Produce a serialized, base64-encoded message and its signature
func SignMessage(key []byte, digest Digest, message interface{}) (string, string, error) {

	m, err := json.Marshal(message)
	if err != nil {
		return "", "", err
	}

	b := base64.StdEncoding.EncodeToString(m)
	s := Sign(key, digest, []byte(b))

	return b, s, nil
}

// Validate a message signature and, if valid, unmarshal the message
func VerifyMessage(key []byte, digest Digest, sig string, message interface{}, data string) error {
	if !Verify(key, digest, sig, []byte(data)) {
		return ErrSignatureNotValid
	}

	m, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(m, message)
	if err != nil {
		return err
	}

	return nil
}
