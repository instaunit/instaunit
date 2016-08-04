package hunit

import (
  "fmt"
  "reflect"
)

/**
 * Assert equality
 */
func assertEqual(e, a interface{}, m string, x ...interface{}) error {
  if !equalValues(e, a) {
    return &AssertionError{e, a, fmt.Sprintf(m, x...)}
  }else{
    return nil
  }
}

/**
 * Are objects equal
 */
func equalObjects(expected, actual interface{}) bool {
  if expected == nil || actual == nil {
    return expected == actual
  }else{
    return reflect.DeepEqual(expected, actual)
  }
}

/**
 * Are objects exactly or semantically equal
 */
func equalValues(expected, actual interface{}) bool {
  if equalObjects(expected, actual) {
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

