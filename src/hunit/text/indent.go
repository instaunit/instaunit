package text

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
  var o string
  if (opt & IndentOptionIndentFirstLine) == IndentOptionIndentFirstLine {
    o += p
  }
  for i := 0; i < len(s); i++ {
    o += string(s[i])
    if s[i] == '\n' {
      o += p
    }
  }
  return o
}
