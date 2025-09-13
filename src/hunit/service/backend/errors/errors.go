package errors

import (
	"fmt"
	"strings"
)

// Alternate errors describes a set of possible errors, one of which actually occurred. 🤷‍♂️
type AlternateErrors []error

func (errs AlternateErrors) Error() string {
	switch len(errs) {
	case 0:
		return "No error"
	case 1:
		return errs[0].Error()
	}

	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("One of %d possible errors occurred:\n", len(errs)))

	for i, err := range errs {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("#%d: %v\n", i+1, err))
	}

	return b.String()
}
