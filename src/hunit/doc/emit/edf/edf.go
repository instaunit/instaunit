package edf

import (
  "io"
  "net/http"
  "hunit/test"
)

import (
  "gopkg.in/yaml.v2"
)

/**
 * An (HUnit) Endpoint Description Format generator. This is essentially
 * just the test suite format with all values literalized, and certain
 * features excluded.
 */
type Generator struct {
  w     io.Writer
  suite *test.Suite
}

/**
 * Produce a new emitter
 */
func New(w io.Writer) *Generator {
  return &Generator{w, nil}
}

/**
 * Init a suite
 */
func (g *Generator) Init(suite *test.Suite) error {
  g.suite = &test.Suite{
    Title:    suite.Title,
    Comments: suite.Comments,
    Config:   suite.Config,
  }
  return nil
}

/**
 * Finish a suite
 */
func (g *Generator) Finalize(suite *test.Suite) error {
  
  data, err := yaml.Marshal(g.suite)
  if err != nil {
    return err
  }
  
  _, err = g.w.Write(data)
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * Generate documentation
 */
func (g *Generator) Case(conf test.Config, c test.Case, req *http.Request, reqdata string, rsp *http.Response, rspdata []byte) error {
  c.Gendoc = c.Documented() // literalize
  c.Request.BasicAuth = nil // clear credentials
  c.Response.Entity = string(rspdata)
  c.Response.Headers = make(map[string]string)
  
  for k, v := range rsp.Header {
    var x string
    if len(v) > 0 {
      x = v[0]
    }
    c.Response.Headers[k] = x
  }
  
  g.suite.Cases = append(g.suite.Cases, c)
  return nil
}
