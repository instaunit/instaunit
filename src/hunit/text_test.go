package hunit

import (
  "fmt"
  "testing"
  "github.com/stretchr/testify/assert"
)

var (
  invalidInputError = fmt.Errorf("Invalid input")
)

var context = map[string]interface{}{
  "a": 123,
  "b": "String value",
}

/**
 * Test interpolate
 */
func TestInterpolateVariables(t *testing.T) {
  testInterpolate(t, "This is the string ${a}, alright.", "This is the string 123, alright.", context)
  testInterpolate(t, "This is the string ${a, alright.", invalidInputError, context)
  testInterpolate(t, `This is the string ${a
}, alright.`, `This is the string 123, alright.`, context) // whitespace is not significant in EPL
}

func testInterpolate(t *testing.T, s string, e, c interface{}) {
  t.Logf("----> %v\n", s)
  
  v, err := interpolate(s, "${", "}", c)
  if e == invalidInputError {
    if assert.NotNil(t, err, "Expect an error") {
      t.Logf("(err) %v\n", err)
    }
    return
  }else{
    if !assert.Nil(t, err, fmt.Sprintf("%v", err)) {
      return
    }
  }
  
  t.Logf("<---- %v\n", v)
  if !assert.Equal(t, e, v) {
    return
  }
  
}
