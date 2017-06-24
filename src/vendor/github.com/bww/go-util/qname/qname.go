package qname

import (
  "fmt"
  "strings"
  "unicode/utf8"
)

const wildcard = "*"

/**
 * Matching options
 */
const (
  MATCH_OPTION_NONE       = 0
  MATCH_OPTION_ENCOMPASS  = 1 << 0
)

/**
 * A qualified name
 */
type QName string

/**
 * Format a QName
 */
func QNamef(f string, a ...interface{}) QName {
  return QName(fmt.Sprintf(f, a...))
}

/**
 * This is the equivalent of MatchesWithOptions(z, MATCH_OPTION_NONE)
 */
func (q QName) Matches(z QName) bool {
  return q.MatchesWithOptions(z, MATCH_OPTION_NONE)
}

/**
 * Determine if this name matches the provided name. A name matches another if all the
 * components of both names are identical (ignoring case), or if all the concrete names
 * in the right (parameter) name match those in the left (this) name accounting for
 * wildcards, or if the left name is the name "*", which matches any other name.
 * 
 * If encompassing is permitted, a left name which ends in a wildcard component will match
 * any number of subsequent components in the right name (e.g., "a.*" matches "a.b", "a.b.c",
 * "a.b.c.d", and so fourth).
 * 
 * For example, the concrete name "a.b" matches: "a.b"
 * 
 * The name "a.*.c" matches: "a.b.c", "a.z.c", "a.*.c" ("*" is interpreted literally in the
 * right name). However, the name "a.*.c" does not match: "*.b.c" or "a.c"
 * 
 * When encompassing, the name "a.b.*" matches "a.b.c", "a.b.c.d", and so fourth
 * However it does not match: "a" or "a.b"
 * 
 * The name "*" matches any name, including an emtpy name or the name "*". An empty name
 * matches nothing.
 */
func (q QName) MatchesWithOptions(z QName, options int) bool {
  a := string(q)
  b := string(z)
  
  enc := (options & MATCH_OPTION_ENCOMPASS) == MATCH_OPTION_ENCOMPASS
  
  if a == "" {
    return false  // nothing can't match anything
  }else if a == "*" {
    return true   // '*' matches anything, including nothing
  }else if b == "" {
    return false  // nothing but '*' can match nothing (!)
  }
  
  var wc bool
  for {
    var ca, cb string
    var aok, bok bool
    
    ca, a, aok = component(a)
    cb, b, bok = component(b)
    
    if !aok && !bok {
      return true       // both sides are finished, we must have matched (even if both sides are empty)
    }else if !aok && bok {
      return wc || enc  // the left side is finished but the right is not; if the left is in a wildcard or if encompassing is enabled it matches the remaining right components
    }else if aok && !bok {
      return false      // the right side ran out of components to match but the left did not
    }
    
    // note our wildcard state
    wc = ca == wildcard
    // if left is either a wildcard or matches right we continue, otherwise we cannot match
    if !wc && !strings.EqualFold(ca, cb) {
      return false
    }
    
  }
  
  // if we haven't returned by this point, no match
  return false
}

/**
 * Determine if a string matches this QName
 */
func (q QName) MatchesString(a string) bool {
  return q.Matches(QName(a))
}

/**
 * Determine if a string matches this QName
 */
func (q QName) MatchesStringWithOptions(a string, options int) bool {
  return q.MatchesWithOptions(QName(a), options)
}

/**
 * Obtain the first component in a QName if one exists. Returns
 * the component, the remaining QName string, and whether or not
 * the component is valid.
 */
func component(q string) (string, string, bool) {
  l := len(q)
  c := ""
  
  if l < 1 {
    return c, "", false
  }
  
  for i := 0; i < l; {
    if i >= l {
      return c, q[i:], i > 0
    }
    r, w := utf8.DecodeRuneInString(q[i:])
    i += w
    if r == '.' {
      return c, q[i:], true
    }else{
      c += string(r)
    }
  }
  
  return c, "", true
}

