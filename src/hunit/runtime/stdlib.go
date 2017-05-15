package runtime

import (
  "fmt"
  "strings"
)

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

/**
 * Generate a random name (Docker style: <adjective> <noun>)
 */
func (s stdlib) RandomName() string {
  l, r := RandomDockerName()
  return fmt.Sprintf("%s %s", strings.Title(l), strings.Title(r))
}
