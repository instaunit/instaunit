package hunit

import (
	"github.com/instaunit/instaunit/hunit/script"
	"github.com/instaunit/instaunit/hunit/text"

	"github.com/davecgh/go-spew/spew"
)

// A script error
type ScriptError struct {
	Message          string
	Expected, Actual interface{}
	Script           *script.Script
}

func (e ScriptError) Error() string {
	m := e.Message
	if m != "" {
		m += ":\n"
	}
	m += "expected: " + spew.Sdump(e.Expected)
	m += "  actual: " + spew.Sdump(e.Actual)
	m += "--\n" + text.IndentWithOptions(e.Script.Source, `{{ printf "%03d: " .Line }}`, text.IndentOptionIndentFirstLine|text.IndentOptionIndentTemplate)
	return m
}
