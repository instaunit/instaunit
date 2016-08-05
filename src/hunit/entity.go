package hunit

import (
  "strings"
)

/**
 * Compare entities for equality
 */
func entitiesEqual(context Context, expected, actual []byte) error {
  var e, a interface{}
  
  if (context.Options & OptionEntityCompareSemantically) == OptionEntityCompareSemantically {
    e = expected
    a = actual
  }else if (context.Options & OptionEntityTrimTrailingWhitespace) == OptionEntityTrimTrailingWhitespace {
    e = strings.TrimRight(string(expected), whitespace)
    a = strings.TrimRight(string(actual), whitespace)
  }
  
  if !equalValues(e, a) {
    return &AssertionError{e, a, "Entities are not equal"}
  }else{
    return nil
  }
}
