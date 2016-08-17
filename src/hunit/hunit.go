package hunit

import (
  "os"
  "io"
  "fmt"
  "time"
  "bytes"
  "strings"
  "strconv"
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
  OptionDebug                         = 1 << 0
  OptionEntityTrimTrailingWhitespace  = 1 << 1
  OptionInterpolateVariables          = 1 << 2
  OptionDisplayRequests               = 1 << 3
  OptionDisplayResponses              = 1 << 4
)

/**
 * A test context
 */
type Context struct {
  BaseURL   string
  Options   Options
  Headers   map[string]string
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
    Headers: c.Headers,
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
  Wait      time.Duration         `yaml:"wait"`
  Request   Request               `yaml:"request"`
  Response  Response              `yaml:"response"`
}

/**
 * Run a test case
 */
func (c Case) Run(context Context) (*Result, error) {
  
  // wait if we need to
  if c.Wait > 0 {
    <- time.After(c.Wait)
  }
  
  // start with an unevaluated result
  result := &Result{Name:fmt.Sprintf("%v %v\n", c.Request.Method, c.Request.URL), Success:true}
  
  method, err := interpolateIfRequired(context, c.Request.Method)
  if err != nil {
    return result.Error(err), nil
  }else if method == "" {
    return nil, fmt.Errorf("Request requires a method (set 'method')")
  }
  
  // incrementally update the name as we evaluate it
  result.Name = fmt.Sprintf("%v %v\n", method, c.Request.URL)
  
  var url string
  if isAbsoluteURL(c.Request.URL) {
    url = c.Request.URL
  }else{
    url = joinPath(context.BaseURL, c.Request.URL)
  }
  
  // incrementally update the name as we evaluate it
  result.Name = fmt.Sprintf("%v %v\n", method, url)
  
  url, err = interpolateIfRequired(context, url)
  if err != nil {
    return result.Error(err), nil
  }else if url == "" {
    return nil, fmt.Errorf("Request requires a URL (set 'url')")
  }
  
  // incrementally update the name as we evaluate it
  result.Name = fmt.Sprintf("%v %v\n", method, url)
  
  var reqdata string
  var entity io.Reader
  if c.Request.Entity != "" {
    reqdata, err = interpolateIfRequired(context, c.Request.Entity)
    if err != nil {
      return result.Error(err), nil
    }else{
      entity = bytes.NewBuffer([]byte(reqdata))
    }
  }
  
  req, err := http.NewRequest(method, url, entity)
  if err != nil {
    return nil, err
  }
  
  if context.Headers != nil {
    for k, v := range context.Headers {
      k, err = interpolateIfRequired(context, k)
      if err != nil {
        return result.Error(err), nil
      }
      v, err = interpolateIfRequired(context, v)
      if err != nil {
        return result.Error(err), nil
      }
      req.Header.Add(k, v)
    }
  }
  
  if c.Request.Headers != nil {
    for k, v := range c.Request.Headers {
      k, err = interpolateIfRequired(context, k)
      if err != nil {
        return result.Error(err), nil
      }
      v, err = interpolateIfRequired(context, v)
      if err != nil {
        return result.Error(err), nil
      }
      req.Header.Add(k, v)
    }
  }
  
  if reqdata != "" {
    req.Header.Add("Content-Length", strconv.FormatInt(int64(len(reqdata)), 10))
  }
  
  if (context.Options & (OptionDisplayRequests | OptionDisplayResponses)) != 0 {
    fmt.Println()
  }
  if (context.Options & OptionDisplayRequests) == OptionDisplayRequests {
    dump := req.Method +" "
    dump += req.URL.Path
    if q := req.URL.RawQuery; q != "" { dump += "?"+ q }
    dump += " "+ req.Proto +"\n"
    
    dump += "Host: "+ req.URL.Host +"\n"
    for k, v := range req.Header {
      dump += k +": "
      for i, e := range v {
        if i > 0 { dump += "," }
        dump += e
      }
      dump += "\n"
    }
    
    dump += "\n"
    if reqdata != "" {
      dump += reqdata +"\n"
    }
    
    fmt.Println(Indent(dump, "> "))
  }
  
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
    return result.Error(err), nil
  }
  
  // check response headers, if necessary
  if headers := c.Response.Headers; headers != nil {
    for k, v := range headers {
      k, err = interpolateIfRequired(context, k)
      if err != nil {
        return result.Error(err), nil
      }
      v, err = interpolateIfRequired(context, v)
      if err != nil {
        return result.Error(err), nil
      }
      result.AssertEqual(v, rsp.Header.Get(k), "Headers do not match: %v", k)
    }
  }
  
  // check response entity, if necessary
  var rspdata []byte
  var rspvalue interface{}
  if entity := c.Response.Entity; entity != "" {
    entity, err = interpolateIfRequired(context, entity)
    if err != nil {
      return result.Error(err), nil
    }
    if rsp.Body == nil {
      result.AssertEqual(entity, "", "Entities do not match")
    }else{
      rspdata, err = ioutil.ReadAll(rsp.Body)
      if err != nil {
        result.Error(err)
      }else if rspvalue, err = entitiesEqual(context, c.Response.Comparison, contentType, []byte(entity), rspdata); err != nil {
        result.Error(err)
      }
    }
  }
  
  if (context.Options & OptionDisplayResponses) == OptionDisplayResponses {
    dump := rsp.Proto +" "+ rsp.Status +"\n"
    
    for k, v := range rsp.Header {
      dump += k +": "
      for i, e := range v {
        if i > 0 { dump += "," }
        dump += e
      }
      dump += "\n"
    }
    
    dump += "\n"
    if rspdata != nil {
      dump += string(rspdata) +"\n"
    }
    
    fmt.Println(Indent(dump, "< "))
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
        "entity": rspdata,
        "value": rspvalue,
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
