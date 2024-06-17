package runtime

import (
	"net/http"

	"github.com/instaunit/instaunit/hunit/doc"
	"github.com/instaunit/instaunit/hunit/expr"
	"github.com/instaunit/instaunit/hunit/test"
)

// A test context
type Context struct {
	BaseURL   string
	Options   test.Options
	Config    test.Config
	Headers   map[string]string
	Debug     bool
	Gendoc    []doc.Generator
	Variables expr.Variables
	Client    *http.Client
}

// Derive a context from the receiver with the provided variables
func (c Context) WithVars(vars ...expr.Variables) Context {
	var v map[string]interface{}
	if len(vars) == 1 {
		v = vars[0]
	} else {
		v = mergeVars(vars...)
	}
	return Context{
		BaseURL:   c.BaseURL,
		Options:   c.Options,
		Config:    c.Config,
		Headers:   c.Headers,
		Debug:     c.Debug,
		Gendoc:    c.Gendoc,
		Client:    c.Client,
		Variables: v,
	}
}

// Merge vars into this context's variables, preferring the parameters
func (c *Context) AddVars(vars ...expr.Variables) {
	c.Variables = mergeVars(append([]expr.Variables{c.Variables}, vars...)...)
}

// Interpolate, if we're configured to do so
func (c *Context) Interpolate(s string) (string, error) {
	if (c.Options & test.OptionInterpolateVariables) == test.OptionInterpolateVariables {
		return expr.Interpolate(s, c.Variables)
	} else {
		return s, nil
	}
}

// Merge maps
func mergeVars(m ...expr.Variables) expr.Variables {
	d := make(expr.Variables)
	for _, e := range m {
		for k, v := range e {
			d[k] = v
		}
	}
	return d
}
