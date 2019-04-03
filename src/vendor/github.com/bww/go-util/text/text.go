package text

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Normalize a string for general purpose matching
func NormalizeString(s string) string {
	return normalize(s, "-_',", ' ')
}

// Normalize query terms
func NormalizeTerms(s string) string {
	return normalize(s, "-_'", ' ')
}

// Identize terms
func IdentizeString(s string) string {
	return normalize(s, "_", '_')
}

func normalize(s, special string, space rune) string {
	n := &strings.Builder{}

	var insp bool
	for _, e := range s {
		if unicode.IsSpace(e) {
			if n.Len() > 0 {
				insp = true
			}
		} else {
			if unicode.IsLetter(e) {
				if insp {
					_, err := n.WriteRune(space)
					if err != nil {
						panic(err)
					}
				}
				_, err := n.WriteRune(unicode.ToLower(e))
				if err != nil {
					panic(err)
				}
				insp = false
			} else if unicode.IsDigit(e) || allowed(e, special) {
				if insp {
					_, err := n.WriteRune(space)
					if err != nil {
						panic(err)
					}
				}
				_, err := n.WriteRune(e)
				if err != nil {
					panic(err)
				}
				insp = false
			} else if n.Len() > 0 {
				insp = true
			}
		}
	}

	return n.String()
}

func allowed(e rune, allow string) bool {
	for _, x := range allow {
		if e == x {
			return true
		}
	}
	return false
}

/**
 * Collapse whitespace
 */
func CollapseSpaces(s string) string {
	n := &strings.Builder{}

	var insp bool
	for _, e := range s {
		if unicode.IsSpace(e) {
			if n.Len() > 0 {
				insp = true
			}
		} else {
			if insp {
				_, err := n.WriteRune(' ')
				if err != nil {
					panic(err)
				}
			}
			_, err := n.WriteRune(e)
			if err != nil {
				panic(err)
			}
			insp = false
		}
	}

	return n.String()
}

// Truncate a string to n characters (not bytes). If the string is truncated,
// the provided suffix is appended. Something like ' [...]' would be appropriate
// as a suffix to indicate that text was removed.
func Truncate(s string, n int, x string) string {
	d := s
	l, b, c := len(s), 0, 0
	for i := 0; i < l && i < n; i++ {
		_, w := utf8.DecodeRuneInString(d)
		c += 1
		b += w
		d = d[w:]
	}
	s = s[:b]
	if b < l {
		s += x
	}
	return s
}

// Normalize a list, using a special final delimiter between the last
// two elements.
func NormalizeJoin(l []string, d, f string) string {
	n := len(l)
	var s string
	for i, e := range l {
		if i > 0 {
			if n-(i+1) == 0 {
				s += f
			} else {
				s += d
			}
		}
		s += e
	}
	return s
}

// Return the first non-empty string from those provided
func Coalesce(v ...string) string {
	for _, e := range v {
		if e != "" {
			return e
		}
	}
	return ""
}

// Strip out control characters in place
func StripControl(s string) string {
	c := make([]byte, len([]byte(s)))
	n := 0
	for _, e := range s {
		if !unicode.IsControl(e) {
			b := []byte(string(e))
			copy(c[n:], b)
			n += len(b)
		}
	}
	return string(c[:n])
}

// Non-space marks
func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

// Convert characters with diacritical marks to their unaccented/base
// counterparts. See also:
// http://stackoverflow.com/questions/26722450/remove-diacritics-using-go
func NormalizeDiacritics(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}
