package hunit

import (
	"time"

	"github.com/instaunit/instaunit/hunit/assert"
	"github.com/instaunit/instaunit/hunit/runtime"
	"github.com/instaunit/instaunit/hunit/testcase"
)

// A test result
type Result struct {
	Name          string          `json:"name"`
	Success       bool            `json:"success"`
	Skipped       bool            `json:"skipped"`
	Errors        []string        `json:"errors,omitempty"`
	FormatReqdata string          `json:"formatted_request_data,omitempty"`
	Reqdata       []byte          `json:"request_data,omitempty"`
	FormatRspdata string          `json:"formatted_response_data,omitempty"`
	Rspdata       []byte          `json:"response_data,omitempty"`
	Context       runtime.Context `json:"context"`
	Runtime       time.Duration   `json:"duration"`
	Case          testcase.Case   `json:"case"`
}

// Assert equality. If the values are not equal an error is added to the result.
func (r *Result) AssertEqual(e, a interface{}, m string, x ...interface{}) bool {
	err := assert.Equal(e, a, m, x...)
	if err != nil {
		r.Error(err)
		return false
	}
	return true
}

// Add an error to the result. The result is marked as unsuccessful and
// the result is returned so calls can be chained.
func (r *Result) Error(e error) *Result {
	r.Success = false
	r.Errors = append(r.Errors, e.Error())
	return r
}
