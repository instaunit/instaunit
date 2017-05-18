package runtime

import (
  "github.com/bww/go-util/rand"
)

/**
 * The standard library
 */
type stdlib struct{}

/**
 * Builtins
 */
var Stdlib stdlib

/**
 * Generate a random string
 */
func (s stdlib) RandomString(n float64) string {
  return rand.RandomString(int(n))
}
