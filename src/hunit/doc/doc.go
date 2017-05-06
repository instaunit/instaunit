package doc

import (
  "fmt"
  "hunit/test"
  "hunit/doc/emit"
  "hunit/doc/emit/markdown"
)

/**
 * Implemented by documentation generators
 */
type Generator interface {
  Generate(io.Writer, test.Case)(error)
}

/**
 * Create a documentation emitter
 */
func New(t emit.Doctype) (Generator, error) {
  switch t {
    case emit.Markdown:
      Generator(markdown.New()), nil
    default:
      return nil, fmt.Errorf("Unsupported doctype: %v", t)
  }
}
