package openapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/instaunit/instaunit/hunit/test"
)

const version = "3.0.3"

type Request struct {
	Req  *http.Request
	Data string
}

type Response struct {
	Rsp  *http.Response
	Data []byte
}

type Specimen struct {
	Suite *test.Suite
	Case  test.Case
	Req   Request
	Rsp   Response
}

func (s Specimen) Id(path string) string {
	if v := s.Case.Route.Id; v != "" {
		return v
	} else if v := s.Case.Route.Path; v != "" {
		return v
	} else {
		return fmt.Sprintf("%s:%s", path, strings.ToLower(s.Req.Req.Method))
	}
}

type Route struct {
	Path  string
	Tests []Specimen
}

type Value struct {
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
	Value       any    `json:"value"`
}

func newValue(ctype string, data []byte) interface{} {
	switch ctype {
	case "application/json":
		return json.RawMessage(data)
	default:
		return Value{Value: string(data)}
	}
}

type Schema struct {
	Type    string      `json:"type,omitempty"`
	Example interface{} `json:"example,omitempty"` // representative object or Value
}

type Payload struct {
	Content map[string]Schema `json:"content,omitempty"`
}

type Reference struct {
	Schema Schema `json:"schema"`
}

type Status struct {
	Summary     string               `json:"summary,omitempty"`
	Description string               `json:"description,omitempty"`
	Status      string               `json:"status"`
	Content     map[string]Reference `json:"content,omitempty"`
}

type Operation struct {
	Id          string                `json:"operationId"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	Tags        []string              `json:"tags,omitempty"`
	Params      []Parameter           `json:"parameters,omitempty"`
	Request     Payload               `json:"requestBody"`
	Responses   map[string]Status     `json:"responses"`
	Security    []SecurityRequirement `json:"security"`
}

type Path struct {
	Path       string
	Operations map[string]Operation // by method
}

func (p Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Operations)
}

type ParameterLocation int

const (
	QueryParameter ParameterLocation = iota
	PathParameter
	parameterLocationCount
)

var parameterLocationNames = []string{
	"query",
	"path",
}

func (p ParameterLocation) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p ParameterLocation) String() string {
	if p < 0 || p > parameterLocationCount {
		return "invalid"
	} else {
		return parameterLocationNames[int(p)]
	}
}

type Parameter struct {
	In          ParameterLocation `json:"in"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Required    bool              `json:"required"`
	Schema      Schema            `json:"schema,omitempty"`
}

type SecurityRequirement map[string][]string

func (a SecurityRequirement) Add(realm string, scopes ...string) SecurityRequirement {
	a[realm] = append(a[realm], scopes...)
	return a
}

type Info struct {
	Title       string `json:"title,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
}

type SecurityScheme struct {
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	In          string `json:"in,omitempty"`
	Scheme      string `json:"scheme,omitempty"`
	Format      string `json:"bearerFormat,omitempty"`
}

type Service struct {
	Standard string           `json:"openapi"`
	Consumes []string         `json:"consumes"`
	Produces []string         `json:"produces"`
	Schemes  []string         `json:"schemes"`
	Security []SecurityScheme `json:"securitySchemes"`
	Info     Info             `json:"info"`
	Host     string           `json:"host"`
	Paths    map[string]Path  `json:"paths"`
}
