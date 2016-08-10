package hunit

import (
  "os"
  "io"
  "fmt"
  "time"
  "bytes"
  "strings"
  "net/http"
  "io/ioutil"
)

import (
  "gopkg.in/yaml.v2"
)

const whitespace = " \n\r\t\v"

/**
 * HTTP client
 */
var client = http.Client{Timeout: time.Second * 30}

/**
 * Options
 */
type Options uint32

const (
  OptionNone                          = 0
  OptionEntityTrimTrailingWhitespace  = 1 << 0
  OptionInterpolateVariables          = 1 << 1
  OptionInterpolateEnvironment        = 1 << 2
)

/**
 * A test context
 */
type Context struct {
  BaseURL   string
  Options   Options
  Debug     bool
  Variables map[string]interface{}
}

/**
 * Derive a subcontext
 */
func (c Context) Subcontext(vars map[string]interface{}) Context {
  return Context{
    BaseURL: c.BaseURL,
    Options: c.Options,
    Debug: c.Debug,
    Variables: vars,
  }
}

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
  Request   Request               `yaml:"request"`
  Response  Response              `yaml:"response"`
}

/**
 * Run a test case
 */
func (c Case) Run(context Context) (*Result, error) {
  
  method, err := interpolateIfRequired(context, c.Request.Method)
  if err != nil {
    return nil, err
  }else if method == "" {
    return nil, fmt.Errorf("Request requires a method (set 'method')")
  }
  
  var url string
  if isAbsoluteURL(c.Request.URL) {
    url = c.Request.URL
  }else{
    url = joinPath(context.BaseURL, c.Request.URL)
  }
  
  url, err = interpolateIfRequired(context, url)
  if err != nil {
    return nil, err
  }else if url == "" {
    return nil, fmt.Errorf("Request requires a URL (set 'url')")
  }
  
  var entity io.Reader
  if c.Request.Entity != "" {
    expand, err := interpolateIfRequired(context, c.Request.Entity)
    if err != nil {
      return nil, err
    }else{
      entity = bytes.NewBuffer([]byte(expand))
    }
  }
  
  req, err := http.NewRequest(method, url, entity)
  if err != nil {
    return nil, err
  }
  
  if c.Request.Headers != nil {
    for k, v := range c.Request.Headers {
      k, err = interpolateIfRequired(context, k)
      if err != nil {
        return nil, err
      }
      v, err = interpolateIfRequired(context, v)
      if err != nil {
        return nil, err
      }
      req.Header.Add(k, v)
    }
  }
  
  result := &Result{Name:fmt.Sprintf("%v %v\n", method, url), Success:true}
  
  rsp, err := client.Do(req)
  if rsp != nil && rsp.Body != nil {
    defer rsp.Body.Close()
  }
  if err != nil {
    return result.Error(err), nil
  }
  
  // check the response status
  result.AssertEqual(c.Response.Status, rsp.StatusCode, "Unexpected status code")
  
  // note the content type; we prefer the explicit format over the content type
  var contentType string
  if c.Response.Format != "" {
    contentType = c.Response.Format
  }else{
    contentType = strings.ToLower(rsp.Header.Get("Content-Type"))
  }
  
  contentType, err = interpolateIfRequired(context, contentType)
  if err != nil {
    return nil, err
  }
  
  // check response headers, if necessary
  if headers := c.Response.Headers; headers != nil {
    for k, v := range headers {
      k, err = interpolateIfRequired(context, k)
      if err != nil {
        return nil, err
      }
      v, err = interpolateIfRequired(context, v)
      if err != nil {
        return nil, err
      }
      result.AssertEqual(v, rsp.Header.Get(k), "Headers do not match: %v", k)
    }
  }
  
  // check response entity, if necessary
  var data []byte
  if entity := c.Response.Entity; entity != "" {
    entity, err = interpolateIfRequired(context, entity)
    if err != nil {
      return nil, err
    }
    if rsp.Body == nil {
      result.AssertEqual(entity, "", "Entities do not match")
    }else{
      data, err = ioutil.ReadAll(rsp.Body)
      if err != nil {
        result.Error(err)
      }else if err = entitiesEqual(context, c.Response.Comparison, contentType, []byte(entity), data); err != nil {
        result.Error(err)
      }
    }
  }
  
  // add to our context if we have an identifier
  if c.Id != "" {
    
    headers := make(map[string]string)
    for k, v := range rsp.Header {
      if len(v) > 0 {
        headers[k] = v[0]
      }
    }
    
    context.Variables[c.Id] = map[string]interface{}{
      "case": c,
      "response": map[string]interface{}{
        "headers": headers,
        "entity": data,
        "status": rsp.StatusCode,
      },
    }
    
  }
  
  return result, nil
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

/**
 * Run a test suite
 */
func (s *Suite) Run(context Context) ([]*Result, error) {
  results := make([]*Result, 0)
  c := context.Subcontext(make(map[string]interface{}))
  
  for _, e := range s.Cases {
    r, err := e.Run(c)
    if err != nil {
      return nil, err
    }
    results = append(results, r)
  }
  
  return results, nil
}

/**
 * Join paths
 */
func joinPath(a, b string) string {
  if a == "" {
    return b
  }
  
  var i, j int
  for i = len(a) - 1; i >= 0; i-- {
    if a[i] != '/' {
      break
    }
  }
  
  for j = 0; j < len(b); j++ {
    if b[j] != '/' {
      break
    }
  }
  
  return a[:i + 1] +"/"+ b[j:]
}
