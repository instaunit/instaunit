package rest

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/instaunit/instaunit/hunit/service/backend/errors"

	yaml "gopkg.in/yaml.v3"
)

// A request
type Request struct {
	sync.Mutex
	Method  string            `yaml:"method"`
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
	Globals   map[string]any `yaml:"vars"`
	Endpoints []Endpoint     `yaml:"service"`
}

// Load a test suite
func LoadSuite(src io.ReadCloser) (*Suite, error) {
	suite := &Suite{}
	var errs []error

	data, err := io.ReadAll(src)
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
		return nil, fmt.Errorf("Could not unmarshal REST service: %w", errors.AlternateErrors(errs))
	}

	return suite, nil
}

func unmarshal(data []byte, dest interface{}) error {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	return dec.Decode(dest)
}
