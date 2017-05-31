package expr

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
  "c": map[string]interface{}{
    "a": []string{ "Zero", "One", "Two", "Three" },
    "b": false,
  },
}

/**
 * Test interpolate
 */
func TestInterpolateVariables(t *testing.T) {
  testInterpolate(t, `Before ${a}, after.`, `Before 123, after.`, context)
  testInterpolate(t, `Before \${a}, after.`, `Before ${a}, after.`, context)
  testInterpolate(t, `Before \\${a}, after.`, `Before \123, after.`, context)
  testInterpolate(t, `Before $${a}}, after.`, `Before $123}, after.`, context)
  testInterpolate(t, `Before $${a}}, after.`, `Before $123}, after.`, context)
  testInterpolate(t, `Before ${c.a[0]}, after.`, `Before Zero, after.`, context)
  testInterpolate(t, `Before ${c["b"]}, after.`, `Before false, after.`, context)
  testInterpolate(t, `Before ${a, after.`, invalidInputError, context)
  testInterpolate(t, `Before ${a
}, after.`, `Before 123, after.`, context) // whitespace is not significant in EPL
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
