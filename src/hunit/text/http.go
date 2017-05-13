package text

import (
  "io"
  "fmt"
  "path"
  "strings"
  "net/http"
)

/**
 * HTTP writing config
 */
type HttpConfig struct {
  RewriteHeaders          map[string]string   `yaml:"rewrite-headers"`
  SuppressHeaders         []string            `yaml:"suppress-headers"`
  AllowHeaders            []string            `yaml:"allow-headers"`
}

/**
 * Initialize patterns, etc
 */
/*
func (c *HttpConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
  err := unmarshal(c)
  if err != nil {
    return err
  }
  return nil
}
*/

/**
 * Obtain a header value, accounting for rewrite rules
 */
func (c HttpConfig) Header(n string, v []string) []string {
  
  // first, suppress specifically excluded headers
  for _, e := range c.SuppressHeaders {
    if strings.EqualFold(e, n) {
      return nil
    }
  }
  
  // then, if allowed headers are defined, all others are suppressed
  // an empty list is allowed and excludes all headers
  if c.AllowHeaders != nil {
    allow := false
    for _, e := range c.AllowHeaders {
      if strings.EqualFold(e, n) {
        allow = true; break
      }
    }
    if !allow {
      return nil
    }
  }
  
  // finally, rewrite headers
  if c.RewriteHeaders != nil {
    for k, r := range c.RewriteHeaders {
      if strings.EqualFold(k, n) {
        return []string{r}
      }
    }
  }
  
  // if none of that applies, just return the default value
  return v
}

/**
 * Write a request to the specified output
 */
func WriteRequest(w io.Writer, conf HttpConfig, req *http.Request, entity string) error {
  var dump string
  
  if req != nil {
    dump += req.Method +" "
    dump += req.URL.Path
    if q := req.URL.RawQuery; q != "" { dump += "?"+ q }
    dump += " "+ req.Proto +"\n"
    
    req.Header.Set("Host", req.URL.Host)
    for k, v := range req.Header {
      if v = conf.Header(k, v); v != nil {
        dump += k +": "
        for i, e := range v {
          if i > 0 { dump += "," }
          dump += e
        }
        dump += "\n"
      }
    }
  }
  
  if entity != "" {
    dump += "\n"
    dump += entity
  }
  
  _, err := w.Write([]byte(dump))
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * Write a response to the specified output
 */
func WriteResponse(w io.Writer, conf HttpConfig, rsp *http.Response, entity []byte) error {
  var dump string
  
  if rsp != nil {
    dump += rsp.Proto +" "+ rsp.Status +"\n"
    
    for k, v := range rsp.Header {
      if v = conf.Header(k, v); v != nil {
        dump += k +": "
        for i, e := range v {
          if i > 0 { dump += "," }
          dump += e
        }
        dump += "\n"
      }
    }
  }
  
  if entity != nil {
    dump += "\n"
    dump += string(entity)
  }
  
  _, err := w.Write([]byte(dump))
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * Determine if the provided request has a particular content type
 */
func HasContentType(req *http.Request, t string) bool {
  return MatchesContentType(t, req.Header.Get("Content-Type"))
}

/**
 * Determine if the provided request has a particular content type
 */
func MatchesContentType(pattern, contentType string) bool {
  
  // trim off the parameters following ';' if we have any
  if i := strings.Index(contentType, ";"); i > 0 {
    contentType = contentType[:i]
  }
  
  // path.Match does glob matching, which is useful it we
  // want to, e.g., test for all image types with `image/*`.
  m, err := path.Match(pattern, contentType)
  if err != nil {
    panic(fmt.Errorf("* * * could not match invalid content-type pattern:", pattern))
  }
  
  return m
}
