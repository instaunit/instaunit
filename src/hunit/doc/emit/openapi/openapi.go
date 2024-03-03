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

	var host string
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
			if host == "" {
				host = e.Req.Req.Host
			}
			m := strings.ToLower(e.Req.Req.Method)
			o, ok := p.Operations[m]
			if !ok {
				var id string
				if v := e.Case.Route.Id; v != "" {
					id = v
				} else {
					id = fmt.Sprintf("%s:%s", k, m)
				}
				var reqcnt Payload
				if req := e.Req; len(req.Data) > 0 {
					ctype := text.Coalesce(firstValue(req.Req.Header["Content-Type"]), "text/plain")
					reqcnt = Payload{
						Content: map[string]Schema{
							ctype: Schema{
								Example: newValue(ctype, []byte(req.Data)),
							},
						},
					}
				}
				o = Operation{
					Id:          id,
					Summary:     text.Coalesce(e.Case.Request.Title, e.Case.Title),
					Description: text.Coalesce(e.Case.Request.Comments, e.Case.Comments),
					Request:     reqcnt,
					Responses:   make(map[string]Status),
				}
			}
			var rspcnt map[string]Reference
			if rsp := e.Rsp; len(rsp.Data) > 0 {
				ctype := text.Coalesce(firstValue(rsp.Rsp.Header["Content-Type"]), "text/plain")
				rspcnt = map[string]Reference{
					ctype: Reference{
						Schema: Schema{
							Type:    "object",
							Example: newValue(ctype, []byte(rsp.Data)),
						},
					},
				}
			}
			o.Responses[e.Rsp.Rsp.Status] = Status{
				Summary:     e.Case.Response.Title,
				Description: e.Case.Response.Comments,
				Status:      e.Rsp.Rsp.Status,
				Content:     rspcnt,
			}
			p.Operations[m] = o
		}
		paths[k] = p
	}

	return enc.Encode(Service{
		Standard: version,
		Consumes: []string{"application/json"},
		Produces: []string{"application/json"},
		Schemes:  []string{"https"},
		Info: Info{
			Title:       text.Coalesce(suite.Title, "API"),
			Description: suite.Comments,
		},
		Host:  host,
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

func firstValue[E any](v []E) E {
	var z E
	if len(v) > 0 {
		return v[0]
	} else {
		return z
	}
}
