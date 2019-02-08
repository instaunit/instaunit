package rest

import (
	"context"
	"fmt"
	"hunit/service"
	"net"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"hunit/net/await"
)

// Don't wait forever
const ioTimeout = time.Second * 10

// Status
const (
	statusMethod = "GET"
	statusPath   = "/_instaunit/status"
)

// REST service
type restService struct {
	conf   service.Config
	suite  *Suite
	server *http.Server
}

// Create a new service
func New(conf service.Config) (service.Service, error) {
	suite, err := LoadSuite(conf.Resource)
	if err != nil {
		return nil, err
	}
	return &restService{
		conf:  conf,
		suite: suite,
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

	// match our internal status endpoint; we don't allow this to be shadowed
	// by defined endpoints in order to monitor the service.
	if req.Method == statusMethod && req.URL.Path == statusPath {
		rsp.Header().Set("Server", "HUnit/1")
		rsp.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rsp.WriteHeader(http.StatusOK)
		return
	}

	// match endpoints
	for _, e := range s.suite.Endpoints {
		if r := e.Request; r != nil {

			if r.methods == nil {
				r.methods = make(map[string]struct{})
				for _, x := range r.Methods {
					r.methods[strings.ToLower(x)] = struct{}{}
				}
			}

			if r.params == nil {
				r.params = make(url.Values)
				u, err := url.Parse(r.Path)
				if err == nil {
					r.params = u.Query()
					r.path = u.Path
				} else {
					r.path = r.Path // just use the full path for matching if the path doesn't parse
				}
			}

			if _, ok := r.methods[strings.ToLower(req.Method)]; ok {
				match, err := path.Match(r.path, req.URL.Path)
				if err != nil {
					fmt.Printf("* * * Invalid path pattern: %v: %v\n", req.URL, err)
				} else if match && paramsMatch(r.params, req.URL.Query()) {
					s.handleRequest(rsp, req, e)
					return
				}
			}

		}
	}

	// nothing matched
	rsp.Header().Set("Server", "HUnit/1")
	rsp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rsp.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(rsp, "Not found.")

}

// Handle requests
func (s *restService) handleRequest(rsp http.ResponseWriter, req *http.Request, endpoint Endpoint) {
	if r := endpoint.Response; r != nil {
		for k, v := range r.Headers {
			rsp.Header().Add(k, v)
		}
		elen := len(r.Entity)
		rsp.Header().Set("Content-Length", strconv.FormatInt(int64(elen), 10))
		if r.Status != 0 {
			rsp.WriteHeader(r.Status)
		} else {
			rsp.WriteHeader(http.StatusOK)
		}
		if elen > 0 {
			_, err := rsp.Write([]byte(r.Entity))
			if err != nil {
				fmt.Printf("* * * Could not write response: %v: %v\n", req.URL, err)
			}
		}
	}
}

// All the parameters in a must be present in b; b may have extra params
func paramsMatch(a, b url.Values) bool {
	for k, v := range a {
		c, ok := b[k]
		if !ok {
			return false
		}
		if !reflect.DeepEqual(v, c) {
			return false
		}
	}
	return true
}
