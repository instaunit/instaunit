package edf

import (
  "io"
  "net/http"
  "hunit/test"
)

/**
 * An (HUnit) Endpoint Description Format generator
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
 * Init a suite
 */
func (g *Generator) Init(suite *test.Suite) error {
  return nil
}

/**
 * Finish a suite
 */
func (g *Generator) Finalize(suite *test.Suite) error {
  return nil
}

/**
 * Generate documentation
 */
func (g *Generator) Case(conf test.Config, c test.Case, req *http.Request, reqdata string, rsp *http.Response, rspdata []byte) error {
  return nil
}
