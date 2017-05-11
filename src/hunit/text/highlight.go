package text

var highlights = map[string]string {
  "application/json":       "json",
  "application/javascript": "javascript",
}

/**
 * Obtain the appropriate highlight for the entity
 */
func EntityHighlight(c string) string {
  for k, v := range highlights {
    if MatchesContentType(k, c) {
      return v
    }
  }
  return ""
}
