package hunit

import (
  "github.com/davecgh/go-spew/spew"
)

/**
 * An assertion error
 */
type AssertionError struct {
  Expected    interface{}
  Actual      interface{}
  Message     string
}

/**
 * Error
 */
func (e AssertionError) Error() string {
  m := e.Message
  if e.Message != "" {
    m += ":\n"
  }
  m += "           expected: "+ spew.Sdump(e.Expected)
  m += "             actual: "+ spew.Sdump(e.Actual)
  return m
}
