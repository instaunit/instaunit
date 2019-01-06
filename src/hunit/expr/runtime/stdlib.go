package runtime

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

import (
	"github.com/bww/epl"
	"github.com/bww/go-util/rand"
	"github.com/bww/go-util/uuid"
)

// The standard library
type stdlib struct{}

// Builtins
var Stdlib stdlib

// Generate a random string
func (s stdlib) RandomString(n float64) string {
	return rand.RandomString(int(n))
}

// Generate a random name (Docker style: <adjective>_<noun>)
func (s stdlib) RandomIdent() string {
	l, r := DockerName()
	return fmt.Sprintf("%s_%s", l, r)
}

// Generate a random person name
func (s stdlib) RandomPersonName() string {
	return fmt.Sprintf("%s %s", strings.Title(firstName()), strings.Title(lastName()))
}

// Generate a random person first name
func (s stdlib) RandomFirstName() string {
	return strings.Title(firstName())
}

// Generate a random scientist last name
func (s stdlib) RandomLastName() string {
	return strings.Title(lastName())
}

// Generate a random name
func (s stdlib) RandomCompanyName() string {
	return companyName()
}

// Generate a random (v4) UUID
func (s stdlib) RandomUUID() string {
	return uuid.New().String()
}

// Escape a URL query component
func (s stdlib) QueryEscape(v string) string {
	return url.QueryEscape(v)
}

// Unescape a URL component
func (s stdlib) QueryUnescape(v string) (string, error) {
	return url.QueryUnescape(v)
}

// Convert a string to lowercase
func (s stdlib) ToLower(v string) string {
	return strings.ToLower(v)
}

// Convert a string to UPPERCASE
func (s stdlib) ToUpper(v string) string {
	return strings.ToUpper(v)
}

// Convert a string to Title Case
func (s stdlib) ToTitle(v string) string {
	return strings.ToTitle(v)
}

// Trim space from both ends of a string
func (s stdlib) TrimSpace(v string) string {
	return strings.TrimSpace(v)
}

// Take any element from a collection
func (s stdlib) Any(v interface{}) (interface{}, error) {
	val := reflect.Indirect(reflect.ValueOf(v))
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		if val.Len() > 0 {
			return val.Index(0).Interface(), nil
		}
	case reflect.Chan:
		if x, ok := val.TryRecv(); ok {
			return x.Interface(), nil
		}
	case reflect.String:
		if s := v.(string); len(s) > 0 {
			return s[0], nil
		}
	default:
		return nil, fmt.Errorf("Invalid type: %T", v)
	}
	return nil, nil // no elements
}

// Filter elements
func (s stdlib) Filter(v interface{}, f string) ([]interface{}, error) {
	val := reflect.Indirect(reflect.ValueOf(v))
	switch val.Kind() {
	case reflect.Array, reflect.Slice: // Cool
	default:
		return nil, fmt.Errorf("Invalid type: %T", v)
	}

	prg, err := epl.Compile(f)
	if err != nil {
		return nil, fmt.Errorf("Invalid filter:\n%v", err)
	}

	sub := make([]interface{}, 0)
	for i := 0; i < val.Len(); i++ {
		e := val.Index(i)
		x := e.Interface()
		res, err := prg.Exec(x)
		if err != nil {
			return nil, fmt.Errorf("Could not filter:\n%v", err)
		}
		match, ok := res.(bool)
		if !ok {
			return nil, fmt.Errorf("Invalid filter result type: %T", res)
		}
		if match {
			sub = append(sub, x)
		}
	}

	return sub, nil
}
