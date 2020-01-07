package text

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndent(t *testing.T) {
	prefix := "> "
	lineno := `{{ printf "%3d: " .Line }}`
	tests := []struct {
		Text    string
		Options IndentOptions
		Prefix  string
		Expect  string
	}{
		{
			"This is text",
			IndentOptionIndentFirstLine,
			prefix,
			fmt.Sprintf("%sThis is text", prefix),
		},
		{
			"This is text\nAnd this is more",
			IndentOptionIndentFirstLine,
			prefix,
			fmt.Sprintf("%sThis is text\n%sAnd this is more", prefix, prefix),
		},
		{
			"This is text\nAnd this is more",
			IndentOptionIndentFirstLine | IndentOptionIndentTemplate,
			lineno,
			"  1: This is text\n  2: And this is more",
		},
	}
	for _, e := range tests {
		v := IndentWithOptions(e.Text, e.Prefix, e.Options)
		fmt.Println(v)
		assert.Equal(t, e.Expect, v, e.Text)
	}
}
