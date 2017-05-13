package markdown

import (
  "io"
  "fmt"
  "bytes"
  "strings"
  "net/http"
  "hunit/test"
  "hunit/text"
  "hunit/text/slug"
)

/**
 * A markdown documentation generator
 */
type Generator struct {
  w         io.Writer
  b         *bytes.Buffer
  sections  map[string]string
  slugs     map[string]int
  
}

/**
 * Produce a new emitter
 */
func New(w io.Writer) *Generator {
  return &Generator{w, nil, make(map[string]string), nil}
}

/**
 * Init a suite
 */
func (g *Generator) Init(suite *test.Suite) error {
  g.b = &bytes.Buffer{}
  return nil
}

/**
 * Finish a suite
 */
func (g *Generator) Finalize(suite *test.Suite) error {
  var err error
  
  err = g.prefix(g.w, suite)
  if err != nil {
    return err
  }
  
  err = g.contents(g.w, suite)
  if err != nil {
    return err
  }
  
  _, err = g.w.Write(g.b.Bytes())
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * Generate documentation
 */
func (g *Generator) Case(conf test.Config, c test.Case, req *http.Request, reqdata string, rsp *http.Response, rspdata []byte) error {
  return g.generate(g.b, conf, c, req, reqdata, rsp, rspdata)
}

/**
 * Generate documentation preamble
 */
func (g *Generator) prefix(w io.Writer, suite *test.Suite) error {
  var err error
  var doc string
  
  if suite.Title != "" {
    doc += fmt.Sprintf("# %s\n\n", strings.TrimSpace(suite.Title))
  }
  if suite.Comments != "" {
    doc += strings.TrimSpace(suite.Comments) +"\n\n"
  }
  
  _, err = fmt.Fprint(w, doc)
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * Table of contents
 */
func (g *Generator) contents(w io.Writer, suite *test.Suite) error {
  var err error
  var doc string
  
  doc += "## Contents\n\n"
  
  for s, t := range g.sections {
    doc += fmt.Sprintf("* [%s](#%s)\n", strings.TrimSpace(t), s)
  }

  doc += "\n"
  
  _, err = fmt.Fprint(w, doc)
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * Generate documentation
 */
func (g *Generator) generate(w io.Writer, conf test.Config, c test.Case, req *http.Request, reqdata string, rsp *http.Response, rspdata []byte) error {
  var err error
  var doc string
  
  var t string
  if c.Title != "" {
    t = strings.TrimSpace(c.Title)
  }else{
    t = fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL)
  }
  
  doc += fmt.Sprintf("## %s\n\n", t)
  var s string
  s, g.slugs = slug.Github(t, g.slugs)
  g.sections[s] = t
  
  if c.Comments != "" {
    doc += strings.TrimSpace(c.Comments) +"\n\n"
  }
  
  if req != nil {
    b := &bytes.Buffer{}
    if conf.Doc.FormatEntities {
      t := text.Coalesce(c.Request.Format, req.Header.Get("Content-Type"))
      f, err := text.FormatEntity([]byte(reqdata), t)
      if err == nil {
        reqdata = string(f)
      }else if err != nil && err != text.ErrUnsupportedContentType {
        fmt.Println("* * * Invalid request entity could not be formatted: %v", t)
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
      t := text.Coalesce(c.Response.Format, rsp.Header.Get("Content-Type"))
      f, err := text.FormatEntity(rspdata, t)
      if err == nil {
        rspdata = f
      }else if err != nil && err != text.ErrUnsupportedContentType {
        fmt.Println("* * * Invalid entity could not be formatted: %v", t)
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
  
  _, err = fmt.Fprint(g.b, doc +"\n\n")
  if err != nil {
    return err
  }
  
  return nil
}
