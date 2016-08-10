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
  "fmt"
)

/**
 * A parser
 */
type parser struct {
  scanner   *scanner
  la        []token
}

/**
 * Create a parser
 */
func newParser(s *scanner) *parser {
  return &parser{s, make([]token, 0, 2)}
}

/**
 * Obtain a look-ahead token without consuming it
 */
func (p *parser) peek(n int) token {
  var t token
  
  l := len(p.la)
  if n < l {
    return p.la[n]
  }else if n + 1 > cap(p.la) {
    panic(fmt.Errorf("Look-ahead overrun: %d >= %d", n + 1, cap(p.la)))
  }
  
  p.la = p.la[:n+1]
  for i := l; i < n + 1; i++ {
    t = p.scanner.scan()
    p.la[i] = t
  }
  
  return t
}

/**
 * Consume the next token
 */
func (p *parser) next() token {
  l := len(p.la)
  if l < 1 {
    return p.scanner.scan()
  }else{
    t := p.la[0]
    for i := 1; i < l; i++ { p.la[i-1] = p.la[i] }
    p.la = p.la[:l-1]
    return t
  }
}

/**
 * Consume the next token asserting that it is one of the provided token types
 */
func (p *parser) nextAssert(valid ...tokenType) (token, error) {
  t := p.next()
  switch t.which {
    case tokenEOF:
      return token{}, fmt.Errorf("Unexpected end-of-input")
    case tokenError:
      return token{}, fmt.Errorf("Error: %v", t)
  }
  for _, v := range valid {
    if t.which == v {
      return t, nil
    }
  }
  return token{}, invalidTokenError(t, valid...)
}

/**
 * Parse
 */
func (p *parser) parse() (*Program, error) {
  e, err := p.parseExpression()
  if err != nil {
    return nil, err
  }else if t := p.peek(0); t.which != tokenEOF {
    return nil, fmt.Errorf("Syntax error: %v", t)
  }else{
    return &Program{e}, nil
  }
}

/**
 * Parse
 */
func (p *parser) parseExpression() (executable, error) {
  return p.parseLogicalOr()
}

/**
 * Parse a logical or
 */
func (p *parser) parseLogicalOr() (executable, error) {
  
  left, err := p.parseLogicalAnd()
  if err != nil {
    return nil, err
  }
  
  op := p.peek(0)
  switch op.which {
    case tokenError:
      return nil, fmt.Errorf("Error: %v", op)
    case tokenLogicalOr:
      break // valid token
    default:
      return left, nil
  }
  
  p.next() // consume the operator
  right, err := p.parseLogicalOr()
  if err != nil {
    return nil, err
  }
  
  return &logicalOrNode{node{encompass(op.span, left.src(), right.src()), &op}, left, right}, nil
}

/**
 * Parse a logical and
 */
func (p *parser) parseLogicalAnd() (executable, error) {
  
  left, err := p.parseRelational()
  if err != nil {
    return nil, err
  }
  
  op := p.peek(0)
  switch op.which {
    case tokenError:
      return nil, fmt.Errorf("Error: %v", op)
    case tokenLogicalAnd:
      break // valid token
    default:
      return left, nil
  }
  
  p.next() // consume the operator
  right, err := p.parseLogicalAnd()
  if err != nil {
    return nil, err
  }
  
  return &logicalAndNode{node{encompass(op.span, left.src(), right.src()), &op}, left, right}, nil
}

/**
 * Parse a relational expression
 */
func (p *parser) parseRelational() (executable, error) {
  
  left, err := p.parseArithmeticL1()
  if err != nil {
    return nil, err
  }
  
  op := p.peek(0)
  switch op.which {
    case tokenError:
      return nil, fmt.Errorf("Error: %v", op)
    case tokenLess, tokenGreater, tokenEqual, tokenLessEqual, tokenGreaterEqual, tokenNotEqual:
      break // valid tokens
    default:
      return left, nil
  }
  
  p.next() // consume the operator
  right, err := p.parseRelational()
  if err != nil {
    return nil, err
  }
  
  return &relationalNode{node{encompass(op.span, left.src(), right.src()), &op}, op, left, right}, nil
}

/**
 * Parse an arithmetic expression
 */
func (p *parser) parseArithmeticL1() (executable, error) {
  
  left, err := p.parseArithmeticL2()
  if err != nil {
    return nil, err
  }
  
  op := p.peek(0)
  switch op.which {
    case tokenError:
      return nil, fmt.Errorf("Error: %v", op)
    case tokenAdd, tokenSub:
      break // valid tokens
    default:
      return left, nil
  }
  
  p.next() // consume the operator
  right, err := p.parseArithmeticL1()
  if err != nil {
    return nil, err
  }
  
  return &arithmeticNode{node{encompass(op.span, left.src(), right.src()), &op}, op, left, right}, nil
}

/**
 * Parse an arithmetic expression
 */
func (p *parser) parseArithmeticL2() (executable, error) {
  
  left, err := p.parseDeref()
  if err != nil {
    return nil, err
  }
  
  op := p.peek(0)
  switch op.which {
    case tokenError:
      return nil, fmt.Errorf("Error: %v", op)
    case tokenMul, tokenDiv, tokenMod:
      break // valid tokens
    default:
      return left, nil
  }
  
  p.next() // consume the operator
  right, err := p.parseArithmeticL2()
  if err != nil {
    return nil, err
  }
  
  return &arithmeticNode{node{encompass(op.span, left.src(), right.src()), &op}, op, left, right}, nil
}

/**
 * Parse a deref expression
 */
func (p *parser) parseDeref() (executable, error) {
  
  left, err := p.parseIndex()
  if err != nil {
    return nil, err
  }
  
  op := p.peek(0)
  switch op.which {
    case tokenError:
      return nil, fmt.Errorf("Error: %v", op)
    case tokenDot:
      break // valid token
    default:
      return left, nil
  }
  
  p.next() // consume the operator
  right, err := p.parseDeref()
  if err != nil {
    return nil, err
  }
  
  switch v := right.(type) {
    case *identNode, *derefNode, *indexNode:
      return &derefNode{node{encompass(op.span, left.src()), &op}, left, v}, nil
    default:
      return nil, fmt.Errorf("Expected ident, deref or subscript: (%T) %v\n%v", right, right, excerptCallout.FormatExcerpt(right.src()))
  }
  
}

/**
 * Parse an index expression
 */
func (p *parser) parseIndex() (executable, error) {
  
  left, err := p.parsePrimary()
  if err != nil {
    return nil, err
  }
  
  return p.parseSubscript(left)
}

/**
 * Parse an index expression
 */
func (p *parser) parseSubscript(left executable) (executable, error) {
  
  op := p.peek(0)
  switch op.which {
    case tokenError:
      return nil, fmt.Errorf("Error: %v", op)
    case tokenLBracket:
      break // valid token
    default:
      return left, nil
  }
  
  p.next() // consume the '['
  right, err := p.parseExpression()
  if err != nil {
    return nil, err
  }
  
  t, err := p.nextAssert(tokenRBracket)
  if err != nil {
    return nil, err
  }
  
  return p.parseSubscript(&indexNode{node{encompass(op.span, left.src(), right.src(), t.span), &op}, left, right})
}

/**
 * Parse a primary expression
 */
func (p *parser) parsePrimary() (executable, error) {
  t := p.next()
  switch t.which {
    case tokenEOF:
      return nil, fmt.Errorf("Unexpected end-of-input")
    case tokenError:
      return nil, fmt.Errorf("Error: %v", t)
    case tokenLParen:
      return p.parseParen()
    case tokenIdentifier:
      return &identNode{node{t.span, &t}, t.value.(string)}, nil
    case tokenNumber, tokenString:
      return &literalNode{node{t.span, &t}, t.value}, nil
    case tokenTrue:
      return &literalNode{node{t.span, &t}, true}, nil
    case tokenFalse:
      return &literalNode{node{t.span, &t}, false}, nil
    case tokenNil:
      return &literalNode{node{t.span, &t}, nil}, nil
    default:
      return nil, fmt.Errorf("Illegal token in primary expression: %v", t)
  }
}

/**
 * Parse a (sub-expression)
 */
func (p *parser) parseParen() (executable, error) {
  
  e, err := p.parseExpression()
  if err != nil {
    return nil, err
  }
  
  t := p.next()
  if t.which != tokenRParen {
    return nil, fmt.Errorf("Expected ')' but found %v", t)
  }
  
  return e, nil
}

/**
 * A parser error
 */
type parserError struct {
  message   string
  span      span
  cause     error
}

/**
 * Error
 */
func (e parserError) Error() string {
  if e.cause != nil {
    return fmt.Sprintf("%s: %v\n%v", e.message, e.cause, excerptCallout.FormatExcerpt(e.span))
  }else{
    return fmt.Sprintf("%s\n%v", e.message, excerptCallout.FormatExcerpt(e.span))
  }
}

/**
 * Invalid token error
 */
func invalidTokenError(t token, e ...tokenType) error {
  
  m := fmt.Sprintf("Invalid token: %v", t.which)
  if e != nil && len(e) > 0 {
    m += " (expected: "
    for i, t := range e {
      if i > 0 { m += ", " }
      m += fmt.Sprintf("%v", t)
    }
    m += ")"
  }
  
  return &parserError{m, t.span, nil}
}
