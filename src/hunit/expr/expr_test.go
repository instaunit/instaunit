package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var context = map[string]interface{}{
	"a": 123,
	"b": "String value",
	"c": map[string]interface{}{
		"a": []string{"Zero", "One", "Two", "Three"},
		"b": false,
	},
}

/**
 * Test interpolate
 */
func TestInterpolateVariables(t *testing.T) {
	tests := []struct {
		Text   string
		Expect string
		Error  error
	}{
		{
			`Before ${a}, after.`, `Before 123, after.`, nil,
		},
		{
			`Before ${"\}"}, after.`, `Before }, after.`, nil,
		},
		{
			`Before ${""}}, after.`, `Before }, after.`, nil,
		},
		{
			`Before \${a}, after.`, `Before ${a}, after.`, nil,
		},
		{
			`Before \\${a}, after.`, `Before \123, after.`, nil,
		},
		{
			`Before $${a}}, after.`, `Before $123}, after.`, nil,
		},
		{
			`Before $${a}}, after.`, `Before $123}, after.`, nil,
		},
		{
			`Before ${c.a[0]}, after.`, `Before Zero, after.`, nil,
		},
		{
			`Before ${a, after.`, "", ErrEndOfInput,
		},
		{
			`Before ${a
}, after.`, `Before 123, after.`, nil,
		},
	}
	for _, e := range tests {
		testInterpolate(t, e.Text, e.Expect, e.Error, context)
	}
}

func testInterpolate(t *testing.T, s string, e string, r error, c interface{}) {
	t.Logf("----> %v\n", s)
	v, err := interpolate(s, "${", "}", c)
	if r != nil {
		t.Logf("<---- (ERR) %v\n", err)
		assert.Equal(t, r, err, s)
	} else {
		t.Logf("<---- %v\n", v)
		assert.Equal(t, e, v)
	}
}
