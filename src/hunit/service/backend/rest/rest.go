package rest

import (
  "fmt"
  "net"
  "time"
  "path"
  "strings"
  "strconv"
  "net/http"
  "hunit/service"
)

// Don't wait forever
const ioTimeout = time.Second * 10

// REST service
type restService struct {
  conf    service.Config
  suite   *Suite
  server  *http.Server
}

// Create a new service
func New(conf service.Config) (service.Service, error) {
  suite, err := LoadSuite(conf.Resource)
  if err != nil {
    return nil, err
  }
  return &restService{
    conf:conf,
    suite:suite,
  }, nil
}

// Start the service
func (s *restService) StartService() (int, error) {
  if s.server != nil {
    return 0, fmt.Errorf("Service is running")
  }
  
  s.server = &http.Server{
    Addr: s.conf.Addr,
    Handler: http.HandlerFunc(s.routeRequest),
    ReadTimeout: ioTimeout,
    WriteTimeout: ioTimeout,
    MaxHeaderBytes: 1 << 20,
  }

  listener, err := net.Listen("tcp", s.server.Addr)
  if err != nil {
    return 0, err
  }
  defer listener.Close()

  port := listener.Addr().(*net.TCPAddr).Port 

  go func(){
    panic(s.server.Serve(listener))
  }()
  
  return port, nil
}

// Stop the service
func (s *restService) StopService() error {
  if s.server == nil {
    return fmt.Errorf("Service is not running")
  }
  err := s.server.Close()
  s.server = nil
  return err
}

// Handle requests
func (s *restService) routeRequest(rsp http.ResponseWriter, req *http.Request) {
  
  // match endpoints
  for _, e := range s.suite.Endpoints {
    if r := e.Request; r != nil {
      if r.methods == nil {
        r.methods = make(map[string]struct{})
        for _, x := range r.Methods {
          r.methods[strings.ToLower(x)] = struct{}{}
        }
      }
      if _, ok := r.methods[strings.ToLower(req.Method)]; ok {
        if match, err := path.Match(r.Path, req.URL.Path); err != nil {
          fmt.Printf("* * * Invalid path pattern: %v: %v\n", req.URL, err)
        }else if match {
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
    }else{
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
