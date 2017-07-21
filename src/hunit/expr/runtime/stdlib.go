package runtime

import (
  "fmt"
  "net/url"
  "strings"
)

import (
  "github.com/bww/go-util/rand"
)

// The standard library
type stdlib struct{}

// Builtins
var Stdlib stdlib

// Generate a random string
func (s stdlib) RandomString(n float64) string {
  return rand.RandomString(int(n))
}

// Generate a random name (Docker style: <adjective>_<noun>)
func (s stdlib) RandomIdent() string {
  l, r := DockerName()
  return fmt.Sprintf("%s_%s", l, r)
}

// Generate a random person name
func (s stdlib) RandomName() string {
  n := PersonName(2)
  return fmt.Sprintf("%s %s", strings.Title(n[0]), strings.Title(n[1]))
}

// Escape a URL query component
func (s stdlib) QueryEscape(v string) string {
  return url.QueryEscape(v)
}

// Unescape a URL component
func (s stdlib) QueryUnescape(v string) (string, error) {
  return url.QueryUnescape(v)
}
