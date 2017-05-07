package doc

import (
  "io"
  "fmt"
  "net/http"
  "hunit/test"
  "hunit/doc/emit"
  "hunit/doc/emit/markdown"
  "hunit/doc/emit/confluence"
)

/**
 * Implemented by documentation generators
 */
type Generator interface {
  Prefix()(error)
  Suffix()(error)
  Generate(test.Case, *http.Request, string, *http.Response, []byte)(error)
}

/**
 * Create a documentation emitter
 */
func New(t emit.Doctype, o io.Writer) (Generator, error) {
  switch t {
    case emit.DoctypeMarkdown:
      return Generator(markdown.New(o)), nil
    case emit.DoctypeConfluence:
      return Generator(confluence.New(o)), nil
    default:
      return nil, fmt.Errorf("Unsupported doctype: %v", t)
  }
}
