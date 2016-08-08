package hunit

import (
  "fmt"
  "strings"
  "encoding/json"
)

/**
 * Compare entities for equality
 */
func entitiesEqual(context Context, comparison Comparison, contentType string, expected, actual []byte) error {
  if comparison == CompareSemantic {
    return semanticEntitiesEqual(context, contentType, expected, actual)
  }else{
    return literalEntitiesEqual(context, contentType, expected, actual)
  }
}

/**
 * Compare entities for equality
 */
func literalEntitiesEqual(context Context, contentType string, expected, actual []byte) error {
  var e, a interface{}
  
  if (context.Options & OptionEntityTrimTrailingWhitespace) == OptionEntityTrimTrailingWhitespace {
    e = strings.TrimRight(string(expected), whitespace)
    a = strings.TrimRight(string(actual), whitespace)
  }else{
    e = expected
    a = actual
  }
  
  if !equalValues(e, a) {
    return &AssertionError{e, a, "Entities are not equal"}
  }else{
    return nil
  }
}

/**
 * Compare entities for equality
 */
func semanticEntitiesEqual(context Context, contentType string, expected, actual []byte) error {
  
  e, err := unmarshalEntity(context, contentType, expected)
  if err != nil {
    return err
  }
  
  a, err := unmarshalEntity(context, contentType, actual)
  if err != nil {
    return err
  }
  
  if !semanticEqual(e, a) {
    return &AssertionError{e, a, "Entities are not equal"}
  }else{
    return nil
  }
}

/**
 * Unmarshal an entity
 */
func unmarshalEntity(context Context, contentType string, entity []byte) (interface{}, error) {
  
  // trim off the parameters following ';' if we have any
  if i := strings.Index(contentType, ";"); i > 0 {
    contentType = contentType[:i]
  }
  
  switch contentType {
    case "application/json":
      return unmarshalJSONEntity(context, entity)
    default:
      return nil, fmt.Errorf("Unsupported content type for semantic comparison: %v", contentType)
  }
  
}

/**
 * Unmarshal a JSON entity
 */
func unmarshalJSONEntity(context Context, entity []byte) (interface{}, error) {
  if entity == nil || len(entity) < 1 {
    return nil, nil
  }
  var value interface{}
  err := json.Unmarshal(entity, &value)
  if err != nil {
    return nil, err
  }
  return value, nil
}

/**
 * Compare results
 */
func semanticEqual(expected, actual interface{}) bool {
  switch a := actual.(type) {
    
    case map[string]interface{}:
      e, ok := expected.(map[string]interface{})
      if !ok {
        return false
      }
      for k, v := range e {
        if !semanticEqual(v, a[k]) {
          return false
        }
      }
      
    case []interface{}:
      e, ok := expected.([]interface{})
      if !ok {
        return false
      }
      for i, v := range e {
        if !semanticEqual(v, a[i]) {
          return false
        }
      }
      
    default:
      return equalValues(expected, actual)
    
  }
  return true
}
