package test

import (
	"time"

	"github.com/instaunit/instaunit/hunit/script"
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
	Line, Column int
	Comments     Comments
}

// A test case
type Case struct {
	Id         string                 `yaml:"id"`
	Wait       time.Duration          `yaml:"wait"`
	Repeat     int                    `yaml:"repeat"`
	Concurrent int                    `yaml:"concurrent"`
	Gendoc     bool                   `yaml:"gendoc"`
	Title      string                 `yaml:"title"`
	Section    string                 `yaml:"section"`
	Comments   string                 `yaml:"doc"`
	Require    bool                   `yaml:"require"`
	Params     map[string]string      `yaml:"params"`
	Request    Request                `yaml:"request"`
	Response   Response               `yaml:"response"`
	Stream     *Stream                `yaml:"websocket"`
	Vars       map[string]interface{} `yaml:"vars"`
	Source     Source
}

// Determine if this case is documented or not
func (c Case) Documented() bool {
	return c.Gendoc || c.Title != "" || c.Comments != ""
}
