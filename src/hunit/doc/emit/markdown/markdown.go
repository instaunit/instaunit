package markdown

import (
  "io"
  "fmt"
  "hunit/test"
)

/**
 * A markdown documentation generator
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
 * Generate documentation
 */
func (g *Generator) Generate(c test.Case, req string, rsp []byte) error {
  var err error
  
  _, err = fmt.Fprintf(g.w, "## %s %s\n", c.Request.Method, c.Request.URL)
  if err != nil {
    return err
  }
  
  if c.Comments != "" {
    _, err = fmt.Fprint(g.w, c.Comments)
    if err != nil {
      return err
    }
  }
  
  _, err = fmt.Fprint(g.w, "\n\n")
  if err != nil {
    return err
  }
  
  return nil
}
