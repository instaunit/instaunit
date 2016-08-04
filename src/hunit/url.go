package hunit

import (
  "regexp"
)

/**
 * Scheme matcher
 */
var urlish = regexp.MustCompile("^[a-z]+://")

/**
 * Determine if a URL appears to be absolute. It is considered so
 * if the string begins with a scheme.
 */
func isAbsoluteURL(u string) bool {
  return urlish.MatchString(u)
}
