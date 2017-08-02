package hunit

import (
  "regexp"
  "net/url"
)

// Scheme matcher
var urlish = regexp.MustCompile("^[a-z]+://")

// Determine if a URL appears to be absolute. It is considered absolute
// if the string begins with a scheme://.
func isAbsoluteURL(u string) bool {
  return urlish.MatchString(u)
}

// Merge query parameters with a map
func mergeQueryParams(u string, p map[string]string) (string, error) {
  if len(p) < 1 {
    return u, nil
  }
  
  v, err := url.Parse(u)
  if err != nil {
    return "", err
  }
  
  q := v.Query()
  for k, v := range p {
    q.Add(k, v)
  }
  
  v.RawQuery = q.Encode()
  return v.String(), nil
}
