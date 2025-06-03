package hunit

import (
	"strings"
	"unicode"

	"github.com/instaunit/instaunit/hunit/assert"
	"github.com/instaunit/instaunit/hunit/entity"
	"github.com/instaunit/instaunit/hunit/runtime"
	"github.com/instaunit/instaunit/hunit/testcase"
)

// Compare entities for equality
func entitiesEqual(context runtime.Context, comparison testcase.Comparison, contentType string, expected []byte, actual interface{}) error {
	if comparison == testcase.CompareSemantic {
		return semanticEntitiesEqual(context, contentType, expected, actual)
	} else {
		return literalEntitiesEqual(context, contentType, expected, actual)
	}
}

// Compare entities for equality
func literalEntitiesEqual(context runtime.Context, contentType string, expected []byte, actual interface{}) error {
	var e, a interface{}
	var ok bool

	var abytes []byte
	if abytes, ok = actual.([]byte); !ok {
		return &assert.AssertionError{expected, actual, "Entities are not equal"}
	}

	if (context.Options & testcase.OptionEntityTrimTrailingWhitespace) == testcase.OptionEntityTrimTrailingWhitespace {
		e = strings.TrimRightFunc(string(expected), unicode.IsSpace)
		a = strings.TrimRightFunc(string(abytes), unicode.IsSpace)
	} else {
		e = expected
		a = abytes
	}

	if !assert.EqualValues(e, a) {
		return &assert.AssertionError{e, a, "Entities are not equal"}
	} else {
		return nil
	}
}

// Compare entities for equality
func semanticEntitiesEqual(context runtime.Context, contentType string, expected []byte, actual interface{}) error {

	e, err := entity.Unmarshal(contentType, expected)
	if err != nil {
		return err
	}

	if !entity.SemanticEqual(e, actual) {
		return &assert.AssertionError{e, actual, "Entities are not equal"}
	} else {
		return nil
	}
}
