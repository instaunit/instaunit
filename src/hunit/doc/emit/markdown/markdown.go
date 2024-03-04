package markdown

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/instaunit/instaunit/hunit/test"
	"github.com/instaunit/instaunit/hunit/text"
	"github.com/instaunit/instaunit/hunit/text/slug"
)

const unassignedSection = ""

type entry struct {
	Slug, Title, Section string
}

type section struct {
	Title   string
	Entries []entry
}

// A markdown documentation generator
type Generator struct {
	docpath string
	w       io.WriteCloser
	b       *bytes.Buffer
	entries []entry
	slugs   map[string]int
}

// Produce a new emitter
func New(docpath string) *Generator {
	return &Generator{docpath, nil, nil, make([]entry, 0), nil}
}

// Init a suite; one doc output per suite
func (g *Generator) Init(suite *test.Suite, docs string) error {
	out, err := os.OpenFile(path.Join(g.docpath, docs), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	g.w = out
	g.b = &bytes.Buffer{}
	return nil
}

// Finish a suite
func (g *Generator) Finalize(suite *test.Suite) error {
	defer func() {
		g.w.Close()
		g.w = nil
	}()

	var err error
	err = g.prefix(g.w, suite)
	if err != nil {
		return err
	}

	err = g.contents(g.w, suite)
	if err != nil {
		return err
	}

	_, err = g.w.Write(g.b.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) Close() error {
	return nil
}

// Generate documentation
func (g *Generator) Case(suite *test.Suite, c test.Case, req *http.Request, reqdata string, rsp *http.Response, rspdata []byte) error {
	return g.generate(g.b, suite, c, req, reqdata, rsp, rspdata)
}

// Generate documentation preamble
func (g *Generator) prefix(w io.Writer, suite *test.Suite) error {
	var err error
	var doc string

	if suite.Title != "" {
		doc += fmt.Sprintf("# %s\n\n", strings.TrimSpace(suite.Title))
	}
	if suite.Comments != "" {
		doc += strings.TrimSpace(suite.Comments) + "\n\n"
	}

	_, err = fmt.Fprint(w, doc)
	if err != nil {
		return err
	}

	return nil
}

// Table of contents
func (g *Generator) contents(w io.Writer, suite *test.Suite) error {
	var err error
	var doc string

	doc += "## Contents\n\n"

	if s := suite.TOC.Comments; s != "" {
		doc += s + "\n\n"
	}

	if len(suite.TOC.Sections) > 0 {
		groups := make(map[string][]entry)
		extra := make([]entry, 0)
		for _, e := range g.entries {
			if e.Section == unassignedSection {
				extra = append(extra, e)
			} else {
				s, ok := groups[e.Section]
				if !ok {
					s = make([]entry, 0, 1)
				}
				s = append(s, e)
				groups[e.Section] = s
			}
		}
		for i, e := range suite.TOC.Sections {
			if i > 0 {
				doc += "\n"
			}
			doc += fmt.Sprintf("### %s\n\n", text.Coalesce(e.Title, e.Key))
			for _, e := range groups[e.Key] {
				doc += fmt.Sprintf("* [%s](#%s)\n", strings.TrimSpace(e.Title), e.Slug)
			}
		}
		if !suite.TOC.SuppressUnassigned && len(extra) > 0 {
			if len(suite.TOC.Sections) > 0 {
				doc += "\n### Additional\n\n"
			}
			for _, e := range extra {
				doc += fmt.Sprintf("* [%s](#%s)\n", strings.TrimSpace(e.Title), e.Slug)
			}
		}
	} else {
		for _, e := range g.entries {
			doc += fmt.Sprintf("* [%s](#%s)\n", strings.TrimSpace(e.Title), e.Slug)
		}
	}

	doc += "\n"

	_, err = fmt.Fprint(w, doc)
	if err != nil {
		return err
	}

	return nil
}

// Generate documentation
func (g *Generator) generate(w io.Writer, suite *test.Suite, c test.Case, req *http.Request, reqdata string, rsp *http.Response, rspdata []byte) error {
	var err error
	var doc string

	var t string
	if c.Title != "" {
		t = strings.TrimSpace(c.Title)
	} else {
		t = fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL)
	}

	doc += fmt.Sprintf("## %s\n\n", t)
	var s string
	s, g.slugs = slug.Github(t, g.slugs)
	g.entries = append(g.entries, entry{
		Slug:    s,
		Title:   t,
		Section: c.Section,
	})

	if c.Comments != "" {
		doc += strings.TrimSpace(c.Comments) + "\n\n"
	}

	if c.Params != nil {
		types, maxkey, maxtype := false, 5, 5
		var tmap map[string]string

		params := make(map[string]string)
		keys := make([]string, len(c.Params))
		i := 0
		for k, _ := range c.Params {
			keys[i] = k
			i++
		}
		sort.Strings(keys)

		for k, v := range c.Params {
			t := strings.TrimSpace(k)
			v = strings.TrimSpace(v)
			if l := len(t); l > maxkey {
				maxkey = l
			}
			if v[0] == '`' {
				types = true
				x := strings.Index(v[1:], "`") + 1
				if x > maxtype {
					maxtype = x
				}
				if tmap == nil {
					tmap = make(map[string]string)
				}
				tmap[k] = v[:x+1]
				params[k] = v[x+1:]
			} else {
				params[k] = v
			}
		}

		doc += "### Query Parameters\n\n"

		var f string
		if types {
			doc += "| Param | Type | Detail |\n"
			doc += "|-------|------|--------|\n"
			f = fmt.Sprintf("| %%%ds | %%%ds | %%5s |\n", maxkey, maxtype)
		} else {
			doc += "| Param | Detail |\n"
			doc += "|-------|--------|\n"
			f = fmt.Sprintf("| %%%ds | %%5s |\n", maxkey)
		}

		for _, k := range keys {
			t := strings.TrimSpace(k)
			v := strings.TrimSpace(params[k])
			if types {
				doc += fmt.Sprintf(f, "`"+t+"`", tmap[k], v)
			} else {
				doc += fmt.Sprintf(f, "`"+t+"`", v)
			}
		}

		doc += "\n\n"
	}

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
		if c.Request.Title != "" {
			doc += fmt.Sprintf("### %s\n\n", strings.TrimSpace(c.Request.Title))
		} else if b.Len() > 0 {
			doc += "### Example request\n\n"
		} else if c.Request.Comments != "" {
			doc += "### Request\n\n"
		}
		if c.Request.Comments != "" {
			doc += strings.TrimSpace(c.Request.Comments) + "\n\n"
		}
		if b.Len() > 0 {
			doc += "```http\n"
			doc += string(b.Bytes()) + "\n"
			doc += "```\n\n"
		}
	}

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
		if c.Response.Title != "" {
			doc += fmt.Sprintf("### %s\n\n", strings.TrimSpace(c.Response.Title))
		} else if b.Len() > 0 {
			doc += "### Example response\n\n"
		} else if c.Response.Comments != "" {
			doc += "### Response\n\n"
		}
		if c.Response.Comments != "" {
			doc += strings.TrimSpace(c.Response.Comments) + "\n\n"
		}
		if b.Len() > 0 {
			doc += "```http\n"
			doc += string(b.Bytes()) + "\n"
			doc += "```\n\n"
		}
	}

	_, err = fmt.Fprint(g.b, doc+"\n")
	if err != nil {
		return err
	}

	return nil
}
