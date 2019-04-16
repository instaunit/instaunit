package report

import (
	"fmt"
	"io"
	"net/http"

	"github.com/instaunit/instaunit/hunit"
	"github.com/instaunit/instaunit/hunit/report/emit"
	"github.com/instaunit/instaunit/hunit/report/emit/junit"
	"github.com/instaunit/instaunit/hunit/test"
)

// Suite results
type Results struct {
	Results []*hunit.Result
	Runtime time.Duration
}

// Implemented by report generators
type Generator interface {
	Init() error
	Suite(test.Config, *test.Suite, *Results) error
	Finalize() error
}

// Create a report emitter
func New(t emit.Doctype, o io.Writer) (Generator, error) {
	switch t {
	case emit.DoctypeJUnitXML:
		return junit.New(o), nil
	default:
		return nil, fmt.Errorf("Unsupported report type: %v", t)
	}
}
