package test

import (
  "os"
  "time"
  "io/ioutil"
)

import (
  "gopkg.in/yaml.v2"
)

/**
 * Options
 */
type Options uint32

const (
  OptionNone                          = 0
  OptionDebug                         = 1 << 0
  OptionEntityTrimTrailingWhitespace  = 1 << 1
  OptionInterpolateVariables          = 1 << 2
  OptionDisplayRequests               = 1 << 3
  OptionDisplayResponses              = 1 << 4
)

/**
 * A test request
 */
type Request struct {
  Method    string                `yaml:"method"`
  URL       string                `yaml:"url"`
  Headers   map[string]string     `yaml:"headers"`
  Entity    string                `yaml:"entity"`
}

/**
 * A test response
 */
type Response struct {
  Status      int                 `yaml:"status"`
  Headers     map[string]string   `yaml:"headers"`
  Entity      string              `yaml:"entity"`
  Comparison  Comparison          `yaml:"compare"`
  Format      string              `yaml:"format"`
}

/**
 * A test case
 */
type Case struct {
  Id        string                `yaml:"id"`
  Wait      time.Duration         `yaml:"wait"`
  Gendoc    bool                  `yaml:"gendoc"`
  Title     string                `yaml:"title"`
  Comments  string                `yaml:"doc"`
  Request   Request               `yaml:"request"`
  Response  Response              `yaml:"response"`
}

/**
 * Determine if this case is documented or not
 */
func (c Case) Documented() bool {
  return c.Gendoc || c.Title != "" || c.Comments != ""
}

/**
 * A test suite
 */
type Suite struct {
  Cases []Case
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
  
  var cases []Case
  err := yaml.Unmarshal(data, &cases)
  if err != nil {
    return nil, err
  }
  
  return &Suite{cases}, nil
}
