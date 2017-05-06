package markdown

import (
  "io"
  "fmt"
  "hunit/test"
  "hunit/text"
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
    _, err = fmt.Fprint(g.w, "\n"+ c.Comments)
    if err != nil {
      return err
    }
  }
  
  if req != "" {
    _, err = fmt.Fprint(g.w, "\n"+ text.Indent(req, "    "))
    if err != nil {
      return err
    }
  }
  
  if len(rsp) > 0 {
    _, err = fmt.Fprint(g.w, "\n"+ text.Indent(string(rsp), "    "))
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
