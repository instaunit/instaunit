package debug

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDump(t *testing.T) {

	type Type struct {
		A string `json:"a"`
		B int    `json:"b"`
	}

	b := &bytes.Buffer{}
	Dumpf(b, Type{"Hello, there", 987})
	assert.Equal(t, "{\n  \"a\": \"Hello, there\",\n  \"b\": 987\n}\n", string(b.Bytes()))

}
