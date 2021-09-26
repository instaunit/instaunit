package markdown

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/instaunit/instaunit/hunit/doc/emit"
	"github.com/instaunit/instaunit/hunit/test"
	"github.com/instaunit/instaunit/hunit/text"
	"github.com/instaunit/instaunit/hunit/text/slug"
)

// A markdown documentation generator
type Generator struct {
	w        io.WriteCloser
	b        *bytes.Buffer
	sections []string
	sectmap  map[string]string
	slugs    map[string]int
}

// Produce a new emitter
func New(w io.WriteCloser) *Generator {
	return &Generator{
		w:       w,
		sectmap: make(map[string]string),
	}
}

// Init a suite
func (g *Generator) Init(suite *test.Suite) error {
	g.b = &bytes.Buffer{}
	return nil
}

// Finish a suite
func (g *Generator) Finalize(suite *test.Suite) error {
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

// Close the wroter
func (g *Generator) Close() error {
	return g.w.Close()
}

// Generate documentation
func (g *Generator) Case(conf test.Config, c emit.Case) error {
	return g.generate(g.b, conf, c)
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

	for _, e := range g.sections {
		s, t := e, g.sectmap[e]
		doc += fmt.Sprintf("* [%s](#%s)\n", strings.TrimSpace(t), s)
	}

	doc += "\n"

	_, err = fmt.Fprint(w, doc)
	if err != nil {
		return err
	}

	return nil
}

// Generate documentation
func (g *Generator) generate(w io.Writer, conf test.Config, c emit.Case) error {
	var err error
	var doc string

	var t string
	if c.Case.Title != "" {
		t = strings.TrimSpace(c.Case.Title)
	} else if c.Route != nil && c.Route.Name != "" {
		t = c.Route.Name
	} else {
		t = fmt.Sprintf("%s %s", c.Case.Request.Method, c.Case.Request.URL)
	}

	doc += fmt.Sprintf("## %s\n\n", t)
	var s string
	s, g.slugs = slug.Github(t, g.slugs)
	if _, ok := g.sectmap[s]; !ok {
		g.sections = append(g.sections, s)
		g.sectmap[s] = t
	}

	if c.Case.Comments != "" {
		doc += strings.TrimSpace(c.Case.Comments) + "\n\n"
	}

	if c.Case.Params != nil {
		types, maxkey, maxtype := false, 5, 5
		var tmap map[string]string

		params := make(map[string]string)
		keys := make([]string, len(c.Case.Params))
		i := 0
		for k, _ := range c.Case.Params {
			keys[i] = k
			i++
		}
		sort.Strings(keys)

		for k, v := range c.Case.Params {
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

	if c.Request != nil {
		b := &bytes.Buffer{}
		if len(c.Reqdata) > 0 && conf.Doc.FormatEntities {
			t := text.Coalesce(c.Case.Request.Format, c.Request.Header.Get("Content-Type"))
			f, err := text.FormatEntity([]byte(c.Reqdata), t)
			if err == nil {
				c.Reqdata = f
			} else if err != nil && err != text.ErrUnsupportedContentType {
				fmt.Printf("* * * Invalid request entity could not be formatted: %v\n", t)
			}
		}
		err = text.WriteRequest(b, c.Request, string(c.Reqdata))
		if err != nil {
			return err
		}
		if c.Case.Request.Title != "" {
			doc += fmt.Sprintf("### %s\n\n", strings.TrimSpace(c.Case.Request.Title))
		} else if b.Len() > 0 {
			doc += "### Example request\n\n"
		} else if c.Case.Request.Comments != "" {
			doc += "### Request\n\n"
		}
		if c.Case.Request.Comments != "" {
			doc += strings.TrimSpace(c.Case.Request.Comments) + "\n\n"
		}
		if b.Len() > 0 {
			doc += "```http\n"
			doc += string(b.Bytes()) + "\n"
			doc += "```\n\n"
		}
	}

	if c.Response != nil {
		b := &bytes.Buffer{}
		if len(c.Rspdata) > 0 && conf.Doc.FormatEntities {
			t := text.Coalesce(c.Case.Response.Format, c.Response.Header.Get("Content-Type"))
			f, err := text.FormatEntity(c.Rspdata, t)
			if err == nil {
				c.Rspdata = f
			} else if err != nil && err != text.ErrUnsupportedContentType {
				fmt.Printf("* * * Invalid entity could not be formatted: %v\n", t)
			}
		}
		err = text.WriteResponse(b, c.Response, c.Rspdata)
		if err != nil {
			return err
		}
		if c.Case.Response.Title != "" {
			doc += fmt.Sprintf("### %s\n\n", strings.TrimSpace(c.Case.Response.Title))
		} else if b.Len() > 0 {
			doc += "### Example response\n\n"
		} else if c.Case.Response.Comments != "" {
			doc += "### Response\n\n"
		}
		if c.Case.Response.Comments != "" {
			doc += strings.TrimSpace(c.Case.Response.Comments) + "\n\n"
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
