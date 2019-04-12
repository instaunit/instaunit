package text

import (
	"fmt"
)

// Convert a value to a reasonable string representation
func Stringer(v interface{}) string {
	switch c := v.(type) {
	case string:
		return c
	case *string:
		return *c
	case fmt.Stringer:
		return c.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
