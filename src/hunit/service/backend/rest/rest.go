package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/instaunit/instaunit/hunit/expr"
	"github.com/instaunit/instaunit/hunit/expr/runtime"
	"github.com/instaunit/instaunit/hunit/net/await"
	"github.com/instaunit/instaunit/hunit/service"

	"github.com/bww/go-router/v1"
	"github.com/bww/go-router/v1/entity"
	"github.com/bww/go-util/v1/debug"
	"github.com/bww/go-util/v1/text"

	humanize "github.com/dustin/go-humanize"
)

// Don't wait forever
const ioTimeout = time.Second * 10

// Status
const (
	statusMethod = "GET"
	statusPath   = "/_instaunit/status"
)

const prefix = "[rest]"

// REST service
type restService struct {
	conf   service.Config
	suite  *Suite
	server *http.Server
	router router.Router
	vars   expr.Variables
}

// Create a new service
func New(conf service.Config) (service.Service, error) {
	suite, err := LoadSuite(conf.Resource)
	if err != nil {
		return nil, err
	}

	vars := expr.Variables{
		"std": runtime.Stdlib,
	}

	handler := func(e Endpoint) router.Handler {
		return func(req *router.Request, cxt router.Context) (*router.Response, error) {
			return handleRequest((*http.Request)(req), cxt, e, vars)
		}
	}

	entityHandler := func(endpoints []Endpoint) router.Handler {
		return func(req *router.Request, cxt router.Context) (*router.Response, error) {
			for _, e := range endpoints {
				if routeMatches(e, req) {
					bodyMatch, err := bodyMatches(e.Request.Entity, req)
					if err != nil {
						return nil, err
					}

					if bodyMatch {
						return handleRequest((*http.Request)(req), cxt, e, vars)
					}
				}
			}
			return nil, fmt.Errorf("%s * * * Could not find matching route based on request body: %v: %v\n", prefix, req.URL, err)
		}
	}

	r := router.New()

	for _, e := range suite.Endpoints {
		if e.Request != nil {
			var routeHandler router.Handler

			// if endpoint has request entity defined, use entityHandler
			if e.Request.Entity != "" {
				routeHandler = entityHandler(suite.Endpoints)
			} else {
				routeHandler = handler(e)
			}

			r.Add(e.Request.Path, routeHandler).Methods(e.Request.Methods...).Params(convertParams(e.Request.Params))
		}
	}

	return &restService{
		conf:   conf,
		suite:  suite,
		router: r,
		vars:   vars,
	}, nil
}

// routeMatches is a simple check to see if request method, path, and query parameters match
// what is defined in the endpoint
func routeMatches(e Endpoint, req *router.Request) bool {
	if req.RequestURI != e.Request.Path {
		return false
	}

	if !stringInSlice(req.Method, e.Request.Methods) {
		return false
	}

	var params url.Values
	for k, v := range e.Request.Params {
		params.Add(k, v)
	}

	if params.Encode() != req.URL.Query().Encode() {
		return false
	}

	return true
}

func bodyMatches(entityBody string, req *router.Request) (bool, error) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return false, err
	}

	// replace request body in case it's needed later
	req.Body = ioutil.NopCloser(bytes.NewReader(reqBody))

	// check if request body is not empty, check if matches this endpoint's entity
	if len(reqBody) != 0 {
		var reqData interface{}
		if err := json.Unmarshal(reqBody, &reqData); err != nil {
			return false, err
		}

		var endpointBody interface{}
		if err := json.Unmarshal([]byte(entityBody), &endpointBody); err != nil {
			return false, err
		}

		return reflect.DeepEqual(endpointBody, reqData), nil
	}

	return false, nil
}

func stringInSlice(needle string, haystack []string) bool {
	for _, s := range haystack {
		if strings.ToLower(needle) == strings.ToLower(s) {
			return true
		}
	}
	return false
}

// Start the service
func (s *restService) Start() error {
	if s.server != nil {
		return fmt.Errorf("Service is running")
	}

	host, port, err := net.SplitHostPort(s.conf.Addr)
	if err != nil {
		return fmt.Errorf("Invalid address: %v", err)
	}
	if host == "" {
		host = "localhost"
	}

	s.server = &http.Server{
		Addr:           s.conf.Addr,
		Handler:        http.HandlerFunc(s.routeRequest),
		ReadTimeout:    ioTimeout,
		WriteTimeout:   ioTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// wait for our service to start up
	status := fmt.Sprintf("http://%s:%s%s", host, port, statusPath)
	err = await.Await(context.Background(), []string{status}, ioTimeout)
	if err == await.ErrTimeout {
		return fmt.Errorf("Timed out waiting for service: %s", status)
	} else if err != nil {
		return err
	}

	return nil
}

// Stop the service
func (s *restService) Stop() error {
	if s.server == nil {
		return fmt.Errorf("Service is not running")
	}
	err := s.server.Close()
	s.server = nil
	return err
}

// Handle requests
func (s *restService) routeRequest(rsp http.ResponseWriter, req *http.Request) {
	if debug.VERBOSE {
		var dlen string
		if req.ContentLength < 0 {
			dlen = "unknown length"
		} else {
			dlen = humanize.Bytes(uint64(req.ContentLength))
		}
		fmt.Printf("%s -> %s %s (%s)\n", prefix, req.Method, req.URL.Path, dlen)
		if req.ContentLength > 0 {
			data, err := ioutil.ReadAll(req.Body)
			if err != nil {
				fmt.Printf("%s * * * Could not handle request: %v: %v\n", prefix, req.URL, err)
				return
			}
			req.Body = ioutil.NopCloser(bytes.NewBuffer(data))
			fmt.Println(text.Indent(string(data), strings.Repeat(" ", len(prefix))+" > "))
		}
	}

	// match our internal status endpoint; we don't allow this to be shadowed
	// by defined endpoints so that we can monitor the service.
	if req.Method == statusMethod && req.URL.Path == statusPath {
		rsp.Header().Set("Server", "HUnit/1")
		rsp.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rsp.WriteHeader(http.StatusOK)
		return
	}

	// handle our route
	res, err := s.router.Handle((*router.Request)(req))
	if err != nil {
		fmt.Printf("%s * * * Could not handle request: %v: %v\n", prefix, req.URL, err)
		return
	}

	// write it
	handleResponse(rsp, req, res)
}

// Handle requests
func handleRequest(req *http.Request, cxt router.Context, endpoint Endpoint, vars expr.Variables) (*router.Response, error) {
	var err error

	r := endpoint.Response
	if r == nil {
		return router.NewResponse(http.StatusOK), nil
	}

	var e string
	if debug.VERBOSE {
		start := time.Now()
		defer func() {
			var query string
			if len(req.URL.RawQuery) > 0 {
				query = "?" + req.URL.RawQuery
			}
			fmt.Printf("%s <- %d/%s (%v) %s %s%s (%s)\n", prefix, r.Status, http.StatusText(r.Status), time.Since(start), req.Method, req.URL.Path, query, humanize.Bytes(uint64(len(e))))
			if len(e) > 0 {
				fmt.Println(text.Indent(e, strings.Repeat(" ", len(prefix))+" < "))
			}
		}()
	}

	cvars := make(map[string]interface{})
	for k, v := range cxt.Vars {
		cvars[k] = v
	}

	cparams := make(map[string]interface{})
	for k, v := range req.URL.Query() {
		if len(v) > 0 {
			cparams[k] = v[0]
		}
	}

	err = req.ParseForm()
	if err != nil {
		return nil, err
	}
	cform := make(map[string]interface{})
	for k, v := range req.Form {
		if len(v) > 0 {
			cform[k] = v[0]
		}
	}

	vars["request"] = map[string]interface{}{
		"vars":   cvars,
		"params": cparams,
		"form":   cform,
	}
	e, err = expr.Interpolate(r.Entity, vars)
	if err != nil {
		return nil, err
	}

	x := router.NewResponse(r.Status)
	if l := len(e); l > 0 {
		ent, err := entity.NewString("binary/octet-stream", e)
		if err != nil {
			return nil, err
		}
		_, err = x.SetEntity(ent)
		if err != nil {
			return nil, err
		}
		x.SetHeader("Content-Length", strconv.FormatInt(int64(l), 10))
	}
	for k, v := range r.Headers {
		x.SetHeader(k, v)
	}

	return x, nil
}

// Handle responses
func handleResponse(rsp http.ResponseWriter, req *http.Request, res *router.Response) {
	for k, v := range res.Header {
		rsp.Header().Set(k, v[0])
	}

	if res.Status != 0 {
		rsp.WriteHeader(res.Status)
	} else {
		rsp.WriteHeader(http.StatusOK)
	}

	if e := res.Entity; e != nil {
		defer e.Close()
		_, err := io.Copy(rsp, e)
		if err != nil {
			fmt.Printf("* * * Could not write response: %v: %v\n", req.URL, err)
		}
	}
}

// Create a
func convertParams(p map[string]string) url.Values {
	var r url.Values
	if len(p) > 0 {
		r = make(url.Values)
		for k, v := range p {
			r.Set(k, v)
		}
	}
	return r
}
