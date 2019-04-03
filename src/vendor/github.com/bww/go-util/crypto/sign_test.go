package crypto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// A message
type message struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// Test signing
func TestSign(t *testing.T) {

	msg := []byte("This is the message")
	key := GenerateKey("c2870e3f55ca5acfd9b49aa8c967c6205f1ede3dc84dce60199f75c4a7ff13774cc845c5d2f3afe79b397f0f0ee0d292cc09b78ee6b2590d70db7fcb445efd91", "Salt", SHA1)
	sig := Sign(key, SHA256, msg)
	assert.Equal(t, true, Verify(key, SHA256, sig, msg))

	key = GenerateKey("c2870e3f55ca5acfd9b49aa8c967c6205f1ede3dc84dce60199f75c4a7ff13774cc845c5d2f3afe79b397f0f0ee0d292cc09b78ee6b2590d70db7fcb445efd91", "Salt", SHA1)
	sig = Sign(key, SHA256, []byte("This is the message"))
	assert.Equal(t, false, Verify(key, SHA256, sig, []byte("this is the message")))

	key1 := GenerateKey("c2870e3f55ca5acfd9b49aa8c967c6205f1ede3dc84dce60199f75c4a7ff13774cc845c5d2f3afe79b397f0f0ee0d292cc09b78ee6b2590d70db7fcb445efd91", "Salt A", SHA1)
	key2 := GenerateKey("c2870e3f55ca5acfd9b49aa8c967c6205f1ede3dc84dce60199f75c4a7ff13774cc845c5d2f3afe79b397f0f0ee0d292cc09b78ee6b2590d70db7fcb445efd91", "Salt B", SHA1)
	sig = Sign(key1, SHA256, msg)
	assert.Equal(t, false, Verify(key2, SHA256, sig, msg))

	key1 = GenerateKey("a2870e3f55ca5acfd9b49aa8c967c6205f1ede3dc84dce60199f75c4a7ff13774cc845c5d2f3afe79b397f0f0ee0d292cc09b78ee6b2590d70db7fcb445efd91", "Salt", SHA1)
	key2 = GenerateKey("b2870e3f55ca5acfd9b49aa8c967c6205f1ede3dc84dce60199f75c4a7ff13774cc845c5d2f3afe79b397f0f0ee0d292cc09b78ee6b2590d70db7fcb445efd91", "Salt", SHA1)
	sig = Sign(key1, SHA256, msg)
	assert.Equal(t, false, Verify(key2, SHA256, sig, msg))

}

// Test message signing
func TestSignMessage(t *testing.T) {
	key := GenerateKey("c2870e3f55ca5acfd9b49aa8c967c6205f1ede3dc84dce60199f75c4a7ff13774cc845c5d2f3afe79b397f0f0ee0d292cc09b78ee6b2590d70db7fcb445efd91", "Salt", SHA1)
	var rsp message

	msg := message{"First", 100}
	enc, sig, err := SignMessage(key, SHA256, msg)
	if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
		err := VerifyMessage(key, SHA256, sig, &rsp, enc)
		if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			assert.Equal(t, msg, rsp)
		}
	}

	msg = message{"Second", 200}
	inv := GenerateKey("c2870e3f55ca5acfd9b49aa8c967c6205f1ede3dc84dce60199f75c4a7ff13774cc845c5d2f3afe79b397f0f0ee0d292cc09b78ee6b2590d70db7fcb445efd91", "Another Salt", SHA1)
	enc, sig, err = SignMessage(key, SHA256, msg)
	if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
		err := VerifyMessage(inv, SHA256, sig, &rsp, enc)
		if assert.NotNil(t, err, fmt.Sprintf("%v", err)) {
			assert.Equal(t, ErrSignatureNotValid, err)
		}
	}

}
