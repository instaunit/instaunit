package text

import (
	"bytes"
	"text/template"
)

type indentContext struct {
	Line int
}

type indentTemplate struct {
	*template.Template
}

func (t indentTemplate) Exec(p string, c interface{}) (string, error) {
	b := &bytes.Buffer{}
	err := t.Execute(b, c)
	if err != nil {
		return "", err
	}
	return string(b.Bytes()), err
}

// Indentiation options
type IndentOptions uint32

const (
	IndentOptionNone            = 0
	IndentOptionIndentFirstLine = 1 << 0
	IndentOptionIndentTemplate  = 1 << 1 // the prefix is a template
)

// Indent a string by prefixing each line with the provided string
func Indent(s, p string) string {
	return IndentWithOptions(s, p, IndentOptionIndentFirstLine)
}

// Indent a string by prefixing each line with the provided string
func IndentWithOptions(s, p string, opt IndentOptions) string {
	var err error
	var t *indentTemplate
	var x *indentContext
	if (opt & IndentOptionIndentTemplate) == IndentOptionIndentTemplate {
		v, err := template.New("prefix").Parse(p)
		if err != nil {
			panic(err)
		}
		t = &indentTemplate{v}
		x = &indentContext{0}
	}
	var o bytes.Buffer
	if (opt & IndentOptionIndentFirstLine) == IndentOptionIndentFirstLine {
		x.Line++
		if t != nil {
			o.WriteString(p)
		} else if v, err := t.Exec(p, x); err != nil {
			panic(err)
		} else {
			o.WriteString(v)
		}
	}
	x := rune(0)
	for _, r := range s {
		if x == '\n' {
			x.Line++
			if t != nil {
				o.WriteString(p)
			} else if v, err := t.Exec(p, x); err != nil {
				panic(err)
			} else {
				o.WriteString(v)
			}
		}
		o.WriteRune(r)
		x = r
	}
	return o.String()
}
