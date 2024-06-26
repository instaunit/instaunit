package openapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
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
	docpath string
	w       io.WriteCloser
	routes  map[string]*Route
}

// Produce a new emitter
func New(docpath string) *Generator {
	return &Generator{docpath, nil, make(map[string]*Route)}
}

// Init a suite
func (g *Generator) Init(suite *test.Suite, docs string) error {
	if g.w == nil {
		out, err := os.OpenFile(path.Join(g.docpath, "service.json"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		g.w = out
	}
	return nil
}

// Finish a suite
func (g *Generator) Finalize(suite *test.Suite) error {
	return nil // nothing to do
}

// Finalize and close the writer
func (g *Generator) Close() error {
	defer func() {
		g.w.Close()
		g.w = nil
	}()

	enc := json.NewEncoder(g.w)
	enc.SetIndent("", "  ")

	var host string
	paths := make(map[string]Path)
	for k, v := range g.routes {
		if len(v.Tests) < 1 {
			continue // no specimens
		}

		p, ok := paths[k]
		if !ok {
			p = Path{
				Path:       k,
				Operations: make(map[string]Operation),
			}
		}

		// The representative speciment for this collection. We try to use the
		// first successful response encountered for this purpose. If one is not
		// found, the first response is used instead.
		var first, rep *Specimen

		// Process responses
		rsps := make(map[string]Status)
		for i, e := range v.Tests {
			if i == 0 {
				first = &e
			}
			if rep == nil && e.Rsp.Rsp.StatusCode == http.StatusOK {
				rep = &e
			}
			if host == "" {
				host = e.Req.Req.Host
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
			rsps[strconv.Itoa(e.Rsp.Rsp.StatusCode)] = Status{
				Summary:     e.Case.Response.Title,
				Description: e.Case.Response.Comments,
				Status:      e.Rsp.Rsp.Status,
				Content:     rspcnt,
			}
		}
		if rep == nil {
			rep = first
		}

		// Process the request
		m := strings.ToLower(rep.Req.Req.Method)
		var id string
		if v := rep.Case.Route.Id; v != "" {
			id = v
		} else if v := rep.Case.Route.Path; v != "" {
			id = v
		} else {
			id = fmt.Sprintf("%s:%s", k, m)
		}
		var reqcnt Payload
		if req := rep.Req; len(req.Data) > 0 {
			ctype := text.Coalesce(firstValue(req.Req.Header["Content-Type"]), "text/plain")
			reqcnt = Payload{
				Content: map[string]Schema{
					ctype: Schema{
						Example: newValue(ctype, []byte(req.Data)),
					},
				},
			}
		}

		var params []Parameter
		for k, v := range rep.Case.Params {
			params = append(params, Parameter{
				In:          QueryParameter,
				Name:        k,
				Schema:      Schema{Type: v.Type},
				Description: v.Description,
				Required:    v.Required,
			})
		}

		p.Operations[m] = Operation{
			Id:          id,
			Summary:     text.Coalesce(rep.Case.Request.Title, rep.Case.Title),
			Description: text.Coalesce(rep.Case.Request.Comments, rep.Case.Comments),
			Tags:        []string{rep.Suite.Title},
			Params:      params,
			Request:     reqcnt,
			Responses:   rsps,
		}

		paths[k] = p
	}

	return enc.Encode(Service{
		Standard: version,
		Consumes: []string{"application/json"},
		Produces: []string{"application/json"},
		Schemes:  []string{"https"},
		Info: Info{
			Title: "API",
			// Description: suite.Comments,
		},
		Host:  host,
		Paths: paths,
	})
}

// Generate documentation
func (g *Generator) Case(suite *test.Suite, c test.Case, req *http.Request, reqdata string, rsp *http.Response, rspdata []byte) error {
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
		Suite: suite,
		Case:  c,
		Req:   Request{Req: req, Data: reqdata},
		Rsp:   Response{Rsp: rsp, Data: rspdata},
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
