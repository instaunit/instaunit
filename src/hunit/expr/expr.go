package expr

import (
  "os"
  "fmt"
  "strings"
  "hunit/expr/runtime"
  "github.com/bww/epl"
)

// Variables
type Variables map[string]interface{}

// Map the environment
func mapenv(v []string) map[string]string {
  env := make(map[string]string)
  for _, e := range v {
    if x := strings.Index(e, "="); x > 0 {
      env[e[:x]] = e[x+1:]
    }else{
      env[e] = ""
    }
  }
  return env
}

// Produce a context with the standard library included
func runtimeContext(v Variables, e []string) interface{} {
  c := make(Variables)
  for k, v := range v {
    c[k] = v
  }
  c["std"] = runtime.Stdlib
  c["env"] = mapenv(e)
  return c
}

// Interpolate expressions in a string
func Interpolate(s string, v Variables) (string, error) {
  return interpolate(s, "${", "}", runtimeContext(v, os.Environ()))
}

// Interpolate
func interpolate(s, pre, suf string, context interface{}) (string, error) {
  if len(pre) < 1 || len(suf) < 1 {
    return "", fmt.Errorf("Invalid variable prefix/suffix")
  }
  
  fp := pre[0]
  fs := suf[0]
  
  var out string
  var i, esc int
  for {
    if i >= len(s) {
      break
    }
    
    if s[i] == '\\' {
      esc++
      if (esc % 2) == 0 {
        out += "\\"
      }
      i++
      continue
    }
    
    if s[i] == fp && (esc % 2) == 0 && matchAhead(s[i:], pre) {
      i += len(pre)
      start := i
      for {
        if i >= len(s) {
          return "", fmt.Errorf("Unexpected end-of-input")
        }
        if s[i] == fs && matchAhead(s[i:], suf) {
          
          prg, err := epl.Compile(s[start:i])
          if err != nil {
            return "", err
          }
          
          res, err := prg.Exec(context)
          if err != nil {
            return "", fmt.Errorf("Could not evaluate expression: {%v}: %v", s[start:i], err)
          }
          
          switch v := res.(type) {
            case string:
              out += v
            default:
              out += fmt.Sprintf("%v", v)
          }
          
          i += len(suf)
          break
        }else{
          i++
        }
      }
    }else{
      out += string(s[i])
      i++
    }
    
    esc = 0
  }
  
  return out, nil
}

// Match ahead
func matchAhead(s, x string) bool {
  if len(s) < len(x) {
    return false
  }
  for i := 0; i < len(x); i++ {
    if s[i] != x[i] {
      return false
    }
  }
  return true
}
