package entity

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/instaunit/instaunit/hunit/assert"
	"github.com/instaunit/instaunit/hunit/httputil/mimetype"
)

var ErrUnsupported = errors.New("Entity type is not supported")

// Unmarshal an entity
func Unmarshal(contentType string, entity []byte) (interface{}, error) {
	// trim off the parameters following ';' if we have any
	if i := strings.Index(contentType, ";"); i > 0 {
		contentType = contentType[:i]
	}

	switch contentType {
	case mimetype.JSON:
		return unmarshalJSON(entity)
	case mimetype.CSV:
		return unmarshalCSV(entity)
	default:
		return nil, fmt.Errorf("%w: %v", ErrUnsupported, contentType)
	}
}

// Unmarshal a JSON entity
func unmarshalJSON(entity []byte) (interface{}, error) {
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
func unmarshalCSV(entity []byte) (interface{}, error) {
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
func SemanticEqual(expected, actual interface{}) bool {
	switch a := actual.(type) {

	case map[string]string:
		e, ok := expected.(map[string]string)
		if !ok {
			return false
		}
		for k, v := range e {
			if !SemanticEqual(v, a[k]) {
				return false
			}
		}

	case map[string]interface{}:
		e, ok := expected.(map[string]interface{})
		if !ok {
			return false
		}
		for k, v := range e {
			if !SemanticEqual(v, a[k]) {
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
			if !SemanticEqual(v, a[i]) {
				return false
			}
		}

	default:
		return assert.EqualValues(expected, actual)

	}
	return true
}
