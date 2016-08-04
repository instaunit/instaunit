package hunit

import (
  "fmt"
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
  m := fmt.Sprintf("%v (expected) != %v (actual)", e.Expected, e.Actual)
  if e.Message != "" {
    m += ": "+ e.Message
  }
  return m
}
