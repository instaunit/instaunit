package script

import (
	"fmt"
	"testing"

	"github.com/instaunit/instaunit/hunit/expr"
	"github.com/stretchr/testify/assert"
)

var (
	invalidInputError = fmt.Errorf("Invalid input")
)

var context = expr.Variables{
	"a": 123,
	"b": "String value",
	"c": map[string]interface{}{
		"a": []string{"Zero", "One", "Two", "Three"},
		"b": false,
	},
}

func TestScripts(t *testing.T) {
	tests := []struct {
		Script Script
		Expect bool
		Error  error
	}{
		{
			Script{"js", `1 == 2`},
			false,
			nil,
		},
		{
			Script{"js", `2 == 2`},
			true,
			nil,
		},
		{
			Script{"js", `false`},
			false,
			nil,
		},
		{
			Script{"js", `"anything"`},
			true,
			nil,
		},
		{
			Script{"js", `c.b`},
			false,
			nil,
		},
		{
			Script{"js", `a == b`},
			false,
			nil,
		},
		{
			Script{"js", `a == a`},
			true,
			nil,
		},
		{
			Script{"js", `c.a[0] == "Zero"`},
			true,
			nil,
		},
		//
		{
			Script{"epl", `1 == 2`},
			false,
			nil,
		},
		{
			Script{"epl", `2 == 2`},
			true,
			nil,
		},
		{
			Script{"epl", `false`},
			false,
			nil,
		},
		{
			Script{"epl", `"Not a bool"`},
			false,
			ErrInvalidTypeError{`"Not a bool"`, "bool", "Not a bool"},
		},
		{
			Script{"epl", `c.b`},
			false,
			nil,
		},
		{
			Script{"epl", `a == b`},
			false,
			nil,
		},
		{
			Script{"epl", `a == a`},
			true,
			nil,
		},
		{
			Script{"epl", `c.a[0] == "Zero"`},
			true,
			nil,
		},
	}
	for _, e := range tests {
		fmt.Println(">>>", e.Script)
		r, err := e.Script.Bool(context)
		if e.Error != nil {
			assert.Equal(t, e.Error, err)
		} else if assert.Nil(t, err, fmt.Sprint(err)) {
			assert.Equal(t, e.Expect, r)
		}
	}
}
