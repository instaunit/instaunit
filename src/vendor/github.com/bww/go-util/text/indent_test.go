package text

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

/**
 * Test indent string
 */
func TestIndent(t *testing.T) {
	assert.Equal(t, "> Hello", Indent("Hello", "> "))
	assert.Equal(t, "> Hello\n> There", Indent("Hello\nThere", "> "))
	assert.Equal(t, "> Hello\n> There\n> ", Indent("Hello\nThere\n", "> "))
}

/**
 * Test indent writer
 */
func TestIndentWriter(t *testing.T) {
	b := &bytes.Buffer{}
	w := NewIndentWriter("> ", IndentOptionIndentFirstLine, b)
	io.WriteString(w, "Hello\nThere\nBr")
	io.WriteString(w, "ah\nChillin.")
	assert.Equal(t, "> Hello\n> There\n> Brah\n> Chillin.", string(b.Bytes()))
}
