package test

import (
	"time"

	"github.com/instaunit/instaunit/hunit/script"

	yaml "gopkg.in/yaml.v2"
)

// Options
type Options uint32

const (
	OptionNone                         = 0
	OptionDebug                        = 1 << 0
	OptionEntityTrimTrailingWhitespace = 1 << 1
	OptionInterpolateVariables         = 1 << 2
	OptionDisplayRequests              = 1 << 3
	OptionDisplayResponses             = 1 << 4
	OptionDisplayRequestsOnFailure     = 1 << 5
	OptionDisplayResponsesOnFailure    = 1 << 6
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

// A test case
type Case struct {
	Id         string            `yaml:"id"`
	Wait       time.Duration     `yaml:"wait"`
	Repeat     int               `yaml:"repeat"`
	Concurrent int               `yaml:"concurrent"`
	Gendoc     bool              `yaml:"gendoc"`
	Title      string            `yaml:"title"`
	Comments   string            `yaml:"doc"`
	Require    bool              `yaml:"require"`
	Params     map[string]string `yaml:"params"`
	Request    Request           `yaml:"request"`
	Response   Response          `yaml:"response"`
	Stream     *Stream           `yaml:"websocket"`
	Vars       yaml.MapSlice     `yaml:"vars"`
}

// Determine if this case is documented or not
func (c Case) Documented() bool {
	return c.Gendoc || c.Title != "" || c.Comments != ""
}
