package hunit

import (
	"reflect"

	"github.com/instaunit/instaunit/hunit/script"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
)

// An assertion error
type AssertionError struct {
	Expected, Actual interface{}
	Message          string
}

func (e AssertionError) Error() string {
	m := e.Message
	if e.Message != "" {
		m += ":\n"
	}

	_, ek := typeAndKind(e.Expected)
	if ek == reflect.String || ek == reflect.Struct || ek == reflect.Map || ek == reflect.Slice || ek == reflect.Array {
		m += cmp.Diff(e.Expected, e.Actual)
	} else {
		m += "expected: " + spew.Sdump(e.Expected)
		m += "  actual: " + spew.Sdump(e.Actual)
	}

	return m
}

// A script error
type ScriptError struct {
	Message          string
	Expected, Actual interface{}
	Script           *script.Script
}

func (e ScriptError) Error() string {
	m := e.Message
	if m != "" {
		m += ":\n"
	}
	m += "expected: " + spew.Sdump(e.Expected)
	m += "  actual: " + spew.Sdump(e.Actual)
	m += "--\n" + e.Script.Source
	return m
}

// Obtain a value's type and kind
func typeAndKind(v interface{}) (reflect.Type, reflect.Kind) {
	t := reflect.TypeOf(v)
	k := t.Kind()
	for k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
	}
	return t, k
}
