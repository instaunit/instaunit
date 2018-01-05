package text

import (
  "fmt"
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestStringer(t *testing.T) {
  v := "Hello"
  assert.Equal(t, "Hello", Stringer(v))
  assert.Equal(t, "Hello", Stringer(&v))
  assert.Equal(t, "1", Stringer(1))
  assert.Equal(t, "1234567891", Stringer(1234567891))
  assert.Equal(t, "12345.67891", Stringer(12345.67891))
  assert.Equal(t, fmt.Sprintf("%v", TestStringer), Stringer(TestStringer))
}
