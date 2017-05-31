package rest

import (
  "io"
  "fmt"
  "time"
  "net/http"
  "hunit/service"
)

// Don't wait forever
const ioTimeout = time.Second * 10

// REST service
type restService struct {
  conf    service.Config
  server  *http.Server
}

// Create a new service
func New(conf service.Config) (service.Service, error) {
  return &restService{conf:conf}, nil
}

// Start the service
func (s *restService) StartService() error {
  if s.server == nil {
    return nil
  }
  
  s.server = &http.Server{
    Addr: conf.Addr,
    Handler: http.HandlerFunc(s.handleRequest),
    ReadTimeout: ioTimeout,
    WriteTimeout: ioTimeout,
    MaxHeaderBytes: 1 << 20,
  }
  
  go func(){
    err := s.server.ListenAndServe()
    if err != nil {
      panic(err)
    }
  }()
  
  return nil
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
func (s *restService) handleRequest(rsp http.ResponseWriter, req *http.Request) error {
  return nil
}
