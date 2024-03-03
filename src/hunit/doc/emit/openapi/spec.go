package openapi

import (
	"encoding/json"
	"net/http"

	"github.com/instaunit/instaunit/hunit/test"
)

const version = "2.0"

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

type Content struct {
	Example Value `json:"example,omitempty"`
}

type Reference struct {
	Summary     string             `json:"summary,omitempty"`
	Description string             `json:"description,omitempty"`
	Content     map[string]Content `json:"content,omitempty"`
}

type Status struct {
	Reference
	Status string `json:"status"`
}

type Operation struct {
	Id          string            `json:"operationId"`
	Summary     string            `json:"summary,omitempty"`
	Description string            `json:"description,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Request     Content           `json:"requestBody"`
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
	Swagger  string          `json:"swagger"`
	Consumes []string        `json:"consumes"`
	Produces []string        `json:"produces"`
	Schemes  []string        `json:"schemes"`
	Info     Info            `json:"info"`
	Host     string          `json:"host"`
	Paths    map[string]Path `json:"paths"`
}
