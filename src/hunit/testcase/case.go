package testcase

import (
	"time"

	"github.com/instaunit/instaunit/hunit/expr"
	"github.com/instaunit/instaunit/hunit/script"
	yaml "gopkg.in/yaml.v3"
)

// Options
type Options uint64

func (o Options) On(v int) bool {
	return (o & Options(v)) == Options(v)
}

const (
	OptionNone                         = iota
	OptionDebug                        = 1 << iota
	OptionQuiet                        = 1 << iota
	OptionEntityTrimTrailingWhitespace = 1 << iota
	OptionInterpolateVariables         = 1 << iota
	OptionDisplayRequests              = 1 << iota
	OptionDisplayResponses             = 1 << iota
	OptionDisplayRequestsOnFailure     = 1 << iota
	OptionDisplayResponsesOnFailure    = 1 << iota
)

// Basic credentials
type BasicCredentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// A test request
type Request struct {
	Method    string            `yaml:"method"`
	URL       string            `yaml:"url"`
	Headers   map[string]string `yaml:"headers"`
	Cookies   map[string]string `yaml:"cookies"`
	Params    map[string]string `yaml:"params"`
	Entity    string            `yaml:"entity"`
	Format    string            `yaml:"format"`
	BasicAuth *BasicCredentials `yaml:"basic-auth"`
	Title     string            `yaml:"title"`
	Comments  string            `yaml:"doc"`
}

// A test response
type Response struct {
	Status     int               `yaml:"status"`
	Headers    map[string]string `yaml:"headers"`
	Cookies    map[string]string `yaml:"cookies"`
	Entity     string            `yaml:"entity"`
	Comparison Comparison        `yaml:"compare"`
	Format     string            `yaml:"format"`
	Transforms []Transform       `yaml:"transform"`
	Assert     *script.Script    `yaml:"assert"`
	Title      string            `yaml:"title"`
	Comments   string            `yaml:"doc"`
}

// A connection message stream
type Stream struct {
	Mode     IOMode            `yaml:"mode"`
	Messages []MessageExchange `yaml:"messages"`
}

// A message exchange consisting of zero or one input and zero or one output
type MessageExchange struct {
	Wait   time.Duration `yaml:"wait"`
	Output *string       `yaml:"send"`
	Input  *string       `yaml:"receive"`
}

// Source comments
type Comments struct {
	Head, Line, Tail string
}

// Source annotation
type Source struct {
	File         string
	Line, Column int
	Comments     Comments
}

// Route descrition for documentation
type Route struct {
	Id   string `yaml:"id"`   // dynamic
	Path string `yaml:"path"` // dynamic; e.g., '/users/{user_id}'
}

// A frame encapsulates a test case and assocaited local vars
type Frame struct {
	Vars map[string]interface{}
	Case Case
}

// Implemented by types that can produce test frames
type Framer interface {
	Frames() []Frame
}

// A matrix of contexts; each test case is run once per variable context
type Matrix struct {
	Vars  []map[string]interface{} `yaml:"with"`
	Cases []*Case                  `yaml:"do"`
}

func (m *Matrix) Annotate(node *yaml.Node, src Source) error {
	l := len(node.Content)
	for i, e := range node.Content {
		if e.Kind == yaml.ScalarNode && e.Value == "do" {
			if i+1 < l && node.Content[i+1].Kind == yaml.SequenceNode {
				return annotate(m.Cases, src.File, node.Content[i+1]) // once we find it, nothing else to do
			}
		}
	}
	return nil
}

// Produce a frame for every case in the matrix
func (m Matrix) Frames() []Frame {
	var r []Frame
	for _, v := range m.Vars {
		for _, c := range m.Cases {
			r = append(r, Frame{
				Vars: v,
				Case: *c,
			})
		}
	}
	return r
}

// A test case
type Case struct {
	Id         string                   `yaml:"id"`
	Wait       time.Duration            `yaml:"wait"`
	Repeat     int                      `yaml:"repeat"`
	Concurrent int                      `yaml:"concurrent"`
	Gendoc     bool                     `yaml:"gendoc"`
	Route      Route                    `yaml:"route"` // the route description for documentation purposes
	Title      string                   `yaml:"title"`
	Section    string                   `yaml:"section"`
	Comments   string                   `yaml:"doc"`
	Require    bool                     `yaml:"require"`
	Verbose    bool                     `yaml:"verbose"` // enable verbose mode for this test case specifically
	Params     map[string]Parameter     `yaml:"params"`
	Security   map[string]AccessControl `yaml:"security"`
	Request    Request                  `yaml:"request"`
	Response   Response                 `yaml:"response"`
	Stream     *Stream                  `yaml:"websocket"`
	Vars       map[string]interface{}   `yaml:"vars"`
	Source     Source
}

// Determine if this case is documented or not
func (c Case) Documented() bool {
	return c.Gendoc || c.Title != "" || c.Comments != ""
}

func (c *Case) Annotate(node *yaml.Node, src Source) error {
	c.Source = src
	return nil
}

// Produce a single case frame representing this case
func (c Case) Frames() []Frame {
	return []Frame{{
		Vars: make(map[string]interface{}),
		Case: c,
	}}
}

// Interpolate dynamic fields using the provided context and
// return a literalized version of the case.
func (c Case) Interpolate(vars expr.Variables) (Case, error) {
	d := c

	if v, err := expr.Interpolate(c.Title, vars); err != nil {
		return d, err
	} else {
		d.Title = v
	}

	if v, err := expr.Interpolate(c.Comments, vars); err != nil {
		return d, err
	} else {
		d.Comments = v
	}

	if len(c.Security) > 0 {
		sec := make(map[string]AccessControl)
		for k, v := range c.Security {
			k, err := expr.Interpolate(k, vars)
			if err != nil {
				return d, err
			}
			a := AccessControl{}
			for _, e := range v.Scopes {
				e, err := expr.Interpolate(e, vars)
				if err != nil {
					return d, err
				}
				a.Scopes = append(a.Scopes, e)
			}
			for _, e := range v.Roles {
				e, err := expr.Interpolate(e, vars)
				if err != nil {
					return d, err
				}
				a.Roles = append(a.Roles, e)
			}
			sec[k] = a
		}
		d.Security = sec
	}

	return d, nil
}

type caseOrMatrix struct {
	Case   `yaml:",inline"`
	Matrix `yaml:",inline"`
}

func (c *caseOrMatrix) Annotate(node *yaml.Node, src Source) error {
	if c.Matrix.Cases != nil {
		return (&c.Matrix).Annotate(node, src)
	} else {
		return (&c.Case).Annotate(node, src)
	}
}

func (c caseOrMatrix) Frames() []Frame {
	if c.Matrix.Cases != nil {
		return c.Matrix.Frames()
	} else {
		return c.Case.Frames()
	}
}
