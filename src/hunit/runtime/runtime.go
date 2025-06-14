package runtime

import (
	"net/http"
	"net/url"

	"github.com/instaunit/instaunit/hunit/doc"
	"github.com/instaunit/instaunit/hunit/expr"
	"github.com/instaunit/instaunit/hunit/protodyn"
	"github.com/instaunit/instaunit/hunit/testcase"
)

// A test context
type Context struct {
	Root      string // test context root (where is the suite run from); this is the directory containing the test suite or '.' if run from STDIN
	BaseURL   *url.URL
	Options   testcase.Options
	Config    testcase.Config
	Headers   map[string]string
	Debug     bool
	Gendoc    []doc.Generator
	Variables expr.Variables
	Service   *protodyn.ServiceRegistry
	Client    *http.Client
}

// Derive a context from the receiver with the provided service registry
func (c Context) WithService(svcreg *protodyn.ServiceRegistry) Context {
	d := c
	d.Service = svcreg
	return d
}

// Derive a context from the receiver with the provided variables
func (c Context) WithVars(vars ...expr.Variables) Context {
	var v map[string]interface{}
	if len(vars) == 1 {
		v = vars[0]
	} else {
		v = mergeVars(vars...)
	}
	d := c
	d.Variables = v
	return d
}

// Merge vars into this context's variables, preferring the parameters
func (c *Context) AddVars(vars ...expr.Variables) {
	c.Variables = mergeVars(append([]expr.Variables{c.Variables}, vars...)...)
}

// ResolveURL resolves a url against this context's base
func (c *Context) ResolveURL(u string) (*url.URL, error) {
	u, err := c.Interpolate(u)
	if err != nil {
		return nil, err
	}
	ref, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	if c.BaseURL == nil {
		return ref, nil
	}
	return c.BaseURL.ResolveReference(ref), nil
}

// Interpolate, if we're configured to do so
func (c *Context) Interpolate(s string) (string, error) {
	if (c.Options & testcase.OptionInterpolateVariables) == testcase.OptionInterpolateVariables {
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
