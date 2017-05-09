package text

import (
	// "bytes"
	"strings"
)

/**
 * Indentiation options
 */
type IndentOptions uint32
const (
  IndentOptionNone            = 0
  IndentOptionIndentFirstLine = 1 << 0
)

/**
 * Indent a string by prefixing each line with the provided string
 */
func Indent(s, p string) string {
  return IndentWithOptions(s, p, IndentOptionIndentFirstLine)
}

/**
 * Indent a string by prefixing each line with the provided string
 */
func IndentWithOptions(s, p string, opt IndentOptions) string {
  // var o string
	// var o bytes.Buffer
	o := []string{}
  if (opt & IndentOptionIndentFirstLine) == IndentOptionIndentFirstLine {
    // o.WriteString(p)
		o = append(o, p)
  }
  for i := 0; i < len(s); i++ {
    // o += string(s[i])
		// o.WriteString(string(s[i]))
		o = append(o, string(s[i]))
    if s[i] == '\n' {
      // o.WriteString(p)
			o = append(o, p)
    }
  }
  return strings.Join(o, ``)
}
