package rest

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/instaunit/instaunit/hunit/expr"
	"github.com/instaunit/instaunit/hunit/expr/runtime"
	"github.com/instaunit/instaunit/hunit/net/await"
	"github.com/instaunit/instaunit/hunit/service"

	"github.com/bww/go-router"
	"github.com/bww/go-util/debug"
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
			return handleRequest((*http.Request)(req), e, vars)
		}
	}

	r := router.New()
	for _, e := range suite.Endpoints {
		if e.Request != nil {
			r.Add(e.Request.Path, handler(e)).Methods(e.Request.Methods...).Params(convertParams(e.Request.Params))
		}
	}

	return &restService{
		conf:   conf,
		suite:  suite,
		router: r,
		vars:   vars,
	}, nil
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
func handleRequest(req *http.Request, endpoint Endpoint, vars expr.Variables) (*router.Response, error) {
	if debug.VERBOSE {
		start := time.Now()
		defer fmt.Printf("%s <- (%v) %s %s\n", prefix, time.Since(start), req.Method, req.URL.Path)
	}

	r := endpoint.Response
	if r == nil {
		return router.NewResponse(http.StatusOK), nil
	}

	e, err := expr.Interpolate(r.Entity, vars)
	if err != nil {
		return nil, err
	}

	x := router.NewResponse(r.Status)
	for k, v := range r.Headers {
		x.SetHeader(k, v)
	}
	if l := len(e); l > 0 {
		x.SetHeader("Content-Length", strconv.FormatInt(int64(l), 10))
		x.SetStringEntity("text/plain", e)
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
