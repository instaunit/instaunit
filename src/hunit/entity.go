package hunit

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/instaunit/instaunit/hunit/runtime"
	"github.com/instaunit/instaunit/hunit/test"
)

// Compare entities for equality
func entitiesEqual(context runtime.Context, comparison test.Comparison, contentType string, expected []byte, actual interface{}) error {
	if comparison == test.CompareSemantic {
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
		return &AssertionError{expected, actual, "Entities are not equal"}
	}

	if (context.Options & test.OptionEntityTrimTrailingWhitespace) == test.OptionEntityTrimTrailingWhitespace {
		e = strings.TrimRightFunc(string(expected), unicode.IsSpace)
		a = strings.TrimRightFunc(string(abytes), unicode.IsSpace)
	} else {
		e = expected
		a = abytes
	}

	if !equalValues(e, a) {
		return &AssertionError{e, a, "Entities are not equal"}
	} else {
		return nil
	}
}

// Compare entities for equality
func semanticEntitiesEqual(context runtime.Context, contentType string, expected []byte, actual interface{}) error {

	e, err := unmarshalEntity(context, contentType, expected)
	if err != nil {
		return err
	}

	if !semanticEqual(e, actual) {
		return &AssertionError{e, actual, "Entities are not equal"}
	} else {
		return nil
	}
}

// Unmarshal an entity
func unmarshalEntity(context runtime.Context, contentType string, entity []byte) (interface{}, error) {

	// trim off the parameters following ';' if we have any
	if i := strings.Index(contentType, ";"); i > 0 {
		contentType = contentType[:i]
	}

	switch contentType {
	case "application/json":
		return unmarshalJSONEntity(context, entity)
	case "text/csv":
		return unmarshalCSVEntity(context, entity)
	default:
		return nil, fmt.Errorf("Unsupported content type for semantic comparison: %v", contentType)
	}

}

// Unmarshal a JSON entity
func unmarshalJSONEntity(context runtime.Context, entity []byte) (interface{}, error) {
	if entity == nil || len(entity) < 1 {
		return nil, nil
	}
	var value interface{}
	err := json.Unmarshal(entity, &value)
	if err != nil {
		return nil, fmt.Errorf("Invalid JSON entity: %v", err)
	}
	return value, nil
}

// Unmarshal a CSV entity
func unmarshalCSVEntity(context runtime.Context, entity []byte) (interface{}, error) {
	if entity == nil || len(entity) < 1 {
		return nil, nil
	}

	value := make([]interface{}, 0)
	var h []string

	r := csv.NewReader(bytes.NewBuffer(entity))
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("Invalid CSV entity: %v (%s)", err, strings.Join(row, ", "))
		}

		if h == nil {
			h = row
			continue
		}

		m := make(map[string]string)
		for i, e := range row {
			m[h[i]] = e
		}
		value = append(value, m)
	}

	return value, nil
}

// Compare results
func semanticEqual(expected, actual interface{}) bool {
	switch a := actual.(type) {

	case map[string]string:
		e, ok := expected.(map[string]string)
		if !ok {
			return false
		}
		for k, v := range e {
			if !semanticEqual(v, a[k]) {
				return false
			}
		}

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
		if len(a) != len(e) {
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
