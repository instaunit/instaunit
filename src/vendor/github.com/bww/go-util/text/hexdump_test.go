package text

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test hexdump string
func TestHexdump(t *testing.T) {
	testHexdump(t, "Hello", "48 65 6c 6c 6f    | Hello", 6)
	testHexdump(t, "Hello", "48 65 | He\n6c 6c | ll\n6f    | o", 2)
	testHexdump(t, "Hello", "48 | H\n65 | e\n6c | l\n6c | l\n6f | o", 1)
	testHexdump(t, "Hello", "48 | H\n65 | e\n6c | l\n6c | l\n6f | o", 0)
	testHexdump(t, "Ok! 👏👏👏", "4f 6b 21 20 f0 9f 91 8f f0 9f | Ok! ......\n91 8f f0 9f 91 8f             | ......", 10)
}

func testHexdump(t *testing.T, v, e string, w int) {
	b := &bytes.Buffer{}
	Hexdump(b, []byte(v), w)
	fmt.Printf(">>>\n%s\n", string(b.Bytes()))
	assert.Equal(t, e, string(b.Bytes()))
}
