// 
// Copyright (c) 2015-2016 Brian W. Wolter, All rights reserved.
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
  "strconv"
)

var excerptCallout excerptFormatter

/**
 * A span formatter
 */
type excerptFormatter struct {
}

/**
 * Format an excerpt
 */
func (f excerptFormatter) FormatExcerpt(s span) string {
  var e string
  
  l := len(s.text)
  c := "^"
  
  // find the beginning of the current line and the line number
  n := int64(1) // base 1
  a := s.offset
  for ; a > 0; a-- {
    if s.text[a] == '\n' {
      a++ // right after the newline
      n++ // newline
      break
    }
  }
  
  // count the rest of the newlines
  for i := a; i > 0; i-- {
    if s.text[i] == '\n' {
      n++ // newline
    }
  }
  
  // find the end of the current line
  z := s.offset
  for z < l {
    if s.text[z] == '\n' {
      break
    }else{
      z++
    }
  }
  
  // add the excerpt line
  x := strconv.FormatInt(n, 10)
  e = x +": "+ s.text[a:z] + "\n"
  lx := len(x) + len(": ")
  
  // add the callout
  p := s.offset - a
  for i := 0; i < (p + lx); i++ {
    e += " "
  }
  for i := 0; i < s.length; i++ {
    e += c
  }
  
  return e + "\n"
}
