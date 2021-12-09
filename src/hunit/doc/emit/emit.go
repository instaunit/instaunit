package emit

import (
	"fmt"
	"net/http"

	"github.com/instaunit/instaunit/hunit/route"
	"github.com/instaunit/instaunit/hunit/test"
)

// A test case to be documented
type Case struct {
	Case     test.Case
	Route    *route.Route
	Request  *http.Request
	Reqdata  []byte
	Response *http.Response
	Rspdata  []byte
}

// Documentation format type
type Doctype uint32

const (
	DoctypeMarkdown Doctype = iota
	DoctypeConfluence
	DoctypeInvalid
)

var doctypeNames = []string{
	"markdown",
	"<invalid>",
}

var doctypeExts = []string{
	".md",
	".???",
}

// Parse a doctype
func ParseDoctype(s string) (Doctype, error) {
	switch s {
	case "markdown":
		return DoctypeMarkdown, nil
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
