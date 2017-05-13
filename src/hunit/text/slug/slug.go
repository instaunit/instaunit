package slug

import (
  "fmt"
  "github.com/bww/go-util/slug"
)

/**
 * Generate a Github-style anchor slug
 */
func Github(t string, c map[string]int) (string, map[string]int) {
  if c == nil {
    c = make(map[string]int)
  }
  s := slug.Slugify(t)
  n, ok := c[s]
  c[s] = n + 1
  if ok {
    s += fmt.Sprintf("-%d", n)
  }
  return s, c
}
