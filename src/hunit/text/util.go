package text

/**
 * Return the first non-empty string
 */
func Coalesce(s ...string) string {
  for _, e := range s {
    if e != "" {
      return e
    }
  }
  return ""
}
