package rest

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type alternateErrors []error

func (errs alternateErrors) Error() string {
	switch len(errs) {
	case 0:
		return "No error"
	case 1:
		return errs[0].Error()
	}

	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("One of %d possible errors occurred:\n", len(errs)))

	for i, err := range errs {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("#%d: %v\n", i+1, err))
	}

	return b.String()
}

// A request
type Request struct {
	sync.Mutex
	Methods []string          `yaml:"methods"`
	Path    string            `yaml:"path"`
	Params  map[string]string `yaml:"params"`
	Headers map[string]string `yaml:"headers"`
	Cookies map[string]string `yaml:"cookies"`
	Entity  string            `yaml:"entity"`
}

// A response
type Response struct {
	Status  int               `yaml:"status"`
	Headers map[string]string `yaml:"headers"`
	Cookies map[string]string `yaml:"cookies"`
	Entity  string            `yaml:"entity"`
}

// An endpoint
type Endpoint struct {
	Wait     time.Duration `yaml:"wait"`
	Request  *Request      `yaml:"endpoint"`
	Response *Response     `yaml:"response"`
}

// A test suite
type Suite struct {
	Endpoints []Endpoint `yaml:"service"`
}

// Load a test suite
func LoadSuite(src io.ReadCloser) (*Suite, error) {
	suite := &Suite{}
	var errs []error

	data, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}

	err = unmarshal(data, suite)
	if err != nil {
		errs = append(errs, err)
	}

	if len(suite.Endpoints) < 1 {
		var endpoints []Endpoint
		err := unmarshal(data, &endpoints)
		if err != nil {
			errs = append(errs, err)
		} else {
			suite.Endpoints = endpoints
		}
	}

	if len(suite.Endpoints) < 1 && len(errs) > 0 {
		return nil, alternateErrors(errs)
	} else {
		return suite, nil
	}
}

func unmarshal(data []byte, dest interface{}) error {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	return dec.Decode(dest)
}

// Return the first non-nil error or nil if there are none.
func coalesce(err ...error) error {
	for _, e := range err {
		if e != nil {
			return e
		}
	}
	return nil
}
