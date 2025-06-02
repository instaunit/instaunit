package assert

import (
	"fmt"
	"reflect"

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

// Assert equality
func Equal(e, a interface{}, m string, x ...interface{}) error {
	if !EqualValues(e, a) {
		return &AssertionError{e, a, fmt.Sprintf(m, x...)}
	} else {
		return nil
	}
}

// Are objects equal
func EqualObjects(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	} else {
		return reflect.DeepEqual(expected, actual)
	}
}

// Are objects exactly or semantically equal
func EqualValues(expected, actual interface{}) bool {
	if EqualObjects(expected, actual) {
		return true
	}

	actualType := reflect.TypeOf(actual)
	if actualType == nil {
		return false
	}

	expectedValue := reflect.ValueOf(expected)
	if expectedValue.IsValid() && expectedValue.Type().ConvertibleTo(actualType) {
		return reflect.DeepEqual(expectedValue.Convert(actualType).Interface(), actual)
	}

	return false
}
