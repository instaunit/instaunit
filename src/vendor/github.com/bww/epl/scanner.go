// 
// Copyright (c) 2014 Brian William Wolter, All rights reserved.
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
// --
// 
// This scanner incorporates routines from the Go package text/scanner:
// http://golang.org/src/pkg/text/scanner/scanner.go
// 
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// 
// http://golang.org/LICENSE
// 

package epl

import (
  "fmt"
  "math"
  "strings"
  "strconv"
  "unicode"
  "unicode/utf8"
)

/**
 * A text span
 */
type span struct {
  text      string
  offset    int
  length    int
}

/**
 * Span (unquoted) excerpt
 */
func (s span) excerpt() string {
  max := float64(len(s.text))
  return s.text[int(math.Max(0, math.Min(max, float64(s.offset)))):int(math.Min(max, float64(s.offset+s.length)))]
}

/**
 * Span (quoted) excerpt
 */
func (s span) String() string {
  return strconv.Quote(s.excerpt())
}

/**
 * Create a new span that encompasses all the provided spans. The underlying text is taken from the first span.
 */
func encompass(a ...span) span {
  var t string
  min, max := 0, 0
  for i, e := range a {
    if i == 0 {
      min, max = e.offset, e.offset + e.length
      t = e.text
    }else{
      if e.offset < min {
        min = e.offset
      }
      if e.offset + e.length > max {
        max = e.offset + e.length
      }
    }
  }
  return span{t, min, max - min}
}

/**
 * Numeric type
 */
type numericType int

const (
  numericInteger numericType = iota
  numericFloat
)

/**
 * Token type
 */
type tokenType int

/**
 * Token types
 */
const (
  
  tokenError tokenType  = iota
  tokenEOF
  tokenBlock
  
  tokenString
  tokenNumber
  tokenIdentifier
  
  tokenTrue
  tokenFalse
  tokenNil
  
  tokenLParen           = '('
  tokenRParen           = ')'
  tokenLBracket         = '['
  tokenRBracket         = ']'
  tokenDot              = '.'
  tokenComma            = ','
  tokenSemi             = ';'
  tokenAdd              = '+'
  tokenSub              = '-'
  tokenMul              = '*'
  tokenDiv              = '/'
  tokenMod              = '%'
  tokenAssign           = '='
  tokenLess             = '<'
  tokenGreater          = '>'
  tokenBang             = '!'
  tokenAmp              = '&'
  tokenPipe             = '|'
  
  tokenPrefixAdd        = 1 << 16
  tokenInc              = tokenPrefixAdd | '+'
  
  tokenPrefixSub        = 1 << 17
  tokenDec              = tokenPrefixSub | '-'
  
  tokenPrefixAmp        = 1 << 18
  tokenLogicalAnd       = tokenPrefixAmp | '&'
  
  tokenPrefixPipe       = 1 << 19
  tokenLogicalOr        = tokenPrefixPipe | '|'
  
  tokenSuffixEqual      = 1 << 20
  tokenEqual            = tokenSuffixEqual | '='
  tokenAddEqual         = tokenSuffixEqual | '+'
  tokenSubEqual         = tokenSuffixEqual | '-'
  tokenMulEqual         = tokenSuffixEqual | '*'
  tokenDivEqual         = tokenSuffixEqual | '/'
  tokenModEqual         = tokenSuffixEqual | '%'
  tokenLessEqual        = tokenSuffixEqual | '<'
  tokenGreaterEqual     = tokenSuffixEqual | '>'
  tokenNotEqual         = tokenSuffixEqual | '!'
  tokenAssignInfer      = tokenSuffixEqual | ':'
  tokenLogicalAndEqual  = tokenSuffixEqual | '&'
  tokenLogicalOrEqual   = tokenSuffixEqual | '|'
  
)

/**
 * Token type string
 */
func (t tokenType) String() string {
  switch t {
    case tokenError:
      return "Error"
    case tokenEOF:
      return "EOF"
    case tokenBlock:
      return "{...}"
    case tokenString:
      return "String"
    case tokenNumber:
      return "Number"
    case tokenIdentifier:
      return "Ident"
    case tokenTrue:
      return "true"
    case tokenFalse:
      return "false"
    case tokenNil:
      return "nil"
    default:
      if t < 128 {
        return fmt.Sprintf("'%v'", string(t))
      }else{
        return fmt.Sprintf("%U", t)
      }
  }
}

/**
 * Token stuff
 */
const (
  eof = -1
)

/**
 * A token
 */
type token struct {
  span      span
  which     tokenType
  value     interface{}
}

/**
 * Stringer
 */
func (t token) String() string {
  switch t.which {
    case tokenError:
      return fmt.Sprintf("<%v %v %v>", t.which, t.span, t.value)
    default:
      return fmt.Sprintf("<%v %v>", t.which, t.span)
  }
}

/**
 * A scanner action
 */
type scannerAction func(*scanner) scannerAction

/**
 * A scanner error
 */
type scannerError struct {
  message   string
  span      span
  cause     error
}

/**
 * Error
 */
func (s *scannerError) Error() string {
  if s.cause != nil {
    return fmt.Sprintf("%s: %v\n%v", s.message, s.cause, excerptCallout.FormatExcerpt(s.span))
  }else{
    return fmt.Sprintf("%s\n%v", s.message, excerptCallout.FormatExcerpt(s.span))
  }
}

/**
 * A scanner
 */
type scanner struct {
  text      string
  index     int
  width     int // current rune width
  start     int // token start position
  depth     int // expression depth
  tokens    chan token
  state     scannerAction
}

/**
 * Create a scanner
 */
func newScanner(text string) *scanner {
  t := make(chan token, 5 /* several tokens may be produced in one iteration */)
  return &scanner{text, 0, 0, 0, 0, t, expressionAction}
}

/**
 * Scan and produce a token
 */
func (s *scanner) scan() token {
  for {
    select {
      case t := <- s.tokens:
        return t
      default:
        s.state = s.state(s)
    }
  }
}

/**
 * Create an error
 */
func (s *scanner) errorf(where span, cause error, format string, args ...interface{}) *scannerError {
  return &scannerError{fmt.Sprintf(format, args...), where, cause}
}

/**
 * Emit a token
 */
func (s *scanner) emit(t token) {
  s.tokens <- t
  s.start = t.span.offset + t.span.length
}

/**
 * Emit an error and return a nil action
 */
func (s *scanner) error(err *scannerError) scannerAction {
  s.tokens <- token{err.span, tokenError, err}
  return nil
}

/**
 * Obtain the next rune from input without consuming it
 */
func (s *scanner) peek() rune {
  r := s.next()
  s.backup()
  return r
}

/**
 * Consume the next rune from input
 */
func (s *scanner) next() rune {
  
  if s.index >= len(s.text) {
    s.width = 0
    return eof
  }
  
  r, w := utf8.DecodeRuneInString(s.text[s.index:])
  s.index += w
  s.width  = w
  
  return r
}

/**
 * Match ahead
 */
func (s *scanner) match(text string) bool {
  return s.matchAt(s.index, text)
}

/**
 * Match ahead
 */
func (s *scanner) matchAt(index int, text string) bool {
  i := index
  
  if i < 0 {
    return false
  }
  
  for n := 0; n < len(text); {
    
    if i >= len(s.text) {
      return false
    }
    
    r, w := utf8.DecodeRuneInString(s.text[i:])
    i += w
    c, z := utf8.DecodeRuneInString(text[n:])
    n += z
    
    if r != c {
      return false
    }
    
  }
  
  return true
}

/**
 * Match ahead. The shortest matching string in the set will succeed.
 */
func (s *scanner) matchAny(texts ...string) (bool, string) {
  return s.matchAnyAt(s.index, texts...)
}

/**
 * Match ahead. The shortest matching string in the set will succeed.
 */
func (s *scanner) matchAnyAt(index int, texts ...string) (bool, string) {
  i := index
  m := 0
  
  if i < 0 {
    return false, ""
  }
  
  for _, v := range texts {
    w := len(v)
    if w > m {
      m = w
    }
  }
  
  for n := 0; n < m; {
    for _, text := range texts {
      
      if i >= len(s.text) {
        return false, ""
      }
      if n >= len(text) {
        continue
      }
      
      r, w := utf8.DecodeRuneInString(s.text[i:])
      i += w
      c, z := utf8.DecodeRuneInString(text[n:])
      n += z
      
      if r != c {
        continue
      }
      if n >= len(text) {
        return true, text
      }
      
    }
  }
  
  return false, ""
}

/**
 * Find the next occurance of any character in the specified string
 */
func (s *scanner) findFrom(index int, any string, invert bool) int {
  i := index
  if !invert {
    return strings.IndexAny(s.text[i:], any)
  }else{
    for {
      
      if i >= len(s.text) {
        return -1
      }
      
      r, w := utf8.DecodeRuneInString(s.text[i:])
      
      if !strings.ContainsRune(any, r) {
        return i
      }else{
        i += w
      }
      
    }
  }
}

/**
 * Shuffle the token start to the current index
 */
func (s *scanner) ignore() {
  s.start = s.index
}

/**
 * Unconsume the previous rune from input (this can be called only once
 * per invocation of next())
 */
func (s *scanner) backup() {
  s.index -= s.width
}

/**
 * Skip past a rune that was previously peeked
 */
func (s *scanner) skip() {
  s.index += s.width
}

/**
 * Skip past a rune that was previously peeked and ignore it
 */
func (s *scanner) skipAndIgnore() {
  s.skip()
  s.ignore()
}

/**
 * Expresion action.
 */
func expressionAction(s *scanner) scannerAction {
  
  for {
    switch r := s.next(); {
      
      case r == eof:
        s.emit(token{span{s.text, len(s.text), 0}, tokenEOF, nil})
        return nil
        
      case unicode.IsSpace(r):
        s.ignore()
        
      case r == '"':
        // consume the open '"'
        return stringAction
        
      case r >= '0' && r <= '9':
        s.backup() // unget the first digit
        return numberAction
        
      case r == '_' || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z'):
        s.backup() // unget the first character
        return identifierAction
        
      case r == '(' || r == ')' || r == '[' || r == ']' || r == '.' || r == ',' || r == ';':
        s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(r), string(r)})
        return expressionAction
        
      case r == '&':
        if n := s.next(); n == '=' {
          s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(tokenSuffixEqual | r), string(r)})
        }else if n == '&' {
          s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(tokenPrefixAmp | r), string(r)})
        }else{
          s.backup()
          s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(r), string(r)})
        }
        return expressionAction
        
      case r == '|':
        if n := s.next(); n == '=' {
          s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(tokenSuffixEqual | r), string(r)})
        }else if n == '|' {
          s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(tokenPrefixPipe | r), string(r)})
        }else{
          s.backup()
          s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(r), string(r)})
        }
        return expressionAction
        
      case r == '+':
        if n := s.next(); n == '=' {
          s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(tokenSuffixEqual | r), string(r)})
        }else if n == '+' {
          s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(tokenPrefixAdd | r), string(r)})
        }else if n >= '0' && n <= '9' {
          s.backup() // unget the first digit; leave the sign
          return numberAction
        }else{
          s.backup()
          s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(r), string(r)})
        }
        return expressionAction
      
      case r == '-':
        if n := s.next(); n == '=' {
          s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(tokenSuffixEqual | r), string(r)})
        }else if n == '-' {
          s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(tokenPrefixSub | r), string(r)})
        }else if n >= '0' && n <= '9' {
          s.backup(); s.backup() // unget the sign and the first digit
          return numberAction
        }else{
          s.backup()
          s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(r), string(r)})
        }
        return expressionAction
      
      case r == '=' || r == '!' || r == '<' || r == '>' || r == ':' || r == '*' || r == '/' || r == '%':
        if n := s.next(); n == '=' {
          s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(tokenSuffixEqual | r), string(r)})
        }else{
          s.backup()
          s.emit(token{span{s.text, s.start, s.index - s.start}, tokenType(r), string(r)})
        }
        return expressionAction
        
      default:
        return s.error(s.errorf(span{s.text, s.index, 1}, nil, "Syntax error"))
        
    }
  }
  
  return expressionAction
}

/**
 * Quoted string
 */
func stringAction(s *scanner) scannerAction {
  if v, err := s.scanString('"', '\\'); err != nil {
    s.error(s.errorf(span{s.text, s.index, 1}, err, "Invalid string"))
  }else{
    s.emit(token{span{s.text, s.start, s.index - s.start}, tokenString, v})
  }
  return expressionAction
}

/**
 * Number string
 */
func numberAction(s *scanner) scannerAction {
  if v, _, err := s.scanNumber(); err != nil {
    s.error(s.errorf(span{s.text, s.index, 1}, err, "Invalid number"))
  }else{
    s.emit(token{span{s.text, s.start, s.index - s.start}, tokenNumber, v})
  }
  return expressionAction
}

/**
 * Identifier
 */
func identifierAction(s *scanner) scannerAction {
  
  v, err := s.scanIdentifierOrUUID()
  if err != nil {
    s.error(s.errorf(span{s.text, s.index, 1}, err, "Invalid identifier"))
  }
  
  t := span{s.text, s.start, s.index - s.start}
  switch v {
    case "true":
      s.emit(token{t, tokenTrue, v})
    case "false":
      s.emit(token{t, tokenFalse, v})
    case "nil":
      s.emit(token{t, tokenNil, nil})
    default:
      s.emit(token{t, tokenIdentifier, v})
  }
  
  return expressionAction
}

/***
 ***  SCANNING PRIMITIVES
 ***/

/**
 * Scan a delimited token with escape sequences. The opening delimiter is
 * expected to have already been consumed.
 */
func (s *scanner) scanString(quote, escape rune) (string, error) {
  var unquoted string
  
  for {
    switch r := s.next(); {
      
      case r == eof:
        return "", s.errorf(span{s.text, s.start, s.index - s.start}, nil, "Unexpected end-of-input")
        
      case r == escape:
        if e, err := s.scanEscape(quote, escape); err != nil {
          return "", s.errorf(span{s.text, s.start, s.index - s.start}, err, "Invalid escape sequence")
        }else{
          unquoted += string(e)
        }
        
      case r == quote:
        return unquoted, nil
        
      default:
        unquoted += string(r)
        
    }
  }
  
  return "", s.errorf(span{s.text, s.start, s.index - s.start}, nil, "Unexpected end-of-input")
}

/**
 * Scan an identifier
 */
func (s *scanner) scanIdentifier() (string, error) {
  start := s.index
	for r := s.next(); r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r); {
		r = s.next()
	}
	s.backup() // unget the last character
	return s.text[start:s.index], nil
}

/**
 * Scan an identifier or UUID
 */
func (s *scanner) scanIdentifierOrUUID() (string, error) {
  var start int
  
  u, _ := s.matchAny("U:", "u:")
  if u {
    s.next(); s.next() // skip "U:" or "u:"
    start = s.index
    
    n := 0
    for r := s.next(); r == '-' || (r >= 'A' && r <= 'F') || (r >= 'a' && r <= 'f' ) || unicode.IsDigit(r); {
      r = s.next()
      if r != '-' {
        n++
      }
    }
    
    if n != 32 {
      return "", s.errorf(span{s.text, start, s.index - start}, nil, "UUID identifier is incorrectly formatted")
    }
  }else{
    
    start = s.index
    for r := s.next(); r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r); {
      r = s.next()
    }
    
  }
  
  s.backup() // unget the last character
  return s.text[start:s.index], nil
}

/**
 * Scan a digit value
 */
func digitValue(ch rune) int {
	switch {
    case '0' <= ch && ch <= '9':
      return int(ch - '0')
    case 'a' <= ch && ch <= 'f':
      return int(ch - 'a' + 10)
    case 'A' <= ch && ch <= 'F':
      return int(ch - 'A' + 10)
	}
	return math.MaxInt32 // too big
}

/**
 * Decimla digit?
 */
func isDecimal(ch rune) bool {
  return '0' <= ch && ch <= '9'
}

/**
 * Scan digits
 */
func (s *scanner) scanDigits(base, n int) (string, error) {
  start := s.index
	for r := s.next(); n > 0 && digitValue(r) < base; {
		r = s.next(); n--
	}
  s.backup() // unget the stop character
	if n > 0 {
		return "", s.errorf(span{s.text, start, s.index - start}, nil, "Not enough digits")
	}else{
	  return s.text[start:s.index], nil
	}
}

/**
 * Scan digits
 */
func (s *scanner) scanDecimal(base, n int) (int64, error) {
  if d, err := s.scanDigits(base, n); err != nil {
    return 0, err
  }else{
    return strconv.ParseInt(d, base, 64)
  }
}

/**
 * Scan digits as a rune
 */
func (s *scanner) scanRune(base, n int) (rune, error) {
  if d, err := s.scanDecimal(base, n); err != nil {
    return 0, err
  }else{
    return rune(d), nil
  }
}

/**
 * Scan an escape
 */
func (s *scanner) scanEscape(quote, esc rune) (rune, error) {
  start := s.index
  r := s.next()
  switch r {
    case 'a':
      return '\a', nil
    case 'b':
      return '\b', nil
    case 'f':
      return '\f', nil
    case 'n':
      return '\n', nil
    case 'r':
      return '\r', nil
    case 't':
      return '\t', nil
    case 'v':
      return '\v', nil
    case esc, quote:
      return r, nil
    case '0', '1', '2', '3', '4', '5', '6', '7':
      return s.scanRune(8, 3)
    case 'x':
      return s.scanRune(16, 2)
    case 'u':
      return s.scanRune(16, 4)
    case 'U':
      return s.scanRune(16, 8)
    default:
      return 0, s.errorf(span{s.text, start, s.index - start}, nil, "Invalid escape sequence")
  }
}

/**
 * Scan a number mantissa
 */
func (s *scanner) scanMantissa(ch rune) rune {
  for isDecimal(ch) {
    ch = s.next()
  }
  return ch
}

/**
 * Scan a number fraction
 */
func (s *scanner) scanFraction(ch rune) rune {
  if ch == '.' {
    ch = s.scanMantissa(s.next())
  }
  return ch
}

/**
 * Scan a number exponent
 */
func (s *scanner) scanExponent(ch rune) rune {
  if ch == 'e' || ch == 'E' {
    ch = s.next()
    if ch == '-' || ch == '+' {
      ch = s.next()
    }
    ch = s.scanMantissa(ch)
  }
  return ch
}

/**
 * Scan a number
 */
func (s *scanner) scanNumber() (float64, numericType, error) {
  var isfloat bool
  start := s.index
  ch := s.next()
  
  if ch == '+' || ch == '-' {
    ch = s.next()
  }
  
  if ch == '0' {
    // int or float
    ch = s.next()
    if ch == 'x' || ch == 'X' {
      
      // hexadecimal int
      ch = s.next()
      
      hasMantissa := false
      for digitValue(ch) < 16 {
        ch = s.next()
        hasMantissa = true
      }
      
      // unscan the stop rune
      s.backup()
      
      if !hasMantissa {
        return 0, 0, s.errorf(span{s.text, start, s.index - start}, nil, "Illegal hexadecimal number")
      }
      
      if v, err := strconv.ParseInt(s.text[start+2:s.index], 16, 64); err != nil {
        return 0, 0, s.errorf(span{s.text, start, s.index - start}, err, "Could not parse number")
      }else{
        return float64(v), numericInteger, nil
      }
      
    } else {
      
      // octal int or float
      has8or9 := false
      for isDecimal(ch) {
        if ch > '7' { has8or9 = true }
        ch = s.next()
      }
      
      // check for a fraction
      if ch == '.' || ch == 'e' || ch == 'E' || ch == 'i' {
        goto fraction
      }
      
      // unscan the stop rune
      s.backup()
      
      // octal int
      if has8or9 {
        s.errorf(span{s.text, start, s.index - start}, nil, "Illegal octal number")
      }
      
      // parse our octal
      t := s.text[start:s.index]
      if v, err := strconv.ParseInt(t, 8, 64); err != nil {
        return 0, 0, s.errorf(span{s.text, start, s.index - start}, err, "Could not parse number")
      }else{
        return float64(v), numericInteger, nil
      }
      
    }
  }
  
  // decimal int or float
  ch = s.scanMantissa(ch)
  
fraction:
  
  // check for a fraction
  if ch == '.' {
    isfloat = true
    ch = s.scanFraction(ch)
  }
  
  // check for an exponent
  if ch == 'e' || ch == 'E' {
    isfloat = true
    ch = s.scanExponent(ch)
  }
  
  // unscan the non-numeric rune
  s.backup()
  
  // parse the base-10 number
  if isfloat {
    if v, err := strconv.ParseFloat(s.text[start:s.index], 64); err != nil {
      return 0, 0, s.errorf(span{s.text, start, s.index - start}, err, "Could not parse number")
    }else{
      return v, numericFloat, nil
    }
  }else{
    if v, err := strconv.ParseInt(s.text[start:s.index], 10, 64); err != nil {
      return 0, 0, s.errorf(span{s.text, start, s.index - start}, err, "Could not parse number")
    }else{
      return float64(v), numericInteger, nil
    }
  }
}

/*
func (s *scanner) scanComment(ch rune) rune {
	// ch == '/' || ch == '*'
	if ch == '/' {
		// line comment
		ch = s.next() // read character after "//"
		for ch != '\n' && ch >= 0 {
			ch = s.next()
		}
		return ch
	}

	// general comment
	ch = s.next() // read character after "/*"
	for {
		if ch < 0 {
			s.error("comment not terminated")
			break
		}
		ch0 := ch
		ch = s.next()
		if ch0 == '*' && ch == '/' {
			ch = s.next()
			break
		}
	}
	return ch
}
*/
