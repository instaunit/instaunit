package hunit

import (
  "reflect"
)

import (
  "github.com/davecgh/go-spew/spew"
  "github.com/pmezard/go-difflib/difflib"
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

  et, ek := typeAndKind(e.Expected)
  at, _ := typeAndKind(e.Actual)
  
  if et != at || (ek != reflect.String && ek != reflect.Struct && ek != reflect.Map && ek != reflect.Slice && ek != reflect.Array) {
    m += "           expected: "+ spew.Sdump(e.Expected)
    m += "             actual: "+ spew.Sdump(e.Actual)
  }else{
    d := diff(e.Expected, e.Actual)
    if d != "" {
      m += indent(d, "           ")
    }
  }
  
  return m
}

/**
 * Indent a string by prefixing each line with the provided string
 */
func indent(s, p string) string {
  o := p
  for i := 0; i < len(s); i++ {
    o += string(s[i])
    if s[i] == '\n' {
      o += p
    }
  }
  return o
}

/**
 * Obtain a diff of values
 */
func diff(expected, actual interface{}) string {
  if expected == nil || actual == nil {
    return ""
  }
  
  et, ek := typeAndKind(expected)
  at, _ := typeAndKind(actual)
  
  if et != at {
    return ""
  }
  if ek != reflect.String && ek != reflect.Struct && ek != reflect.Map && ek != reflect.Slice && ek != reflect.Array {
    return ""
  }
  
  var e, a string
  if ek == reflect.String {
    e = expected.(string)
    a = actual.(string)
  }else{
    spew.Config.SortKeys = true
    e = spew.Sdump(expected)
    a = spew.Sdump(actual)
  }
  
  diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
    A:        difflib.SplitLines(e),
    B:        difflib.SplitLines(a),
    FromFile: "Expected",
    ToFile:   "Actual",
    Context:  1,
  })
  
  return diff
}

/**
 * Obtain a value's type and kind
 */
func typeAndKind(v interface{}) (reflect.Type, reflect.Kind) {
  t := reflect.TypeOf(v)
  k := t.Kind()
  if k == reflect.Ptr {
    t = t.Elem()
    k = t.Kind()
  }
  return t, k
}
