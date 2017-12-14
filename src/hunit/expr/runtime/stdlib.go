package runtime

import (
  "fmt"
  "net/url"
  "strings"
)

import (
  "github.com/bww/go-util/rand"
  "github.com/bww/go-util/uuid"
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
func (s stdlib) RandomPersonName() string {
  return fmt.Sprintf("%s %s", strings.Title(firstName()), strings.Title(lastName()))
}

// Generate a random person first name
func (s stdlib) RandomFirstName() string {
  return strings.Title(firstName())
}

// Generate a random scientist last name
func (s stdlib) RandomLastName() string {
  return strings.Title(lastName())
}

// Generate a random name
func (s stdlib) RandomCompanyName() string {
  return companyName()
}

// Generate a random (v4) UUID
func (s stdlib) RandomUUID() string {
  return uuid.New().String()
}

// Escape a URL query component
func (s stdlib) QueryEscape(v string) string {
  return url.QueryEscape(v)
}

// Unescape a URL component
func (s stdlib) QueryUnescape(v string) (string, error) {
  return url.QueryUnescape(v)
}

// Convert a string to lowercase
func (s stdlib) ToLower(v string) string {
  return strings.ToLower(v)
}

// Convert a string to UPPERCASE
func (s stdlib) ToUpper(v string) string {
  return strings.ToUpper(v)
}

// Convert a string to Title Case
func (s stdlib) ToTitle(v string) string {
  return strings.ToTitle(v)
}

// Trim space from both ends of a string
func (s stdlib) TrimSpace(v string) string {
  return strings.TrimSpace(v)
}
