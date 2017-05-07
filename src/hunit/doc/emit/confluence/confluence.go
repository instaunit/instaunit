package confluence

import (
  "io"
  "fmt"
  // "bytes"
  "net/http"
  "hunit/test"
  // "hunit/text"
)

/**
 * A confluence documentation generator
 */
type Generator struct {
  w io.Writer
}

/**
 * Produce a new emitter
 */
func New(w io.Writer) *Generator {
  return &Generator{w}
}

/**
 * Generate documentation preamble
 */
func (g *Generator) Prefix() error {
  return nil
}

/**
 * Generate documentation suffix
 */
func (g *Generator) Suffix() error {
  return nil
}

/**
 * Generate documentation
 */
func (g *Generator) Generate(c test.Case, req *http.Request, reqdata string, rsp *http.Response, rspdata []byte) error {
  var err error
  var doc string
  
  doc += "Documentation goes here..."
  
  _, err = fmt.Fprint(g.w, doc +"\n\n")
  if err != nil {
    return err
  }
  
  return nil
}
