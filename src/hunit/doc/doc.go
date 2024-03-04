package doc

import (
	"fmt"
	"net/http"

	"github.com/instaunit/instaunit/hunit/doc/emit"
	"github.com/instaunit/instaunit/hunit/doc/emit/instadoc"
	"github.com/instaunit/instaunit/hunit/doc/emit/markdown"
	"github.com/instaunit/instaunit/hunit/doc/emit/openapi"
	"github.com/instaunit/instaunit/hunit/test"
)

// Implemented by documentation generators
type Generator interface {
	Init(*test.Suite, string) error
	Case(*test.Suite, test.Case, *http.Request, string, *http.Response, []byte) error
	Finalize(*test.Suite) error
	Close() error
}

// Create a documentation emitter
func New(t emit.Doctype, base string) (Generator, error) {
	switch t {
	case emit.DoctypeMarkdown:
		return Generator(markdown.New(base)), nil
	case emit.DoctypeInstadoc:
		return Generator(instadoc.New(base)), nil
	case emit.DoctypeOpenAPI:
		return Generator(openapi.New(base)), nil
	default:
		return nil, fmt.Errorf("Unsupported doctype: %v", t)
	}
}
