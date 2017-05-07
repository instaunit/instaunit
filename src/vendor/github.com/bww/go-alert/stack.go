// 
// Go Alert
// Copyright (c) 2015 Brian W. Wolter, All rights reserved.
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
//   * Neither the names of Brian W. Wolter nor the names of the contributors may
//     be used to endorse or promote products derived from this software without
//     specific prior written permission.
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

package alt

import (
  "fmt"
  "path"
  "strings"
  "runtime"
  "crypto/sha1"
  "encoding/json"
)

/**
 * A stacktrace frame
 */
type Frame struct {
  Filename    string
  LineNumber  int
  FilePath    string
  Function    string
  Module      string
}

/**
 * A stacktrace
 */
type Stacktrace struct {
  Frames []Frame
}

/**
 * Fingerprint
 */
func (s Stacktrace) Fingerprint() string {
  data, _ := json.Marshal(s)
  return fmt.Sprintf("%x", sha1.Sum(data))
}

/**
 * Generate a stacktrace
 */
func generateStacktrace() Stacktrace {
  return generateStacktraceWithOptions(1 /* skip this call itself */, nil)
}

/**
 * Generate a stacktrace with options
 */
func generateStacktraceWithOptions(skip int, exclude []string) Stacktrace {
  var stacktrace Stacktrace
  
  maxDepth := 10
  for depth := 1 /* skip this call itself */ + skip; depth < maxDepth; depth++ {
    
    pc, filePath, line, ok := runtime.Caller(depth)
    if !ok {
      break
    }
    
    f := runtime.FuncForPC(pc)
    fname := f.Name()
    
    if strings.Contains(fname, "runtime") {
      break // Stop when reaching runtime
    }
    if exclude != nil {
      for _, e := range exclude {
        if strings.Contains(fname, e) {
          continue // Skip excluded calls
        }
      }
    }
    
    var moduleName string
    if strings.Contains(f.Name(), "(") {
      components := strings.SplitN(f.Name(), ".(", 2)
      fname = "(" + components[1]
      moduleName = components[0]
    }
    
    fileName := path.Base(filePath)
    frame := Frame{Filename: fileName, LineNumber: line, FilePath: filePath, Function: fname, Module: moduleName}
    stacktrace.Frames = append(stacktrace.Frames, frame)
  }
  
  return stacktrace
}
