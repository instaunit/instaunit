package openapi

import (
	"encoding/json"
	"net/http"

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
	Case test.Case
	Req  Request
	Rsp  Response
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
	Id          string            `json:"operationId"`
	Summary     string            `json:"summary,omitempty"`
	Description string            `json:"description,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Request     Payload           `json:"requestBody"`
	Responses   map[string]Status `json:"responses"`
}

type Path struct {
	Path       string
	Operations map[string]Operation // by method
}

func (p Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Operations)
}

type Info struct {
	Title       string `json:"title,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
}

type Service struct {
	Standard string          `json:"openapi"`
	Consumes []string        `json:"consumes"`
	Produces []string        `json:"produces"`
	Schemes  []string        `json:"schemes"`
	Info     Info            `json:"info"`
	Host     string          `json:"host"`
	Paths    map[string]Path `json:"paths"`
}
