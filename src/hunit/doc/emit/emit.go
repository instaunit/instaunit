package emit

import (
	"fmt"
)

// Documentation format type
type Doctype uint32

const (
	DoctypeMarkdown Doctype = iota
	DoctypeInstadoc
	DoctypeOpenAPI
	DoctypeInvalid
)

var doctypeNames = []string{
	"markdown",
	"instadoc",
	"openapi",
	"<invalid>",
}

var doctypeExts = []string{
	".md",
	".instadoc",
	".json",
	".???",
}

// Parse a doctype
func ParseDoctype(s string) (Doctype, error) {
	switch s {
	case "markdown":
		return DoctypeMarkdown, nil
	case "instadoc":
		return DoctypeInstadoc, nil
	case "openapi":
		return DoctypeOpenAPI, nil
	default:
		return DoctypeInvalid, fmt.Errorf("Unsupported type: %v", s)
	}
}

// Extension
func (c Doctype) Ext() string {
	if c < 0 || c >= DoctypeInvalid {
		return ""
	} else {
		return doctypeExts[int(c)]
	}
}

// Stringer
func (c Doctype) String() string {
	if c < 0 || c >= DoctypeInvalid {
		return "<invalid>"
	} else {
		return doctypeNames[int(c)]
	}
}

// Marshal
func (d Doctype) MarshalYAML() (interface{}, error) {
	return d.String(), nil
}

// Unmarshal
func (d *Doctype) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	*d, err = ParseDoctype(s)
	if err != nil {
		return err
	}
	return nil
}
