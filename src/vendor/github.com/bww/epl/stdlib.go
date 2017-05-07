// 
// Copyright (c) 2015 Brian William Wolter, All rights reserved.
// EPL - A little Embeddable Predicate Language
// 
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
// 
//   * Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
// 
//   * Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//     
//   * Neither the names of Brian William Wolter, Wolter Group New York, nor the
//     names of its contributors may be used to endorse or promote products derived
//     from this software without specific prior written permission.
//     
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
// INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
// OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
// OF THE POSSIBILITY OF SUCH DAMAGE.
// 

package epl

import (
  "os"
  "fmt"
  "regexp"
  "reflect"
)

/**
 * The environment context
 */
type environment struct {}

/**
 * Obtain an environment variable
 */
func (e environment) Variable(name string) (interface{}, error) {
  v := os.Getenv(name)
  if v != "" {
    return v, nil
  }else{
    return nil, undefinedVariableError
  }
}

/**
 * Standard library
 */
var stdlib = map[string]interface{}{
  "env": environment{},
  "len": builtInLen,
  "match": builtInMatch,
  "printf": builtInPrintf,
}

/**
 * len()
 */
func builtInLen(a interface{}) (int, error) {
  v, _ := derefValue(reflect.ValueOf(a))
  switch v.Kind() {
    case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan, reflect.String:
      return v.Len(), nil
    default:
      return 0, fmt.Errorf("Invalid parameter for builtin 'len'")
  }
}

/**
 * Regex match
 */
func builtInMatch(e, v interface{}) (bool, error) {
  var se, sv string
  var ok bool
  if se, ok = e.(string); !ok {
    return false, fmt.Errorf("Invalid parameter to: match(string, string)")
  }
  if sv, ok = v.(string); !ok {
    return false, fmt.Errorf("Invalid parameter to: match(string, string)")
  }
  return regexp.MatchString(se, sv)
}

/**
 * Print and return true
 */
func builtInPrintf(s *State, f interface{}, a ...interface{}) (bool, error) {
  var sf string
  var ok bool
  if sf, ok = f.(string); !ok {
    return false, fmt.Errorf("Invalid parameter to: printf(string, ...<any>)")
  }
  fmt.Fprintf(s.Runtime.Stdout, sf +"\n", a...)
  return true, nil
}
