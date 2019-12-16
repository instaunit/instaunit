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
	"github.com/instaunit/instaunit/hunit/test"
	"github.com/instaunit/instaunit/hunit/text"

	gtext "github.com/bww/go-util/text"
	"github.com/gorilla/websocket"
)

const localVarsId = "vars"

// HTTP client
var client = http.Client{Timeout: time.Second * 30}

// A test context
type Context struct {
	BaseURL   string
	Options   test.Options
	Config    test.Config
	Headers   map[string]string
	Debug     bool
	Gendoc    []doc.Generator
	Variables map[string]interface{}
}

// Derive a subcontext
func (c Context) Subcontext(vars map[string]interface{}) Context {
	return Context{
		BaseURL:   c.BaseURL,
		Options:   c.Options,
		Config:    c.Config,
		Headers:   c.Headers,
		Debug:     c.Debug,
		Gendoc:    c.Gendoc,
		Variables: vars,
	}
}

// Run a test suite
func RunSuite(s *test.Suite, context Context) ([]*Result, error) {
	var futures []FutureResult
	results := make([]*Result, 0)
	c := context.Subcontext(make(map[string]interface{}))

	for _, e := range context.Gendoc {
		err := e.Init(s)
		if err != nil {
			return nil, err
		}
	}

	for _, e := range s.Cases {
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
				r, f, err := RunTest(e, c)
				lock.Lock()
				if err != nil {
					errs = append(errs, err)
				} else {
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
		err := e.Finalize(s)
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
func RunTest(c test.Case, context Context) (*Result, FutureResult, error) {
	start := time.Now()

	// wait if we need to
	if c.Wait > 0 {
		<-time.After(c.Wait)
	}

	// start with an unevaluated result
	result := &Result{Name: fmt.Sprintf("%v %v\n", c.Request.Method, c.Request.URL), Success: true}
	defer func() {
		result.Runtime = time.Since(start)
	}()

	// process variables first, they can be referenced by this case, itself
	var vars map[string]interface{}
	for _, e := range c.Vars {
		k, v := e.Key.(string), gtext.Stringer(e.Value)
		e, err := interpolateIfRequired(context, v)
		if err != nil {
			return result.Error(err), nil, nil
		}
		if x, ok := context.Variables[localVarsId]; ok {
			vars = x.(map[string]interface{})
		} else {
			vars = make(map[string]interface{})
			context.Variables[localVarsId] = vars
		}
		vars[k] = e
	}

	// update the method
	method, err := interpolateIfRequired(context, c.Request.Method)
	if err != nil {
		return result.Error(err), nil, nil
	} else if method == "" {
		return nil, nil, fmt.Errorf("Request requires a method (set 'method')")
	}

	// incrementally update the name as we evaluate it
	result.Name = fmt.Sprintf("%v %v\n", method, c.Request.URL)

	var url string
	if isAbsoluteURL(c.Request.URL) {
		url = c.Request.URL
	} else {
		url = joinPath(context.BaseURL, c.Request.URL)
	}

	if len(c.Request.Params) > 0 {
		url, err = mergeQueryParams(url, c.Request.Params, context)
		if err != nil {
			return result.Error(err), nil, nil
		}
	}

	// incrementally update the name as we evaluate it
	result.Name = fmt.Sprintf("%v %v\n", method, url)

	url, err = interpolateIfRequired(context, url)
	if err != nil {
		return result.Error(err), nil, nil
	} else if url == "" {
		return nil, nil, fmt.Errorf("Request requires a URL (set 'url')")
	}

	// incrementally update the name as we evaluate it
	result.Name = fmt.Sprintf("%v %v\n", method, url)

	header := make(http.Header)
	if context.Headers != nil {
		for k, v := range context.Headers {
			k, err = interpolateIfRequired(context, k)
			if err != nil {
				return result.Error(err), nil, nil
			}
			v, err = interpolateIfRequired(context, v)
			if err != nil {
				return result.Error(err), nil, nil
			}
			header.Add(k, v)
		}
	}

	if c.Request.Headers != nil {
		for k, v := range c.Request.Headers {
			k, err = interpolateIfRequired(context, k)
			if err != nil {
				return result.Error(err), nil, nil
			}
			v, err = interpolateIfRequired(context, v)
			if err != nil {
				return result.Error(err), nil, nil
			}
			header.Add(k, v)
		}
	}

	if c.Request.BasicAuth != nil {
		u, err := interpolateIfRequired(context, c.Request.BasicAuth.Username)
		if err != nil {
			return result.Error(err), nil, nil
		}
		p, err := interpolateIfRequired(context, c.Request.BasicAuth.Password)
		if err != nil {
			return result.Error(err), nil, nil
		}
		header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(u+":"+p))))
	}

	// if we expect a stream, we must set it up and return our future result here
	if c.Stream != nil {
		messages := c.Stream.Messages
		if len(messages) < 1 {
			return result.Error(fmt.Errorf("No messages are exchanged over websocket.")), nil, nil
		}

		dialer := websocket.Dialer{
			NetDial: func(n, a string) (net.Conn, error) {
				return net.DialTimeout(n, a, time.Second*3)
			},
		}
		url, err := urlWithScheme("ws", url)
		if err != nil {
			return result.Error(err), nil, nil
		}
		conn, _, err := dialer.Dial(url, header)
		if err != nil {
			return result.Error(err), nil, nil
		}

		monitor := NewStreamMonitor(url, context, conn, messages)
		err = monitor.Run(result)
		if err != nil {
			return nil, nil, err
		}

		// add to our context if we have an identifier
		if c.Id != "" {
			context.Variables[c.Id] = map[string]interface{}{
				"case": c,
				"vars": dupmap(vars),
				"websocket": map[string]interface{}{
					"url": url,
				},
			}
		}

		// depending on the I/O mode, resolve or return a future
		switch m := c.Stream.Mode; m {
		case test.IOModeSync:
			r, err := monitor.Finish(time.Time{})
			if err != nil {
				return result.Error(err), nil, nil
			}
			return r, nil, nil
		case test.IOModeAsync:
			return nil, monitor, nil
		default:
			return nil, nil, fmt.Errorf("No such I/O mode: %v", m)
		}
	}

	var reqdata string
	var entity io.Reader
	if c.Request.Entity != "" {
		reqdata, err = interpolateIfRequired(context, c.Request.Entity)
		if err != nil {
			return result.Error(err), nil, nil
		} else {
			entity = bytes.NewBuffer([]byte(reqdata))
		}
	}
	if reqdata != "" {
		header.Add("Content-Length", strconv.FormatInt(int64(len(reqdata)), 10))
	}

	req, err := http.NewRequest(method, url, entity)
	if err != nil {
		return nil, nil, err
	}

	req.Header = header

	if c.Request.Cookies != nil {
		for k, v := range c.Request.Cookies {
			k, err = interpolateIfRequired(context, k)
			if err != nil {
				return result.Error(err), nil, nil
			}
			v, err = interpolateIfRequired(context, v)
			if err != nil {
				return result.Error(err), nil, nil
			}
			req.AddCookie(&http.Cookie{Name: k, Value: v})
		}
	}

	reqbuf := &bytes.Buffer{}
	err = text.WriteRequest(reqbuf, req, reqdata)
	if err != nil {
		return result.Error(err), nil, nil
	} else {
		result.Reqdata = reqbuf.Bytes()
	}

	rsp, err := client.Do(req)
	if rsp != nil && rsp.Body != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return result.Error(err), nil, nil
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
		return result.Error(err), nil, nil
	}

	// check response headers, if necessary
	if headers := c.Response.Headers; headers != nil {
		for k, v := range headers {
			k, err = interpolateIfRequired(context, k)
			if err != nil {
				return result.Error(err), nil, nil
			}
			v, err = interpolateIfRequired(context, v)
			if err != nil {
				return result.Error(err), nil, nil
			}
			result.AssertEqual(v, rsp.Header.Get(k), "Headers do not match: %v", k)
		}
	}

	// handle the response entity
	var rspdata []byte
	if rsp.Body != nil {
		rspdata, err = ioutil.ReadAll(rsp.Body)
		if err != nil {
			result.Error(err)
		}
	}

	// parse response entity if necessry
	var rspvalue interface{} = rspdata
	if c.Response.Comparison == test.CompareSemantic {
		rspvalue, err = unmarshalEntity(context, contentType, rspdata)
		if err != nil {
			return result.Error(err), nil, nil
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
			return result.Error(err), nil, nil
		}
		if len(rspdata) == 0 {
			result.AssertEqual(entity, "", "Entities do not match")
		} else if err = entitiesEqual(context, c.Response.Comparison, contentType, []byte(entity), rspvalue); err != nil {
			result.Error(err)
		}
	}

	rspbuf := &bytes.Buffer{}
	err = text.WriteResponse(rspbuf, rsp, rspdata)
	if err != nil {
		return result.Error(err), nil, nil
	} else {
		result.Rspdata = rspbuf.Bytes()
	}

	// add to our context if we have an identifier
	if c.Id != "" {
		headers := flattenHeader(rsp.Header)
		context.Variables[c.Id] = map[string]interface{}{
			"case": c,
			"vars": dupmap(vars),
			"response": map[string]interface{}{
				"headers": headers,
				"entity":  rspdata,
				"value":   rspvalue,
				"status":  rsp.StatusCode,
			},
		}
	}

	// generate documentation if necessary
	if c.Documented() && len(context.Gendoc) > 0 {
		for _, e := range context.Gendoc {
			err := e.Case(context.Config, c, req, reqdata, rsp, rspdata)
			if err != nil {
				return nil, nil, fmt.Errorf("Could not generate documentation: %v", err)
			}
		}
	}

	return result, nil, nil
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

// Duplicate a variable map
func dupmap(m map[string]interface{}) map[string]interface{} {
	d := make(map[string]interface{})
	for k, v := range m {
		d[k] = v
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
