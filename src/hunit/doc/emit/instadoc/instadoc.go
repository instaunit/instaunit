package instadoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/instaunit/instaunit/hunit/testcase"
	"github.com/instaunit/instaunit/hunit/text"
)

const (
	typePlain    = "text/plain"
	typeMarkdown = "text/markdown"
)

// A markdown documentation generator
type Generator struct {
	docpath string
	w       io.WriteCloser
	toc     *TOC
	routes  []*Route
}

// Produce a new emitter
func New(docpath string) *Generator {
	return &Generator{docpath, nil, nil, nil}
}

// Init a suite; one doc output per suite
func (g *Generator) Init(suite *testcase.Suite, docs string) error {
	out, err := os.OpenFile(path.Join(g.docpath, docs), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	var sects []*Section
	for _, e := range suite.TOC.Sections {
		sects = append(sects, &Section{
			Key:   e.Key,
			Title: e.Title,
		})
	}

	if len(sects) > 0 {
		var detail *Content
		if t := suite.TOC.Comments; t != "" {
			detail = &Content{
				Type: typeMarkdown,
				Data: t,
			}
		}
		g.toc = &TOC{
			Detail:   detail,
			Sections: sects,
		}
	}

	g.routes = make([]*Route, 0)
	g.w = out
	return nil
}

// Finish a suite
func (g *Generator) Finalize(suite *testcase.Suite) error {
	defer func() {
		g.w.Close()
		g.w = nil
	}()

	var detail *Content
	if v := suite.Comments; v != "" {
		detail = &Content{
			Type: typeMarkdown,
			Data: v,
		}
	}

	enc := json.NewEncoder(g.w)
	enc.SetIndent("", "  ")
	return enc.Encode(Suite{
		Title:  suite.Title,
		Detail: detail,
		TOC:    g.toc,
		Routes: g.routes,
	})
}

func (g *Generator) Close() error {
	return nil
}

// Generate documentation
func (g *Generator) Case(suite *testcase.Suite, c testcase.Case, req *http.Request, reqdata string, rsp *http.Response, rspdata []byte) error {
	return g.generate(suite, c, req, reqdata, rsp, rspdata)
}

// Generate documentation
func (g *Generator) generate(suite *testcase.Suite, c testcase.Case, req *http.Request, reqdata string, rsp *http.Response, rspdata []byte) error {
	var err error

	var sections []string
	if c.Section != "" {
		sections = []string{c.Section}
	}

	var title string
	if c.Title != "" {
		title = strings.TrimSpace(c.Title)
	} else {
		title = fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL)
	}

	var comment *Content
	if c.Comments != "" {
		comment = &Content{
			Type: typeMarkdown,
			Data: strings.TrimSpace(c.Comments),
		}
	}

	var params []*Parameter
	if c.Params != nil {
		params = make([]*Parameter, 0, len(c.Params))

		keys := make([]string, 0, len(c.Params))
		for k, _ := range c.Params {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		for _, k := range keys {
			params = append(params, &Parameter{
				Name: strings.TrimSpace(k),
				Type: typePlain,
				Detail: &Content{
					Type: typePlain,
					Data: strings.TrimSpace(c.Params[k].Description),
				},
			})
		}
	}

	var request *Listing
	if req != nil {
		b := &bytes.Buffer{}
		if len(reqdata) > 0 && suite.Config.Doc.FormatEntities {
			t := text.Coalesce(c.Request.Format, req.Header.Get("Content-Type"))
			f, err := text.FormatEntity([]byte(reqdata), t)
			if err == nil {
				reqdata = string(f)
			} else if err != nil && err != text.ErrUnsupportedContentType {
				fmt.Printf("* * * Invalid request entity could not be formatted: %v\n", t)
			}
		}
		err = text.WriteRequest(b, req, reqdata)
		if err != nil {
			return err
		}
		var comment *Content
		if c.Request.Comments != "" {
			comment = &Content{
				Type: typeMarkdown,
				Data: strings.TrimSpace(c.Request.Comments),
			}
		}
		request = &Listing{
			Title:  strings.TrimSpace(c.Request.Title),
			Detail: comment,
			Data:   string(b.Bytes()),
		}
	}

	var response *Listing
	if rsp != nil {
		b := &bytes.Buffer{}
		if len(rspdata) > 0 && suite.Config.Doc.FormatEntities {
			t := text.Coalesce(c.Response.Format, rsp.Header.Get("Content-Type"))
			f, err := text.FormatEntity(rspdata, t)
			if err == nil {
				rspdata = f
			} else if err != nil && err != text.ErrUnsupportedContentType {
				fmt.Printf("* * * Invalid entity could not be formatted: %v\n", t)
			}
		}
		err = text.WriteResponse(b, rsp, rspdata)
		if err != nil {
			return err
		}
		var comment *Content
		if c.Request.Comments != "" {
			comment = &Content{
				Type: typeMarkdown,
				Data: strings.TrimSpace(c.Request.Comments),
			}
		}
		response = &Listing{
			Title:  strings.TrimSpace(c.Request.Title),
			Detail: comment,
			Data:   string(b.Bytes()),
		}
	}

	var examples []*Example
	if request != nil {
		examples = []*Example{
			{
				Request:  request,
				Response: response,
			},
		}
	}

	g.routes = append(g.routes, &Route{
		Sections: sections,
		Title:    title,
		Detail:   comment,
		Method:   req.Method,
		Resource: req.URL.Path,
		Params:   params,
		Examples: examples,
	})
	return nil
}
