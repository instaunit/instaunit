package junit

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/instaunit/instaunit/hunit/report/emit"
	"github.com/instaunit/instaunit/hunit/test"
)

type testfail struct {
	Message string `xml:"message,attr,omitempty"`
	Type    string `xml:"type,attr"`
	Detail  string `xml:",cdata"`
}

type testcase struct {
	Id       string     `xml:"id,attr,omitempty"`
	Name     string     `xml:"name,attr,omitempty"`
	Duration float64    `xml:"time,attr"`
	Failures []testfail `xml:"failure,omitempty"`
}

type testsuite struct {
	Id       string     `xml:"id,attr,omitempty"`
	Name     string     `xml:"name,attr,omitempty"`
	Duration float64    `xml:"time,attr"`
	Cases    []testcase `xml:"testcase"`
}

type testsuites struct {
	Id       string      `xml:"id,attr,omitempty"`
	Name     string      `xml:"name,attr,omitempty"`
	Tests    int         `xml:"tests,attr"`
	Failures int         `xml:"failures,attr"`
	Duration float64     `xml:"time,attr"`
	Suites   []testsuite `xml:"testsuite,omitempty"`
}

// A junit report generator
type Generator struct {
	w               io.WriteCloser
	id              string
	tests, failures int
	duration        time.Duration
	suites          []testsuite
}

// Produce a new emitter
func New(w io.WriteCloser, id string) *Generator {
	return &Generator{w: w, id: id}
}

// Initialize the report
func (g *Generator) Init() error {
	g.suites = nil
	return nil
}

// Finalize the report
func (g *Generator) Finalize() error {
	ts := testsuites{
		Id:       g.id,
		Name:     "Name",
		Tests:    g.tests,
		Failures: g.failures,
		Duration: float64(g.duration) / float64(time.Second),
		Suites:   g.suites,
	}

	enc := xml.NewEncoder(g.w)
	enc.Indent("", "  ")
	err := enc.Encode(ts)
	if err != nil {
		return err
	}

	return g.w.Close()
}

// Generate a report for the provided suite
func (g *Generator) Suite(conf test.Config, suite *test.Suite, results *emit.Results) error {
	var success, failure int
	for _, e := range results.Results {
		if e.Success {
			success++
		} else {
			failure++
		}
	}

	sid := len(g.suites) + 1
	tc := make([]testcase, len(results.Results))
	for i, e := range results.Results {
		var tf []testfail
		if len(e.Errors) > 0 {
			for _, err := range e.Errors {
				tf = append(tf, testfail{
					Type:   "ERROR",
					Detail: err.Error(),
				})
			}
		} else if !e.Success {
			tf = append(tf, testfail{
				Type:    "ERROR",
				Message: "The test failed. That's all we know.",
			})
		}
		tc[i] = testcase{
			Id:       fmt.Sprintf("%s_%d_%d", g.id, sid, i+1),
			Name:     e.Name,
			Duration: float64(e.Runtime) / float64(time.Second),
			Failures: tf,
		}
	}

	ts := testsuite{
		Id:       fmt.Sprintf("%s_%d", g.id, sid),
		Name:     suite.Title,
		Duration: float64(results.Runtime) / float64(time.Second),
		Cases:    tc,
	}

	g.suites = append(g.suites, ts)
	g.tests += len(results.Results)
	g.failures += failure
	g.duration += results.Runtime
	return nil
}
