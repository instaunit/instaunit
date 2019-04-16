package junit

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/instaunit/instaunit/hunit"
	"github.com/instaunit/instaunit/hunit/test"
)

// A junit report generator
type Generator struct {
	w io.Writer
	b *bytes.Buffer
}

// Produce a new emitter
func New(w io.Writer) *Generator {
	return &Generator{w, nil}
}

// Init a suite
func (g *Generator) Init() error {
	g.b = &bytes.Buffer{}
	return nil
}

// Finish a suite
func (g *Generator) Finalize() error {
	return nil
}

// Generate a report for the provided suite
func (g *Generator) Sute(conf test.Config, suite *test.Suite, results *Results) error {
	var success, failure int
	for _, e := range results.Results {
		if e.Success {
			success++
		} else {
			failure++
		}
	}

	fmt.Fprintf(g.b, `<?xml version="1.0" encoding="UTF-8" ?>
  <testsuites id="20140612_170519" name="New_configuration (14/06/12 17:05:19)" tests="%d" failures="%d" time="%f">`,
		suite.Title, len(results), failure, results.Runtime)
	return nil
}
