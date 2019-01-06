package slug

import (
	"fmt"
	"github.com/bww/go-util/slug"
	"unicode"
)

// Generate a Github-style anchor slug
func Github(t string, c map[string]int) (string, map[string]int) {
	if c == nil {
		c = make(map[string]int)
	}
	var s string
	var ps bool
	for _, e := range t {
		if unicode.IsLetter(e) || unicode.IsDigit(e) {
			s += string(unicode.ToLower(e))
		} else if e == '_' {
			s += string(e)
		} else if unicode.IsSpace(e) && !ps {
			s += "-"
		}
		ps = unicode.IsSpace(e)
	}
	n, ok := c[s]
	c[s] = n + 1
	if ok {
		s += fmt.Sprintf("-%d", n)
	}
	return s, c
}

// Generate a Rails-style anchor slug
func Rails(t string, c map[string]int) (string, map[string]int) {
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
