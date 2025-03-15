package testcase

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/instaunit/instaunit/hunit/exec"

	"github.com/bww/go-util/v1/maps"
	yaml "gopkg.in/yaml.v3"
)

var errMalformedSuite = errors.New("Malformed suite")

type Annotater interface {
	Annotate(node *yaml.Node, src Source) error
}

// Suite options
type Config struct {
	Net struct {
		StreamIOGracePeriod time.Duration `yaml:"stream-io-grace-period"`
	} `yaml:",inline"`
	Doc struct {
		AnchorStyle         AnchorStyle `yaml:"anchor-style"`
		FormatEntities      bool        `yaml:"format-entities"`
		IncludeRequestHTTP  bool        `yaml:"doc-include-request-http"`
		IncludeResponseHTTP bool        `yaml:"doc-include-response-http"`
		IncludeHTTP         bool        `yaml:"doc-include-http"`
	} `yaml:",inline"`
}

// A set of dependencies
type Dependencies struct {
	Resources []string      `yaml:"resources"`
	Timeout   time.Duration `yaml:"timeout"`
}

// A contents section
type Section struct {
	Key   string `yaml:"key"`
	Title string `yaml:"title"`
}

// Table of contents
type TOC struct {
	Sections           []Section `yaml:"sections"`
	Comments           string    `yaml:"doc"`
	SuppressUnassigned bool      `yaml:"suppress-unassigned"`
}

// A test suite
type Suite struct {
	Title     string                    `yaml:"title"`
	Comments  string                    `yaml:"doc"`
	Link      string                    `yaml:"link"` // link to external documentation
	Imports   []string                  `yaml:"import"`
	Authns    map[string]Authentication `yaml:"authentication"`
	TOC       TOC                       `yaml:"toc"`
	Route     Route                     `yaml:"route"` // the route description for documentation purposes; this may be dynamic and shared by all routes in the suite
	Cases     []*caseOrMatrix           `yaml:"tests"`
	Config    Config                    `yaml:"options"`
	Setup     []*exec.Command           `yaml:"setup"`
	Teardown  []*exec.Command           `yaml:"teardown"`
	Transform TransformCollection       `yaml:"transform"`
	Exec      *exec.Command             `yaml:"process"`
	Deps      *Dependencies             `yaml:"depends"`
	Globals   map[string]interface{}    `yaml:"vars"`
}

// Determine if this suite is documented or not
func (s Suite) Documented() bool {
	return s.Title != "" || s.Comments != ""
}

// Produce frames for every case in the suite
func (s Suite) Frames() []Frame {
	var f []Frame
	for _, e := range s.Cases {
		f = append(f, e.Frames()...)
	}
	return f
}

// Load a test suite
func LoadSuiteFromFile(c *Config, p string) (*Suite, error) {
	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return LoadSuiteFromReader(c, p, path.Dir(p), file)
}

// Load a test suite
func LoadSuiteFromReader(c *Config, p, b string, r io.Reader) (*Suite, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return LoadSuiteFromData(c, p, b, data)
}

// Load a test suite
func LoadSuiteFromData(conf *Config, p, b string, data []byte) (*Suite, error) {
	node := &yaml.Node{}
	err := yaml.Unmarshal(data, node)
	if err != nil {
		return nil, err
	}
	if node.Kind != yaml.DocumentNode || len(node.Content) < 1 {
		err = errMalformedSuite
	}

	node = node.Content[0]
	suite := &Suite{Config: *conf}
	switch node.Kind {
	case yaml.SequenceNode:
		err = decodeCases(suite, p, node)
	case yaml.MappingNode:
		err = decodeSuite(suite, p, node)
	default:
		err = errMalformedSuite
	}

	// Imported records are appended in declaration order before
	// those of the current suite...
	var globals map[string]interface{}
	var cases []*caseOrMatrix
	var authns map[string]Authentication
	for _, e := range suite.Imports {
		sub, err := LoadSuiteFromFile(conf, filepath.Join(b, e))
		if err != nil {
			return nil, fmt.Errorf("Could not read imported suite: %w", err)
		}
		// globals
		if len(sub.Globals) > 0 {
			if globals == nil {
				globals = make(map[string]interface{})
			}
			maps.Merge(globals, sub.Globals)
		}
		// test cases
		cases = append(cases, sub.Cases...)
		// authentications
		if len(sub.Authns) > 0 {
			if authns == nil {
				authns = make(map[string]Authentication)
			}
			maps.Merge(authns, sub.Authns)
		}
	}

	// ...globals
	if globals != nil {
		suite.Globals = maps.Merged(globals, suite.Globals)
	}

	// ...current test cases
	cases = append(cases, suite.Cases...)
	suite.Cases = cases

	// ...current authentications
	if authns != nil {
		suite.Authns = maps.Merged(authns, suite.Authns)
	}

	*conf = suite.Config
	return suite, nil
}

// Decode the full suite document format
func decodeSuite(suite *Suite, file string, node *yaml.Node) error {
	err := node.Decode(suite)
	if err != nil {
		return err
	}
	l := len(node.Content)
	for i, e := range node.Content {
		if e.Kind == yaml.ScalarNode && e.Value == "tests" {
			if i+1 < l && node.Content[i+1].Kind == yaml.SequenceNode {
				return annotate(suite.Cases, file, node.Content[i+1])
			}
		}
	}
	return nil
}

// Decode a sequence of cases; this is the simple suite document format
func decodeCases(suite *Suite, file string, node *yaml.Node) error {
	err := node.Decode(&suite.Cases)
	if err != nil {
		return err
	}
	return annotate(suite.Cases, file, node)
}

// Annotate test cases using the provided sequence node
func annotate[E Annotater](cases []E, file string, node *yaml.Node) error {
	if len(cases) != len(node.Content) {
		return errMalformedSuite
	}
	for i, e := range cases {
		n := node.Content[i]
		err := e.Annotate(n, Source{
			File:   file,
			Line:   n.Line,
			Column: n.Column,
			Comments: Comments{
				Head: n.HeadComment,
				Line: n.LineComment,
				Tail: n.FootComment,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Return the first non-nil error or nil if there are none.
func coalesce(err ...error) error {
	for _, e := range err {
		if e != nil {
			return e
		}
	}
	return nil
}
