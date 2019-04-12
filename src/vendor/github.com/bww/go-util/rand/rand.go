package rand

import (
	"crypto/rand"
	"fmt"
	"io"
	"net"
)

/**
 * Cached MAC address
 */
var macaddr net.HardwareAddr

/**
 * Content types
 */
const (
	CONTENT_TYPE_JSON = "application/json"
)

/**
 * Setup
 */
func init() {

	// search our network interfaces for a hardware MAC address
	if interfaces, err := net.Interfaces(); err == nil {
		for _, i := range interfaces {
			if (i.Flags&net.FlagLoopback) == 0 && len(i.HardwareAddr) > 0 {
				macaddr = i.HardwareAddr
				break
			}
		}
	}

	// if we failed to obtain the MAC address of the current computer, we will use a randomly generated 6 byte sequence instead and set the multicast bit as recommended in RFC 4122.
	if macaddr == nil {
		macaddr = make(net.HardwareAddr, 6)
		randomBytes(macaddr)
		macaddr[0] = macaddr[0] | 0x01
	}

}

/**
 * Obtain the current host's MAC address
 */
func HardwareAddr() net.HardwareAddr {
	return macaddr
}

/**
 * Obtain the current host's MAC address as a hex string
 */
func HardwareKey() string {
	return fmt.Sprintf("%x", macaddr)
}

/**
 * Generate random bytes
 */
func RandomBytes(n int) []byte {
	b := make([]byte, n)
	randomBytes(b)
	return b
}

/**
 * Read random bytes
 */
func ReadRandom(b []byte) {
	randomBytes(b)
}

/**
 * Generate random bytes
 */
func randomBytes(b []byte) {
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err.Error()) // rand should never fail
	}
}

/**
 * Characters used in random strings
 */
const randomStringChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

/**
 * Generate a random string representing the specified number of characters
 */
func RandomString(n int) string {
	b := RandomBytes(n)
	s := ""
	l := len(randomStringChars)
	for i := 0; i < len(b); i++ {
		s += string(randomStringChars[int(b[i])%l])
	}
	return s
}
