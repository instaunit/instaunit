import test

import (
  "os"
  "io/ioutil"
)

import (
  "gopkg.in/yaml.v2"
)

/**
 * Suite options
 */
type SuiteOptions struct {
  Doc struct {
    IncludeRequestHTTP  bool  `yaml:"doc-include-request-http"`
    IncludeResponseHTTP bool  `yaml:"doc-include-response-http"`
    IncludeHTTP         bool  `yaml:"doc-include-http"`
  }                           `yaml:",inline"`
}

/**
 * A test suite
 */
type Suite struct {
  Title       string          `yaml:"title"`
  Comments    string          `yaml:"doc"`
  Cases       []Case          `yaml:"tests"`
  Options     SuiteOptions    `yaml:"options"`
}

/**
 * Load a test suite
 */
func LoadSuiteFromFile(p string) (*Suite, error) {
  
  file, err := os.Open(p)
  if err != nil {
    return nil, err
  }
  
  data, err := ioutil.ReadAll(file)
  if err != nil {
    return nil, err
  }
  
  return LoadSuite(data)
}

/**
 * Load a test suite
 */
func LoadSuite(data []byte) (*Suite, error) {
  var ferr error
  
  suite := &Suite{}
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
  
  fmt.Println("OPTIONS:", suite.Options)
  
  return suite, nil
}

/**
 * Return the first non-nil error or nil if there are none.
 */
func coalesce(err ...error) error {
  for _, e := range err {
    if e != nil {
      return e
    }
  }
  return nil
}
