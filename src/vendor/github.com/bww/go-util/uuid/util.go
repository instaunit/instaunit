package uuid

import (
  "regexp"
)

var uuidPattern = regexp.MustCompile("[a-fA-f0-9]{8}-([a-fA-f0-9]{4}-){3}[a-fA-f0-9]{12}")

// Determine if an input string resembles a UUID
func ResemblesUUID(s string) bool {
  return len(s) == 36 && uuidPattern.MatchString(s)
}
