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
  "io"
  "fmt"
  "reflect"
)

var undefinedVariableError = fmt.Errorf("undefined")

/**
 * A runtime context
 */
type Context interface {
  Variable(name string)(interface{}, error)
}

/**
 * A variable provider
 */
type VariableProvider func(name string)(interface{}, error)

/**
 * Executable context
 */
type Runtime struct {
  Stdout    io.Writer
}

/**
 * Executable context
 */
type context struct {
  stack []interface{}
}

/**
 * Create a new context
 */
func newContext(f ...interface{}) *context {
  return &context{f}
}

/**
 * Push a frame
 */
func (c *context) push(f interface{}) *context {
  if l := len(c.stack); cap(c.stack) > l {
    c.stack = c.stack[:l+1]
    c.stack[l] = f
  }else{
    c.stack = append(c.stack, f)
  }
  return c
}

/**
 * Pop a frame
 */
func (c *context) pop() *context {
  if l := len(c.stack); l > 0 {
    c.stack = c.stack[:l-1]
  }
  return c
}

/**
 * Obtain the top of the stack
 */
func (c *context) top() interface{} {
  if l := len(c.stack); l > 0 {
    return c.stack[l-1]
  }else{
    return nil
  }
}

/**
 * Obtain a value
 */
func (c *context) get(r *Runtime, s span, n string) (interface{}, error) {
  return c.sget(r, s, n, c.stack, derefOptionDerefFunctions)
}

/**
 * Obtain a value, without dereferencing
 */
func (c *context) value(r *Runtime, s span, n string) (interface{}, error) {
  return c.sget(r, s, n, c.stack, 0)
}

/**
 * Obtain a value
 */
func (c *context) sget(r *Runtime, s span, n string, k []interface{}, opts derefOptions) (interface{}, error) {
  l := len(k)
  if l < 1 {
    return nil, nil
  }
  
  v, err := derefProp(r, c, s, k[l-1], n, opts)
  if err == undefinedVariableError && l > 1 {
    return c.sget(r, s, n, k[:l-1], opts)
  }else if err != nil {
    return nil, err
  }
  
  return v, nil
}

/**
 * Execution state
 */
type State struct {
  Runtime   *Runtime
  Context   interface{}
}

var typeOfState   = reflect.TypeOf(&State{})
var typeOfRuntime = reflect.TypeOf(&Runtime{})
var typeOfError   = reflect.TypeOf((*error)(nil)).Elem()

/**
 * Print options
 */
type PrintOptions int
const (
  PrintOptionNone = 0
)

/**
 * A single indent level
 */
const indentLevel = "  "

/**
 * Print state
 */
type printState struct {
  depth   int
}

/**
 * Descend
 */
func (s printState) Desc() printState {
  return printState{s.depth + 1}
}

/**
 * Indent for the current state depth
 */
func (s printState) Indent() string {
  var v string
  for i := 0; i < s.depth; i++ {
    v += indentLevel
  }
  return v
}

/**
 * An executable expression
 */
type executable interface {
  src()(span)
  exec(*Runtime, *context)(interface{}, error)
  print(io.Writer, PrintOptions, printState)(error)
}

/**
 * Variable arguments
 */
type varargs []executable

/**
 * Obtain the node src span
 */
func (n varargs) src() span {
  if l := len(n); l > 1 {
    return encompass(n[0].src(), n[l-1].src())
  }else if l > 0 {
    return n[0].src()
  }else{
    return span{}
  }
}

/**
 * Execute
 */
func (n varargs) exec(runtime *Runtime, context *context) (interface{}, error) {
  args := make([]interface{}, len(n))
  for i, e := range n {
    v, err := e.exec(runtime, context)
    if err != nil {
      return nil, err
    }
    args[i] = v
  }
  return args, nil
}

/**
 * Print
 */
func (n varargs) print(w io.Writer, opts PrintOptions, state printState) error {
  for _, e := range n {
    err := e.print(w, opts, state)
    if err != nil {
      return err
    }
  }
  return nil
}
 
/**
 * An AST node
 */
type node struct {
  span      span
  token     *token
}

/**
 * Obtain the node src span
 */
func (n node) src() span {
  return n.span
}

/**
 * Execute
 */
func (n *node) exec(runtime *Runtime, context *context) (interface{}, error) {
  return nil, fmt.Errorf("No implementation")
}

/**
 * Print
 */
func (n *node) print(w io.Writer, opts PrintOptions, state printState) error {
  return fmt.Errorf("No implementation")
}

/**
 * A program
 */
type Program struct {
  root executable
}

/**
 * Execute
 */
func (p *Program) Exec(context interface{}) (interface{}, error) {
  return p.root.exec(&Runtime{os.Stdout}, newContext(stdlib, context))
}

/**
 * Print
 */
func (p *Program) Print(w io.Writer, opts PrintOptions) error {
  return p.root.print(w, opts, printState{})
}

/**
 * A program
 */
type emptyNode struct {
  node
}

/**
 * A logical OR node
 */
type logicalOrNode struct {
  node
  left, right executable
}

/**
 * Execute
 */
func (n *logicalOrNode) exec(runtime *Runtime, context *context) (interface{}, error) {
  
  lvi, err := n.left.exec(runtime, context)
  if err != nil {
    return nil, err
  }
  lv, err := asBool(n.left.src(), lvi)
  if err != nil {
    return nil, err
  }
  
  if lv {
    return true, nil
  }
  
  rvi, err := n.right.exec(runtime, context)
  if err != nil {
    return nil, err
  }
  rv, err := asBool(n.right.src(), rvi)
  if err != nil {
    return nil, err
  }
  
  return rv, nil
}

/**
 * Print
 */
func (n *logicalOrNode) print(w io.Writer, opts PrintOptions, state printState) error {
  indent := state.Indent()
  
  _, err := w.Write([]byte(indent + fmt.Sprintf("%T (\n", n)))
  if err != nil {
    return err
  }
  
  n.left.print(w, opts, state.Desc())
  
  _, err = w.Write([]byte("\n"+ indent +"||\n"))
  if err != nil {
    return err
  }
  
  n.right.print(w, opts, state.Desc())
  
  _, err = w.Write([]byte("\n"+ indent +")\n"))
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * A logical AND node
 */
type logicalAndNode struct {
  node
  left, right executable
}

/**
 * Execute
 */
func (n *logicalAndNode) exec(runtime *Runtime, context *context) (interface{}, error) {
  
  lvi, err := n.left.exec(runtime, context)
  if err != nil {
    return nil, err
  }
  lv, err := asBool(n.left.src(), lvi)
  if err != nil {
    return nil, err
  }
  
  if !lv {
    return false, nil
  }
  
  rvi, err := n.right.exec(runtime, context)
  if err != nil {
    return nil, err
  }
  rv, err := asBool(n.right.src(), rvi)
  if err != nil {
    return nil, err
  }
  
  return rv, nil
}

/**
 * Print
 */
func (n *logicalAndNode) print(w io.Writer, opts PrintOptions, state printState) error {
  indent := state.Indent()
  
  _, err := w.Write([]byte(indent + fmt.Sprintf("%T (\n", n)))
  if err != nil {
    return err
  }
  
  n.left.print(w, opts, state.Desc())
  
  _, err = w.Write([]byte("\n"+ indent +"&&\n"))
  if err != nil {
    return err
  }
  
  n.right.print(w, opts, state.Desc())
  
  _, err = w.Write([]byte("\n"+ indent +")\n"))
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * An arithmetic expression node
 */
type arithmeticNode struct {
  node
  op          token
  left, right executable
}

/**
 * Execute
 */
func (n *arithmeticNode) exec(runtime *Runtime, context *context) (interface{}, error) {
  
  lvi, err := n.left.exec(runtime, context)
  if err != nil {
    return nil, err
  }
  lv, err := asNumber(n.left.src(), lvi)
  if err != nil {
    return nil, err
  }
  
  rvi, err := n.right.exec(runtime, context)
  if err != nil {
    return nil, err
  }
  rv, err := asNumber(n.right.src(), rvi)
  if err != nil {
    return nil, err
  }
  
  switch n.op.which {
    case tokenAdd:
      return lv + rv, nil
    case tokenSub:
      return lv - rv, nil
    case tokenMul:
      return lv * rv, nil
    case tokenDiv:
      return lv / rv, nil
    case tokenMod: // truncates to int
      return int64(lv) % int64(rv), nil
    default:
      return nil, fmt.Errorf("Invalid operator: %v", n.op)
  }
  
}

/**
 * Print
 */
func (n *arithmeticNode) print(w io.Writer, opts PrintOptions, state printState) error {
  indent := state.Indent()
  
  _, err := w.Write([]byte(indent + fmt.Sprintf("%T (\n", n)))
  if err != nil {
    return err
  }
  
  n.left.print(w, opts, state.Desc())
  
  var op string
  switch n.op.which {
    case tokenAdd:
      op = "+"
    case tokenSub:
      op = "-"
    case tokenMul:
      op = "*"
    case tokenDiv:
      op = "/"
    case tokenMod:
      op = "%"
    default:
      return fmt.Errorf("Invalid operator: %v", n.op)
  }
  
  _, err = w.Write([]byte("\n"+ indent + op +"\n"))
  if err != nil {
    return err
  }
  
  n.right.print(w, opts, state.Desc())
  
  _, err = w.Write([]byte("\n"+ indent +")\n"))
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * An relational expression node
 */
type relationalNode struct {
  node
  op          token
  left, right executable
}

/**
 * Execute
 */
func (n *relationalNode) exec(runtime *Runtime, context *context) (interface{}, error) {
  
  lvi, err := n.left.exec(runtime, context)
  if err != nil {
    return nil, err
  }
  rvi, err := n.right.exec(runtime, context)
  if err != nil {
    return nil, err
  }
  
  switch n.op.which {
    case tokenEqual:
      return lvi == rvi, nil
    case tokenNotEqual:
      return lvi != rvi, nil
  }
  
  lv, err := asNumber(n.left.src(), lvi)
  if err != nil {
    return nil, err
  }
  rv, err := asNumber(n.right.src(), rvi)
  if err != nil {
    return nil, err
  }
  
  switch n.op.which {
    case tokenLess:
      return lv < rv, nil
    case tokenGreater:
      return lv > rv, nil
    case tokenLessEqual:
      return lv <= rv, nil
    case tokenGreaterEqual:
      return lv >= rv, nil
    default:
      return nil, fmt.Errorf("Invalid operator: %v", n.op)
  }
  
}

/**
 * Print
 */
func (n *relationalNode) print(w io.Writer, opts PrintOptions, state printState) error {
  indent := state.Indent()
  
  _, err := w.Write([]byte(indent + fmt.Sprintf("%T (\n", n)))
  if err != nil {
    return err
  }
  
  n.left.print(w, opts, state.Desc())
  
  var op string
  switch n.op.which {
    case tokenLess:
      op = "<"
    case tokenGreater:
      op = ">"
    case tokenLessEqual:
      op = "<="
    case tokenGreaterEqual:
      op = ">="
    default:
      return fmt.Errorf("Invalid operator: %v", n.op)
  }
  
  _, err = w.Write([]byte("\n"+ indent + op +"\n"))
  if err != nil {
    return err
  }
  
  n.right.print(w, opts, state.Desc())
  
  _, err = w.Write([]byte("\n"+ indent +")\n"))
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * A dereference expression node
 */
type derefNode struct {
  node
  left, right executable
}

/**
 * Execute
 */
func (n *derefNode) exec(runtime *Runtime, context *context) (interface{}, error) {
  
  v, err := n.left.exec(runtime, context)
  if err != nil {
    return nil, err
  }
  
  context.push(v)
  defer context.pop()
  
  var z interface{}
  switch v := n.right.(type) {
    case *identNode:
      z, err = context.get(runtime, n.span, v.ident)
    case *derefNode, *indexNode, *invokeNode:
      z, err = v.exec(runtime, context)
    default:
      return nil, fmt.Errorf("Invalid right operand to . (dereference): %v (%T)", v, v)
  }
  
  return z, err
}

/**
 * Print
 */
func (n *derefNode) print(w io.Writer, opts PrintOptions, state printState) error {
  indent := state.Indent()
  
  _, err := w.Write([]byte(indent + fmt.Sprintf("%T (\n", n)))
  if err != nil {
    return err
  }
  
  n.left.print(w, opts, state.Desc())
  
  _, err = w.Write([]byte("\n"+ indent +".\n"))
  if err != nil {
    return err
  }
  
  n.right.print(w, opts, state.Desc())
  
  _, err = w.Write([]byte("\n"+ indent +")\n"))
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * A dereference expression node
 */
type indexNode struct {
  node
  left, right executable
}

/**
 * Execute
 */
func (n *indexNode) exec(runtime *Runtime, context *context) (interface{}, error) {
  
  left, err := n.left.exec(runtime, context)
  if err != nil {
    return nil, err
  }
  
  sub, err := n.right.exec(runtime, context)
  if err != nil {
    return nil, err
  }
  
  val := reflect.ValueOf(sub)
  if val.Kind() == reflect.Invalid {
    return nil, runtimeErrorf(n.right.src(), "Subscript expression is nil")
  }
  
  context.push(left)
  defer context.pop()
  
  deref, _ := derefValue(reflect.ValueOf(left))
  switch deref.Kind() {
    case reflect.String: // character index
      return n.execString(runtime, context, deref, val)
    case reflect.Array:
      return n.execArray(runtime, context, deref, val)
    case reflect.Slice:
      return n.execArray(runtime, context, deref, val)
    case reflect.Map:
      return n.execMap(runtime, context, deref, val)
    default:
      return nil, runtimeErrorf(n.span, "Expression result is not indexable: %v", displayType(deref))
  }
  
}

/**
 * Execute
 */
func (n *indexNode) execString(runtime *Runtime, context *context, val reflect.Value, index reflect.Value) (interface{}, error) {
  val = reflect.ValueOf([]rune(val.Interface().(string))) // convert to []rune
  
  i, err := asNumberValue(n.right.src(), index)
  if err != nil {
    return nil, err
  }
  
  l := val.Len()
  if int(i) < 0 || int(i) >= l {
    return nil, runtimeErrorf(n.span, "Index out-of-bounds: %v", i)
  }
  
  return string(val.Index(int(i)).Interface().(rune)), nil
}

/**
 * Execute
 */
func (n *indexNode) execArray(runtime *Runtime, context *context, val reflect.Value, index reflect.Value) (interface{}, error) {
  
  i, err := asNumberValue(n.right.src(), index)
  if err != nil {
    return nil, err
  }
  
  l := val.Len()
  if int(i) < 0 || int(i) >= l {
    return nil, runtimeErrorf(n.span, "Index out-of-bounds: %v", i)
  }
  
  return val.Index(int(i)).Interface(), nil
}

/**
 * Execute
 */
func (n *indexNode) execMap(runtime *Runtime, context *context, val reflect.Value, key reflect.Value) (interface{}, error) {
  
  if !key.Type().AssignableTo(val.Type().Key()) {
    return nil, runtimeErrorf(n.span, "Expression result is not assignable to map key type: %v != %v", key.Type(), val.Type().Key())
  }
  
  return val.MapIndex(key).Interface(), nil
}

/**
 * Print
 */
func (n *indexNode) print(w io.Writer, opts PrintOptions, state printState) error {
  indent := state.Indent()
  
  _, err := w.Write([]byte(indent + fmt.Sprintf("%T (\n", n)))
  if err != nil {
    return err
  }
  
  n.left.print(w, opts, state.Desc())
  
  _, err = w.Write([]byte("\n"+ indent +"[\n"))
  if err != nil {
    return err
  }
  
  n.right.print(w, opts, state.Desc())
  
  _, err = w.Write([]byte("\n"+ indent +"])\n"))
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * A function invocation expression node
 */
type invokeNode struct {
  node
  left, right executable
  params      []executable
}

/**
 * Execute
 */
func (n *invokeNode) exec(runtime *Runtime, context *context) (interface{}, error) {
  var liv interface{}
  var err error
  
  var name string
  switch v := n.right.(type) {
    case *identNode:
      name = v.ident
    default:
      return nil, runtimeErrorf(n.span, "Invalid node type for function call: %T", v)
  }
  
  if n.left != nil {
    liv, err = n.left.exec(runtime, context)
    if err != nil {
      return nil, err
    }
  }
  
  return invokeFunction(runtime, context, n.span, liv, name, n.params)
}

/**
 * Print
 */
func (n *invokeNode) print(w io.Writer, opts PrintOptions, state printState) error {
  indent := state.Indent()
  
  _, err := w.Write([]byte(indent + fmt.Sprintf("%T (\n", n)))
  if err != nil {
    return err
  }
  
  if(n.left != nil){
    n.left.print(w, opts, state.Desc())
  }else{
    _, err = w.Write([]byte(indent + indentLevel +"<nil>"))
    if err != nil {
      return err
    }
  }
  
  _, err = w.Write([]byte("\n"+ indent +".\n"))
  if err != nil {
    return err
  }
  
  n.right.print(w, opts, state.Desc())
  
  _, err = w.Write([]byte("\n"+ indent +"(\n"))
  if err != nil {
    return err
  }
  
  for i, p := range n.params {
    if i > 0 {
      _, err = w.Write([]byte(indent +",\n"))
      if err != nil {
        return err
      }
    }
    p.print(w, opts, state.Desc())
  }
  
  _, err = w.Write([]byte("\n"+ indent +"))\n"))
  if err != nil {
    return err
  }
  
  return nil
}

/**
 * An identifier expression node
 */
type identNode struct {
  node
  ident string
}

/**
 * Execute
 */
func (n *identNode) exec(runtime *Runtime, context *context) (interface{}, error) {
  return context.get(runtime, n.span, n.ident)
}

/**
 * Print
 */
func (n *identNode) print(w io.Writer, opts PrintOptions, state printState) error {
  _, err := w.Write([]byte(state.Indent() + "ident:"+ n.ident))
  if err != nil {
    return err
  }
  return nil
}

/**
 * A literal expression node
 */
type literalNode struct {
  node
  value interface{}
}

/**
 * Execute
 */
func (n *literalNode) exec(runtime *Runtime, context *context) (interface{}, error) {
  return n.value, nil
}

/**
 * Print
 */
func (n *literalNode) print(w io.Writer, opts PrintOptions, state printState) error {
  _, err := w.Write([]byte(state.Indent() + fmt.Sprintf("literal<%T>:%v", n.value, n.value)))
  if err != nil {
    return err
  }
  return nil
}

/**
 * Obtain an interface value as a bool
 */
func asBool(s span, value interface{}) (bool, error) {
  v := reflect.ValueOf(value)
  switch v.Kind() {
    case reflect.Bool:
      return v.Bool(), nil
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
      return v.Int() != 0, nil
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
      return v.Uint() != 0, nil
    case reflect.Float32, reflect.Float64:
      return v.Float() != 0, nil
    default:
      return false, runtimeErrorf(s, "Cannot cast %v to bool", displayType(v))
  }
}

/**
 * Obtain an interface value as a number
 */
func asNumber(s span, v interface{}) (float64, error) {
  return asNumberValue(s, reflect.ValueOf(v))
}

/**
 * Obtain an interface value as a number
 */
func asNumberValue(s span, v reflect.Value) (float64, error) {
  switch v.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
      return float64(v.Int()), nil
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
      return float64(v.Uint()), nil
    case reflect.Float32, reflect.Float64:
      return v.Float(), nil
    default:
      return 0, runtimeErrorf(s, "Cannot cast %v to numeric", displayType(v))
  }
}

/**
 * Invoke a function
 */
func invokeFunction(runtime *Runtime, context *context, s span, liv interface{}, name string, ins []executable) (interface{}, error) {
  
  var f reflect.Value
  if liv != nil {
    lrv := reflect.ValueOf(liv)
    f = lrv.MethodByName(name)
    if !f.IsValid() {
      return nil, runtimeErrorf(s, "No such method '%v' for type %v or method is not exported", name, lrv.Type())
    }
  }else{
    var err error
    liv, err = context.value(runtime, s, name)
    if err == undefinedVariableError {
      return nil, undefinedVariableError
    }else if liv == nil {
      return nil, runtimeErrorf(s, "No such function '%v'", name)
    }
    f = reflect.ValueOf(liv)
    if f.Kind() != reflect.Func {
      return nil, runtimeErrorf(s, "Variable '%v' (%T) is not a function", name, displayType(f))
    }
  }
  
  ft := f.Type()
  lp := len(ins)
  
  cout := ft.NumOut()
  if cout > 2 {
    return nil, runtimeErrorf(s, "Function %v returns %v values (expected: 0, 1 or 2)", name, cout)
  }
  
  in, extra := 0, 1
  args := make([]reflect.Value, 0)
  cin := ft.NumIn()
  
  if ft.IsVariadic() {
    off := 0
    if ft.In(in) == typeOfState {
      off++
    }
    if lp < cin - off - 1 {
      return nil, runtimeErrorf(s, "Function %v takes %v arguments but is given %v", name, cin - off - 1, lp)
    }
    vargs := ins[cin - off - 1:]
    vatmp := make([]executable, len(ins) - len(vargs))
    copy(vatmp, ins[:cin - off - 1])
    ins = append(vatmp, varargs(vargs))
    lp  = len(ins)
  }
  
  if cin != lp {
    if cin - extra != lp /* allow for runtime parameter */ {
      return nil, runtimeErrorf(s, "Function %v takes %v arguments but is given %v", name, cin, lp)
    }
    if ft.In(in) != typeOfState {
      return nil, runtimeErrorf(s, "Function %v takes %v arguments but is given %v; first native argument must receive %v", name, cin - extra, lp, typeOfState)
    }
    args = append(args, reflect.ValueOf(&State{runtime, context}))
    in++
  }
  
  for _, e := range ins {
    v, err := e.exec(runtime, context)
    if err != nil {
      return nil, err
    }
    t := ft.In(in)
    var a reflect.Value
    if v == nil { // we need a typed zero value if the value is nil
      a = reflect.Zero(t)
    }else{
      a = reflect.ValueOf(v)
    }
    if !a.IsValid() {
      return nil, runtimeErrorf(e.src(), "Invalid parameter")
    }
    if !a.Type().AssignableTo(t) {
      return nil, runtimeErrorf(e.src(), "Cannot use %v as %v", displayType(a), t.String())
    }
    args = append(args, a)
    in++
  }
  
  var r []reflect.Value
  if ft.IsVariadic() {
    r = f.CallSlice(args)
  }else{
    r = f.Call(args)
  }
  if r == nil {
    return nil, nil
  }else if l := len(r); l > 2 {
    return nil, runtimeErrorf(s, "Function %v must return either (void), (interface{}) or (interface{}, error)", name)
  }else if l == 0 {
    return nil, nil
  }else if l == 1 {
    if ft.Out(0) == typeOfError {
      if !r[0].IsNil() {
        return nil, r[0].Interface().(error)
      }else{
        return nil, nil
      }
    }else{
      return r[0].Interface(), nil
    }
  }else if l == 2 {
    r0 := r[0].Interface()
    r1 := r[1].Interface()
    if r1 == nil {
      return r0, nil
    }else if e, ok := r1.(error); ok {
      return r0, e
    }else{
      return nil, runtimeErrorf(s, "Function %v must return either (void), (interface{}) or (interface{}, error)", name)
    }
  }
  
  return nil, undefinedVariableError
}

/**
 * Dereference options
 */
type derefOptions int
const (
  derefOptionNone           = derefOptions(0)
  derefOptionDerefFunctions = derefOptions(1 << 0)
)


/**
 * Dereference
 */
func derefProp(runtime *Runtime, context *context, s span, val interface{}, ident string, opts derefOptions) (interface{}, error) {
  
  switch v := val.(type) {
    case Context:
      return v.Variable(ident)
    case VariableProvider:
      return v(ident)
    case func(string)(interface{}, error):
      return v(ident)
    case map[string]interface{}:
      res, ok := v[ident]
      if ok {
        return res, nil
      }else{
        return nil, undefinedVariableError
      }
  }
  
  switch v := reflect.ValueOf(val); v.Kind() {
    case reflect.Map:
      return derefMap(s, v, ident)
    case reflect.Ptr, reflect.Struct:
      return derefMember(runtime, context, s, val, ident, opts)
    default:
      return nil, runtimeErrorf(s, "Cannot dereference variable: %v", displayType(v))
  }
  
}

/**
 * Execute
 */
func derefMap(s span, val reflect.Value, property string) (interface{}, error) {
  key := reflect.ValueOf(property)
  
  if !key.Type().AssignableTo(val.Type().Key()) {
    return nil, runtimeErrorf(s, "Expression result is not assignable to map key type: %v != %v", key.Type(), val.Type().Key())
  }
  
  res := val.MapIndex(key)
  if res.IsValid() {
    return res.Interface(), nil
  }else{
    return nil, undefinedVariableError
  }
}

/**
 * Execute
 */
func derefMember(runtime *Runtime, context *context, s span, val interface{}, property string, opts derefOptions) (interface{}, error) {
  raw  := reflect.ValueOf(val)
  base := raw
  
  if base.Kind() == reflect.Ptr {
    base, _ = derefValue(base)
  }
  if base.Kind() != reflect.Struct {
    return nil, runtimeErrorf(s, "Cannot dereference variable: %v", displayType(base))
  }
  
  v := raw.MethodByName(property)
  if v.IsValid() {
    if (opts & derefOptionDerefFunctions) == derefOptionDerefFunctions {
      f := v.Type()
      if n := f.NumOut(); n < 1 {
        return nil, runtimeErrorf(s, "Method %v of %v returns no values, which cannot be used as a dereference", v, displayType(base))
      }
      if f.Out(0) == typeOfError {
        return nil, runtimeErrorf(s, "Method %v of %v returns only an error, which cannot be used as a dereference", v, displayType(base))
      }
      return invokeFunction(runtime, context, s, val, property, nil)
    }else{
      if !v.CanInterface() {
        return nil, runtimeErrorf(s, "Cannot access %v of %v", property, displayType(raw))
      }
      return v.Interface(), nil
    }
  }
  
  v = base.FieldByName(property)
  if v.IsValid() {
    if !v.CanInterface() {
      return nil, runtimeErrorf(s, "Cannot access %v of %v", property, displayType(raw))
    }
    return v.Interface(), nil
  }
  
  return nil, undefinedVariableError
}

/**
 * Dereference a value
 */
func derefValue(value reflect.Value) (reflect.Value, int) {
  v := value
  c := 0
  for ; v.Kind() == reflect.Ptr; {
    v = v.Elem()
    c++
  }
  return v, c
}

/**
 * Obtain the presentation type of a value
 */
func displayType(v reflect.Value) string {
  if v.Kind() == reflect.Invalid {
    return "<nil>"
  }else{
    return v.Type().String()
  }
}

/**
 * A runtime error
 */
type runtimeError struct {
  message   string
  span      span
  cause     error
}

/**
 * Format a runtime error
 */
func runtimeErrorf(s span, f string, a ...interface{}) *runtimeError {
  return &runtimeError{fmt.Sprintf(f, a...), s, nil}
}

/**
 * Error
 */
func (e runtimeError) Error() string {
  if e.cause != nil {
    return fmt.Sprintf("%s: %v\n%v", e.message, e.cause, excerptCallout.FormatExcerpt(e.span))
  }else{
    return fmt.Sprintf("%s\n%v", e.message, excerptCallout.FormatExcerpt(e.span))
  }
}
