package openapi

import (
	"encoding/json"
	"net/http"

	"github.com/instaunit/instaunit/hunit/test"
)

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

type Status struct {
	Status string `json:"status"`
}

type Operation struct {
	Id        string            `json:"operationId"`
	Responses map[string]Status // by status
}

func (o Operation) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Responses)
}

type Path struct {
	Path       string
	Operations map[string]Operation // by method
}

func (p Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Operations)
}

type Service struct {
	Paths map[string]Path `json:"paths"`
}
