package hunit

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/instaunit/instaunit/hunit/doc"
	"github.com/instaunit/instaunit/hunit/expr"
	"github.com/instaunit/instaunit/hunit/test"
	"github.com/instaunit/instaunit/hunit/text"

	"github.com/bww/go-util/v1/debug"
	textutil "github.com/bww/go-util/v1/text"
	"github.com/gorilla/websocket"
)

// A test context
type Context struct {
	BaseURL   string
	Options   test.Options
	Config    test.Config
	Headers   map[string]string
	Debug     bool
	Gendoc    []doc.Generator
	Variables expr.Variables
	Client    *http.Client
}

// Derive a context from the receiver with the provided variables
func (c Context) WithVars(vars ...expr.Variables) Context {
	var v map[string]interface{}
	if len(vars) == 1 {
		v = vars[0]
	} else {
		v = mergeVars(vars...)
	}
	return Context{
		BaseURL:   c.BaseURL,
		Options:   c.Options,
		Config:    c.Config,
		Headers:   c.Headers,
		Debug:     c.Debug,
		Gendoc:    c.Gendoc,
		Client:    c.Client,
		Variables: v,
	}
}

// Merge vars into this context's variables, preferring the parameters
func (c *Context) AddVars(vars ...expr.Variables) {
	c.Variables = mergeVars(append([]expr.Variables{c.Variables}, vars...)...)
}

// Run a test suite
func RunSuite(suite *test.Suite, context Context) ([]*Result, error) {
	var futures []FutureResult
	results := make([]*Result, 0)
	globals := dupVars(suite.Globals)

	for _, e := range context.Gendoc {
		err := e.Init(suite)
		if err != nil {
			return nil, err
		}
	}

	precond := true
	for _, e := range suite.Cases {
		if !precond {
			results = append(results, &Result{Name: fmt.Sprintf("%v %v (dependency failed)\n", e.Request.Method, e.Request.URL), Skipped: true})
			continue
		}

		n := e.Repeat
		if n < 1 {
			n = 1
		}

		p := e.Concurrent
		if p < 1 {
			p = 1
		}

		sem := make(chan struct{}, p)
		var wg sync.WaitGroup
		var lock sync.Mutex
		var errs []error

		for i := 0; i < n; i++ {
			lock.Lock()
			nerr := len(errs)
			lock.Unlock()
			if nerr > 0 {
				break
			}

			wg.Add(1)
			sem <- struct{}{}
			go func() {
				defer func() {
					<-sem
					wg.Done()
				}()
				lock.Lock()
				g := dupVars(globals)
				lock.Unlock()
				r, f, v, err := RunTest(suite, e, context.WithVars(g))
				lock.Lock()
				if v != nil && e.Id != "" {
					globals[e.Id] = v
				}
				if err != nil {
					if e.Require {
						precond = false
					}
					errs = append(errs, err)
				} else {
					if e.Require {
						precond = precond && r.Success
					}
					if r != nil {
						results = append(results, r)
					}
					if f != nil {
						futures = append(futures, f)
					}
				}
				lock.Unlock()
			}()
		}

		wg.Wait()

		if len(errs) > 0 {
			return nil, errs[0]
		}
	}

	for _, e := range context.Gendoc {
		err := e.Finalize(suite)
		if err != nil {
			return nil, err
		}
	}

	if len(futures) > 0 {
		d := time.Now() // no grace period by default
		if p := context.Config.Net.StreamIOGracePeriod; p > 0 {
			d = d.Add(p)
		}
		for _, e := range futures {
			r, err := e.Finish(d)
			if err != nil {
				return nil, err
			}
			results = append(results, r)
		}
	}

	return results, nil
}

// Run a test case
func RunTest(suite *test.Suite, c test.Case, context Context) (*Result, FutureResult, expr.Variables, error) {
	var vdef expr.Variables
	start := time.Now()

	// wait if we need to
	if c.Wait > 0 {
		<-time.After(c.Wait)
	}

	// start with an unevaluated result
	result := &Result{Name: formatName(c, c.Request.Method, c.Request.URL), Success: true}
	defer func() {
		result.Runtime = time.Since(start)
	}()

	// process variables first, they can be referenced by this case, itself
	locals := make(expr.Variables)
	for k, e := range c.Vars {
		v := textutil.Stringer(e)
		r, err := interpolateIfRequired(context, v)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, nil, nil
		}
		locals[k] = r
	}

	// test case variables
	vars := expr.Variables{
		"test": c,
		"vars": locals,
	}
	context.AddVars(vars)

	// update the method
	method, err := interpolateIfRequired(context, c.Request.Method)
	if err != nil {
		return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
	} else if method == "" {
		return nil, nil, nil, fmt.Errorf("Request requires a method (set 'method')")
	}

	// update the url
	url, err := interpolateIfRequired(context, c.Request.URL)
	if err != nil {
		return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
	}
	if !isAbsoluteURL(url) {
		url = joinPath(context.BaseURL, url)
	}

	// incrementally update the name as we evaluate it
	result.Name = formatName(c, method, url)

	if len(c.Request.Params) > 0 {
		url, err = mergeQueryParams(url, c.Request.Params, context)
		if err != nil {
			return nil, nil, vars, fmt.Errorf("Test case declared on line %d: %v", c.Source.Line, err)
		}
	}

	// incrementally update the name as we evaluate it
	result.Name = formatName(c, method, url)

	url, err = interpolateIfRequired(context, url)
	if err != nil {
		return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
	} else if url == "" {
		return nil, nil, nil, fmt.Errorf("Request requires a URL (set 'url')")
	}

	// incrementally update the name as we evaluate it
	result.Name = formatName(c, method, url)

	header := make(http.Header)
	if context.Headers != nil {
		for k, v := range context.Headers {
			k, err = interpolateIfRequired(context, k)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			v, err = interpolateIfRequired(context, v)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			header.Add(k, v)
		}
	}

	if c.Request.Headers != nil {
		for k, v := range c.Request.Headers {
			k, err = interpolateIfRequired(context, k)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			v, err = interpolateIfRequired(context, v)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			header.Add(k, v)
		}
	}

	if c.Request.BasicAuth != nil {
		u, err := interpolateIfRequired(context, c.Request.BasicAuth.Username)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
		}
		p, err := interpolateIfRequired(context, c.Request.BasicAuth.Password)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
		}
		header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(u+":"+p))))
	}

	// if we expect a stream, we must set it up and return our future result here
	if c.Stream != nil {
		messages := c.Stream.Messages
		if len(messages) < 1 {
			return result.Error(fmt.Errorf("No messages are exchanged over websocket.")), nil, vars, nil
		}

		dialer := websocket.Dialer{
			NetDial: func(n, a string) (net.Conn, error) {
				return net.DialTimeout(n, a, time.Second*3)
			},
		}
		url, err := urlWithScheme("ws", url)
		if err != nil {
			return result.Error(fmt.Errorf("Could not upgrade URL scheme: %w", err)), nil, vars, nil
		}
		conn, _, err := dialer.Dial(url, header)
		if err != nil {
			return result.Error(fmt.Errorf("Could not dial websocket: %w", err)), nil, vars, nil
		}

		monitor := NewStreamMonitor(url, context, conn, messages)
		err = monitor.Run(result)
		if err != nil {
			return nil, nil, nil, err
		}

		// websocket variables
		vdef = expr.Variables{
			"url": url,
		}
		vars["websocket"] = vdef
		context.AddVars(expr.Variables{
			"websocket": vdef,
		})

		// depending on the I/O mode, resolve or return a future
		switch m := c.Stream.Mode; m {
		case test.IOModeSync:
			r, err := monitor.Finish(time.Time{})
			if err != nil {
				return result.Error(fmt.Errorf("Could not finish I/O: %w", err)), nil, vars, nil
			}
			return r, nil, vars, nil
		case test.IOModeAsync:
			return nil, monitor, vars, nil
		default:
			return nil, nil, nil, fmt.Errorf("No such I/O mode: %v", m)
		}
	}

	var reqdata string
	var entity io.Reader
	if c.Request.Entity != "" {
		reqdata, err = interpolateIfRequired(context, c.Request.Entity)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
		} else {
			entity = bytes.NewBuffer([]byte(reqdata))
		}
	}
	if reqdata != "" {
		header.Add("Content-Length", strconv.FormatInt(int64(len(reqdata)), 10))
	}

	req, err := http.NewRequest(method, url, entity)
	if err != nil {
		return nil, nil, nil, err
	}

	req.Header = header

	if c.Request.Cookies != nil {
		for k, v := range c.Request.Cookies {
			k, err = interpolateIfRequired(context, k)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			v, err = interpolateIfRequired(context, v)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			req.AddCookie(&http.Cookie{Name: k, Value: v})
		}
	}

	reqbuf := &bytes.Buffer{}
	err = text.WriteRequest(reqbuf, req, reqdata)
	if err != nil {
		return result.Error(fmt.Errorf("Could not write: %w", err)), nil, vars, nil
	} else {
		result.Reqdata = reqbuf.Bytes()
	}

	rsp, err := context.Client.Do(req)
	if rsp != nil && rsp.Body != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return result.Error(fmt.Errorf("Could not read response body: %w", err)), nil, vars, nil
	}

	// check the response status
	result.AssertEqual(c.Response.Status, rsp.StatusCode, "Unexpected status code")

	// note the content type; we prefer the explicit format over the content type
	var contentType string
	if c.Response.Format != "" {
		contentType = c.Response.Format
	} else {
		contentType = strings.ToLower(rsp.Header.Get("Content-Type"))
	}

	contentType, err = interpolateIfRequired(context, contentType)
	if err != nil {
		return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
	}

	// check response headers, if necessary
	if headers := c.Response.Headers; headers != nil {
		for k, v := range headers {
			k, err = interpolateIfRequired(context, k)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			v, err = interpolateIfRequired(context, v)
			if err != nil {
				return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
			}
			result.AssertEqual(v, rsp.Header.Get(k), "Headers do not match: %v", k)
		}
	}

	// handle the response entity
	var rspdata []byte
	if rsp.Body != nil {
		rspdata, err = ioutil.ReadAll(rsp.Body)
		if err != nil {
			result.Error(fmt.Errorf("Could not read response body: %w", err))
		}
	}

	// parse response entity if necessry
	var rspvalue interface{} = rspdata
	if c.Response.Comparison == test.CompareSemantic {
		rspvalue, err = unmarshalEntity(context, contentType, rspdata)
		if err != nil {
			return result.Error(fmt.Errorf("Could not unmarshal entity: %w", err)), nil, vars, nil
		}
	} else if c.Id != "" { // attempt it but don't produce an error if we fail
		val, err := unmarshalEntity(context, contentType, rspdata)
		if err == nil {
			rspvalue = val
		}
	}

	// check response entity, if necessary
	if entity := c.Response.Entity; entity != "" {
		entity, err = interpolateIfRequired(context, entity)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
		}
		if len(rspdata) == 0 {
			result.AssertEqual(entity, "", "Entities do not match")
		} else if err = entitiesEqual(context, c.Response.Comparison, contentType, []byte(entity), rspvalue); err != nil {
			result.Error(fmt.Errorf("Could not compare entities: %w", err))
		}
	}

	rspbuf := &bytes.Buffer{}
	err = text.WriteResponse(rspbuf, rsp, rspdata)
	if err != nil {
		return result.Error(fmt.Errorf("Could not write: %w", err)), nil, vars, nil
	} else {
		result.Rspdata = rspbuf.Bytes()
	}

	// response variables
	vdef = expr.Variables{
		"headers": flattenHeader(rsp.Header),
		"cookies": flattenCookies(rsp.Cookies()),
		"entity":  rspdata,
		"value":   rspvalue,
		"status":  rsp.StatusCode,
	}
	vars["response"] = vdef
	context.AddVars(expr.Variables{
		"response": vdef,
	})

	// update request with final context
	result.Context = context

	// update test case dynamic post-fields with response
	if r := c.Route.Id; r != "" {
		r, err = interpolateIfRequired(context, r)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
		}
		c.Route.Id = r
	}
	if r := c.Route.Path; r != "" {
		r, err = interpolateIfRequired(context, r)
		if err != nil {
			return result.Error(fmt.Errorf("Could not interpolate: %w", err)), nil, vars, nil
		}
		c.Route.Path = r
	}

	// assertions
	if assert := c.Response.Assert; assert != nil {
		ok, err := assert.Bool(context.Variables)
		if err != nil {
			b := &strings.Builder{}
			debug.Dumpf(b, context.Variables)
			return nil, nil, nil, fmt.Errorf("Could not evaluate assertion: %v\n%s", err, b.String())
		}
		if !ok {
			result.Error(&ScriptError{"Script assertion failed", true, ok, assert})
		}
	}

	// generate documentation if necessary
	if c.Documented() && len(context.Gendoc) > 0 {
		for _, e := range context.Gendoc {
			err := e.Case(suite, c, req, reqdata, rsp, rspdata)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("Could not generate documentation: %v", err)
			}
		}
	}

	return result, nil, vars, nil
}

func formatName(c test.Case, method, url string) string {
	return fmt.Sprintf("%v %v @ line %d\n", method, url, c.Source.Line)
}

// Flatten a header to a one-to-one key-to-value map
func flattenHeader(header http.Header) map[string]string {
	f := make(map[string]string)
	for k, v := range header {
		if len(v) > 0 {
			f[k] = v[0]
		} else {
			f[k] = ""
		}
	}
	return f
}

// Flatten cookies to a one-to-one key-to-value map
func flattenCookies(cookies []*http.Cookie) map[string]string {
	f := make(map[string]string)
	for _, v := range cookies {
		f[v.Name] = v.Value
	}
	return f
}

// Duplicate a variable map
func dupVars(m expr.Variables) expr.Variables {
	d := make(expr.Variables)
	for k, v := range m {
		d[k] = v
	}
	return d
}

// Merge maps
func mergeVars(m ...expr.Variables) expr.Variables {
	d := make(expr.Variables)
	for _, e := range m {
		for k, v := range e {
			d[k] = v
		}
	}
	return d
}

// Join paths
func joinPath(a, b string) string {
	if a == "" {
		return b
	}

	var i, j int
	for i = len(a) - 1; i >= 0; i-- {
		if a[i] != '/' {
			break
		}
	}

	for j = 0; j < len(b); j++ {
		if b[j] != '/' {
			break
		}
	}

	return a[:i+1] + "/" + b[j:]
}
