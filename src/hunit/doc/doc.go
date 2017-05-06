package doc

import (
  "io"
  "fmt"
  "hunit/test"
  "hunit/doc/emit"
  "hunit/doc/emit/markdown"
)

/**
 * Implemented by documentation generators
 */
type Generator interface {
  Generate(test.Case, string, []byte)(error)
}

/**
 * Create a documentation emitter
 */
func New(t emit.Doctype, o io.Writer) (Generator, error) {
  switch t {
    case emit.DoctypeMarkdown:
      return Generator(markdown.New(o)), nil
    default:
      return nil, fmt.Errorf("Unsupported doctype: %v", t)
  }
}
