package text

import (
	"bytes"
  "encoding/json"
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
	var o bytes.Buffer
  if (opt & IndentOptionIndentFirstLine) == IndentOptionIndentFirstLine {
    o.WriteString(p)
  }
  for i := 0; i < len(s); i++ {
		o.WriteString(string(s[i]))
    if s[i] == '\n' {
      o.WriteString(p)
    }
  }
  return o.String()
}

func IsJSON(s string) bool {
    var js map[string]interface{}
    return json.Unmarshal([]byte(s), &js) == nil
}
