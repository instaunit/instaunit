package doc

import (
	"fmt"
	"io"
	"net/http"

	"github.com/instaunit/instaunit/hunit/doc/emit"
	"github.com/instaunit/instaunit/hunit/doc/emit/markdown"
	"github.com/instaunit/instaunit/hunit/test"
)

// Implemented by documentation generators
type Generator interface {
	Init(*test.Suite) error
	Case(*test.Suite, test.Case, *http.Request, string, *http.Response, []byte) error
	Finalize(*test.Suite) error
	Close() error
}

// Create a documentation emitter
func New(t emit.Doctype, o io.WriteCloser) (Generator, error) {
	switch t {
	case emit.DoctypeMarkdown:
		return Generator(markdown.New(o)), nil
	default:
		return nil, fmt.Errorf("Unsupported doctype: %v", t)
	}
}
