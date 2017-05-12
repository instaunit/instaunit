package markdown

import (
  "io"
  "fmt"
  "bytes"
  "strings"
  "net/http"
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
 * Generate documentation preamble
 */
func (g *Generator) Prefix(suite *test.Suite) error {
  var err error
  var doc string
  
  if suite.Title != "" {
    doc += fmt.Sprintf("# %s\n\n", strings.TrimSpace(suite.Title))
  }
  if suite.Comments != "" {
    doc += strings.TrimSpace(suite.Comments) +"\n\n"
  }
  
  _, err = fmt.Fprint(g.w, doc)
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * Generate documentation suffix
 */
func (g *Generator) Suffix(suite *test.Suite) error {
  return nil
}

/**
 * Generate documentation
 */
func (g *Generator) Generate(conf test.Config, c test.Case, req *http.Request, reqdata string, rsp *http.Response, rspdata []byte) error {
  var err error
  var doc string
  
  if c.Title != "" {
    doc += fmt.Sprintf("## %s\n\n", strings.TrimSpace(c.Title))
  }else{
    doc += fmt.Sprintf("## %s %s\n\n", c.Request.Method, c.Request.URL)
  }
  
  if c.Comments != "" {
    doc += strings.TrimSpace(c.Comments) +"\n\n"
  }
  
  if req != nil {
    b := &bytes.Buffer{}
    if conf.Doc.FormatEntities {
      f, err := text.FormatEntity(req, []byte(reqdata))
      if err == nil {
        reqdata = string(f)
      }else if err != nil && err != text.ErrUnsupportedContentType {
        fmt.Println("* * * Invalid request entity could not be formatted: %v", req.Header.Get("Content-Type"))
      }
    }
    err = text.WriteRequest(b, req, reqdata)
    if err != nil {
      return err
    }
    if b.Len() > 0 {
      doc += "### Example request\n\n"
      doc += "```http\n"
      doc += string(b.Bytes()) +"\n"
      doc += "```\n"
    }
  }
  
  if rsp != nil {
    b := &bytes.Buffer{}
    if conf.Doc.FormatEntities {
      f, err := text.FormatEntity(req, rspdata)
      if err == nil {
        rspdata = f
      }else if err != nil && err != text.ErrUnsupportedContentType {
        fmt.Println("* * * Invalid entity could not be formatted: %v", req.Header.Get("Content-Type"))
      }
    }
    err = text.WriteResponse(b, rsp, rspdata)
    if err != nil {
      return err
    }
    if b.Len() > 0 {
      doc += "### Example response\n\n"
      doc += "```http\n"
      doc += string(b.Bytes()) +"\n"
      doc += "```\n"
    }
  }
  
  _, err = fmt.Fprint(g.w, doc +"\n\n")
  if err != nil {
    return err
  }
  
  return nil
}
