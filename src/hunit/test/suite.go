package test

import (
  "os"
  "time"
  "io/ioutil"
  
  "hunit/exec"
)

import (
  "gopkg.in/yaml.v2"
)

// Suite options
type Config struct {
  Net struct {
    StreamIOGracePeriod   time.Duration       `yaml:"stream-io-grace-period"`
  }                                           `yaml:",inline"`
  Doc struct {
    AnchorStyle           AnchorStyle         `yaml:"anchor-style"`
    FormatEntities        bool                `yaml:"format-entities"`
    IncludeRequestHTTP    bool                `yaml:"doc-include-request-http"`
    IncludeResponseHTTP   bool                `yaml:"doc-include-response-http"`
    IncludeHTTP           bool                `yaml:"doc-include-http"`
  }                                           `yaml:",inline"`
}

// A test suite
type Suite struct {
  Title         string            `yaml:"title"`
  Comments      string            `yaml:"doc"`
  Cases         []Case            `yaml:"tests"`
  Config        Config            `yaml:"options"`
  Setup         []exec.Command    `yaml:"setup"`
  Teardown      []exec.Command    `yaml:"teardown"`
}

// Load a test suite
func LoadSuiteFromFile(c *Config, p string) (*Suite, error) {
  
  file, err := os.Open(p)
  if err != nil {
    return nil, err
  }
  
  data, err := ioutil.ReadAll(file)
  if err != nil {
    return nil, err
  }
  
  return LoadSuite(c, data)
}

// Load a test suite
func LoadSuite(conf *Config, data []byte) (*Suite, error) {
  var ferr error
  
  suite := &Suite{Config:*conf}
  err := yaml.Unmarshal(data, suite)
  if err != nil {
    ferr = err
  }
  
  if len(suite.Cases) < 1 {
    var cases []Case
    err := yaml.Unmarshal(data, &cases)
    if err != nil {
      return nil, coalesce(ferr, err)
    }else{
      suite.Cases = cases
    }
  }
  
  *conf = suite.Config
  return suite, nil
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
