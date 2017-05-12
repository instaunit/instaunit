package hunit

import (
  "io"
  "fmt"
  "time"
  "bytes"
  "strings"
  "strconv"
  "net/http"
  "io/ioutil"
  "hunit/doc"
  "hunit/test"
  "hunit/text"
)

/**
 * HTTP client
 */
var client = http.Client{Timeout: time.Second * 30}

/**
 * A test context
 */
type Context struct {
  BaseURL   string
  Options   test.Options
  Config    test.Config
  Headers   map[string]string
  Debug     bool
  Gendoc    []doc.Generator
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
    Gendoc: c.Gendoc,
    Variables: vars,
  }
}

/**
 * Run a test suite
 */
func RunSuite(s *test.Suite, context Context) ([]*Result, error) {
  results := make([]*Result, 0)
  c := context.Subcontext(make(map[string]interface{}))
  
  for _, e := range context.Gendoc {
    err := e.Prefix(s)
    if err != nil {
      return nil, err
    }
  }
  
  for _, e := range s.Cases {
    r, err := RunTest(e, c)
    if err != nil {
      return nil, err
    }
    results = append(results, r)
  }
  
  for _, e := range context.Gendoc {
    err := e.Suffix(s)
    if err != nil {
      return nil, err
    }
  }
  
  return results, nil
}

/**
 * Run a test case
 */
func RunTest(c test.Case, context Context) (*Result, error) {
  
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
  
  if (context.Options & (test.OptionDisplayRequests | test.OptionDisplayResponses)) != 0 {
    fmt.Println()
  }
  if (context.Options & test.OptionDisplayRequests) == test.OptionDisplayRequests {
    b := &bytes.Buffer{}
    err = text.WriteRequest(b, req, reqdata)
    if err != nil {
      return result.Error(err), nil
    }
    fmt.Println(text.Indent(string(b.Bytes()), "> "))
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
  
  // handle the response entity
  var rspdata []byte
  if rsp.Body != nil {
    rspdata, err = ioutil.ReadAll(rsp.Body)
    if err != nil {
      result.Error(err)
    }
  }
  
  // parse response entity if necessry
  var rspvalue interface{} = rspdata
  if c.Response.Comparison == test.CompareSemantic {
    rspvalue, err = unmarshalEntity(context, contentType, rspdata)
    if err != nil {
      return result.Error(err), nil
    }
  }else if c.Id != "" { // attempt it but don't produce an error if we fail
    val, err := unmarshalEntity(context, contentType, rspdata)
    if err == nil {
      rspvalue = val
    }
  }
  
  // check response entity, if necessary
  if entity := c.Response.Entity; entity != "" {
    entity, err = interpolateIfRequired(context, entity)
    if err != nil {
      return result.Error(err), nil
    }
    if len(rspdata) == 0 {
      result.AssertEqual(entity, "", "Entities do not match")
    }else if err = entitiesEqual(context, c.Response.Comparison, contentType, []byte(entity), rspvalue); err != nil {
      result.Error(err)
    }
  }
  
  if (context.Options & test.OptionDisplayResponses) == test.OptionDisplayResponses {
    b := &bytes.Buffer{}
    err = text.WriteResponse(b, rsp, rspdata)
    if err != nil {
      return result.Error(err), nil
    }
    fmt.Println(text.Indent(string(b.Bytes()), "< "))
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
  
  // generate documentation if necessary
  if c.Documented() && len(context.Gendoc) > 0 {
    for _, e := range context.Gendoc {
      err := e.Generate(c, req, reqdata, rsp, rspdata)
      if err != nil {
        return nil, fmt.Errorf("Could not generate documentation: %v", err)
      }
    }
  }
  
  return result, nil
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
