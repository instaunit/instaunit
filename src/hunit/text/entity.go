package text

import (
  "fmt"
  "strings"
  "encoding/json"
)

/**
 * A type error
 */
type contentTypeError string

func (e contentTypeError) Error() string {
  return string(e)
}

func IsContentTypeError(e error) bool {
  _, ok := e.(contentTypeError)
  return ok
}

// formats entities
type entityFormatter func([]byte)([]byte, error)

// content types to formatters
var stdEntityFormatters = map[string]entityFormatter {
  "application/json": jsonEntityFormatter,
}

/**
 * Format an entity, if possible
 */
func FormatEntity(entity []byte, contentType string) ([]byte, error) {
  contentType = strings.TrimSpace(contentType)
  for k, f := range stdEntityFormatters {
    if MatchesContentType(k, contentType) {
      return f(entity)
    }
  }
  return nil, contentTypeError(fmt.Sprintf("Unsupported content type: %v", contentType))
}

/**
 * Format json
 */
func jsonEntityFormatter(entity []byte) ([]byte, error) {
  var v interface{}
  err := json.Unmarshal(entity, &v)
  if err != nil {
    return nil, err
  }
  return json.MarshalIndent(v, "", "  ")
}
