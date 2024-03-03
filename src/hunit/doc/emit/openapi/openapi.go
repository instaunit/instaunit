package openapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/instaunit/instaunit/hunit/test"
	"github.com/instaunit/instaunit/hunit/text"
)

const (
	typePlain    = "text/plain"
	typeMarkdown = "text/markdown"
)

// An OpenAPI documentation generator
type Generator struct {
	w      io.WriteCloser
	routes map[string]*Route
}

// Produce a new emitter
func New(w io.WriteCloser) *Generator {
	return &Generator{w, nil}
}

// Init a suite
func (g *Generator) Init(suite *test.Suite) error {
	g.routes = make(map[string]*Route)
	return nil
}

// Finish a suite
func (g *Generator) Finalize(suite *test.Suite) error {
	enc := json.NewEncoder(g.w)
	enc.SetIndent("", "  ")

	paths := make(map[string]Path)
	for k, v := range g.routes {
		p, ok := paths[k]
		if !ok {
			p = Path{
				Path:       k,
				Operations: make(map[string]Operation),
			}
		}
		for _, e := range v.Tests {
			m := strings.ToLower(e.Req.Req.Method)
			o, ok := p.Operations[m]
			if !ok {
				var id string
				if v := e.Case.Route.Id; v != "" {
					id = v
				} else {
					id = fmt.Sprintf("%s:%s", k, m)
				}
				o = Operation{
					Id:          id,
					Summary:     text.Coalesce(e.Case.Request.Title, e.Case.Title),
					Description: text.Coalesce(e.Case.Request.Comments, e.Case.Comments),
					Responses:   make(map[string]Status),
				}
			}
			o.Responses[e.Rsp.Rsp.Status] = Status{
				Status:      e.Rsp.Rsp.Status,
				Summary:     e.Case.Response.Title,
				Description: e.Case.Response.Comments,
			}
			p.Operations[m] = o
		}
		paths[k] = p
	}

	return enc.Encode(Service{
		Paths: paths,
	})
}

// Close the writer
func (g *Generator) Close() error {
	return g.w.Close()
}

// Generate documentation
func (g *Generator) Case(suite *test.Suite, c test.Case, req *http.Request, reqdata string, rsp *http.Response, rspdata []byte) error {
	return g.generate(suite, c, req, reqdata, rsp, rspdata)
}

// Generate documentation
func (g *Generator) generate(suite *test.Suite, c test.Case, req *http.Request, reqdata string, rsp *http.Response, rspdata []byte) error {
	var path string
	if r := c.Route.Path; r != "" {
		path = strings.TrimSpace(r)
	} else if t := c.Title; t != "" {
		path = strings.TrimSpace(t)
	} else {
		path = strings.TrimSpace(c.Request.URL)
	}
	route, ok := g.routes[path]
	if !ok {
		route = &Route{Path: path}
	}
	route.Tests = append(route.Tests, Specimen{
		Case: c,
		Req:  Request{Req: req, Data: reqdata},
		Rsp:  Response{Rsp: rsp, Data: rspdata},
	})
	g.routes[path] = route
	return nil
}
