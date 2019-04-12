package text

import (
	"bytes"
	"io"
)

// Indentiation options
type IndentOptions uint32

const (
	IndentOptionNone            = 0
	IndentOptionIndentFirstLine = 1 << 0
)

// Indent a string by prefixing each line with the provided string
func Indent(s, p string) string {
	return IndentWithOptions(s, p, IndentOptionIndentFirstLine)
}

// Indent a string by prefixing each line with the provided string
func IndentWithOptions(s, p string, opt IndentOptions) string {
	b := &bytes.Buffer{}
	_, err := indent(s, p, opt, b)
	if err != nil {
		panic(err)
	}
	return string(b.Bytes())
}

// Indent and write
type indentWriter struct {
	prefix string
	opt    IndentOptions
	dst    io.Writer
	first  bool
}

// Create an indent writer
func NewIndentWriter(p string, opt IndentOptions, w io.Writer) io.Writer {
	return &indentWriter{p, opt, w, true}
}

// Write
func (w *indentWriter) Write(s []byte) (int, error) {
	opt := w.opt
	if w.first {
		opt |= IndentOptionIndentFirstLine
	} else {
		opt = opt &^ IndentOptionIndentFirstLine
	}
	w.first = false
	return indent(string(s), w.prefix, opt, w.dst)
}

// Indent a string by prefixing each line with the provided string
func indent(s, p string, opt IndentOptions, dst io.Writer) (int, error) {
	var n int
	if (opt & IndentOptionIndentFirstLine) == IndentOptionIndentFirstLine {
		x, err := io.WriteString(dst, p)
		if err != nil {
			return n, err
		}
		n += x
	}
	for i := 0; i < len(s); i++ {
		x, err := io.WriteString(dst, string(s[i]))
		if err != nil {
			return n, err
		}
		n += x
		if s[i] == '\n' {
			x, err = io.WriteString(dst, p)
			if err != nil {
				return n, err
			}
			n += x
		}
	}
	return n, nil
}
