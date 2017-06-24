package text

import (
  "unicode"
)

/**
 * Normalize text
 */
func NormalizeString(s string) string {
  var n string
  
  var insp bool
  for _, e := range s {
    if unicode.IsSpace(e) {
      if len(n) > 0 {
        insp = true
      }
    }else{
      if unicode.IsLetter(e) {
        if insp { n += " " }
        n += string(unicode.ToLower(e))
        insp = false
      }else if unicode.IsDigit(e) || e == '-' || e == '_' || e == '\'' || e == ',' {
        if insp { n += " " }
        n += string(e)
        insp = false
      }else if len(n) > 0 {
        insp = true
      }
    }
  }
  
  return n
}

/**
 * Normalize query terms
 */
func NormalizeTerms(s string) string {
  var n string
  
  var insp bool
  for _, e := range s {
    if unicode.IsSpace(e) {
      if len(n) > 0 {
        insp = true
      }
    }else{
      if unicode.IsLetter(e) {
        if insp { n += " " }
        n += string(unicode.ToLower(e))
        insp = false
      }else if unicode.IsDigit(e) || e == '-' || e == '_' || e == '\'' {
        if insp { n += " " }
        n += string(e)
        insp = false
      }else if len(n) > 0 {
        insp = true
      }
    }
  }
  
  return n
}

/**
 * Collapse whitespace
 */
func CollapseSpaces(s string) string {
  var n string
  
  var insp bool
  for _, e := range s {
    if unicode.IsSpace(e) {
      if len(n) > 0 {
        insp = true
      }
    }else{
      if insp { n += " " }
      n += string(e)
      insp = false
    }
  }
  
  return n
}

/**
 * Normalize a list, using a special final delimiter between the last
 * two elements.
 */
func NormalizeJoin(l []string, d, f string) string {
  n := len(l)
  var s string
  for i, e := range l {
    if i > 0 {
      if n - (i + 1) == 0 {
        s += f
      }else{
        s += d
      }
    }
    s += e
  }
  return s
}
