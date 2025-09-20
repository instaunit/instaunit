package status

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"sync"

	"github.com/bww/go-util/v1/debug"
)

const addr = ":9222"

type Status string

const (
	Pending Status = "pending"
	Ready   Status = "ready"
	Stopped Status = "stopped"
	Error   Status = "error"
)

type serviceStatus struct {
	Status Status `json:"status"`
}

// status service
type Service struct {
	sync.RWMutex
	server   *http.Server
	services map[string]Status
}

// Create a new status service
func New() (*Service, error) {
	s := &Service{
		services: make(map[string]Status),
	}

	mux := http.NewServeMux()
	mux.Handle("GET /status", http.HandlerFunc(s.handleListStatus))
	mux.Handle("GET /status/{service}", http.HandlerFunc(s.handleFetchStatus))

	svr := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	go func() {
		err := svr.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Status service is shutting down...")
		} else {
			panic(err)
		}
	}()

	s.server = svr

	if debug.VERBOSE {
		fmt.Printf("Status service is running on %s\n", addr)
	}
	return s, nil
}

func (s *Service) Get(svc string) (Status, bool) {
	s.RLock()
	status, ok := s.services[svc]
	s.RUnlock()
	return status, ok
}

func (s *Service) clone() map[string]Status {
	s.RLock()
	svcs := maps.Clone(s.services)
	s.RUnlock()
	return svcs
}

func (s *Service) Set(svc string, status Status) {
	s.Lock()
	s.services[svc] = status
	s.Unlock()
}

func (s *Service) Del(svc string) {
	s.Lock()
	delete(s.services, svc)
	s.Unlock()
}

func (s *Service) handleListStatus(rsp http.ResponseWriter, req *http.Request) {
	svcs := s.clone()

	hdr := rsp.Header()
	hdr.Set("Content-Type", "application/json")

	rsp.WriteHeader(http.StatusOK)
	json.NewEncoder(rsp).Encode(svcs)
}

func (s *Service) handleFetchStatus(rsp http.ResponseWriter, req *http.Request) {
	status, ok := s.Get(req.PathValue("service"))
	if !ok {
		rsp.WriteHeader(http.StatusNotFound)
		return
	}

	hdr := rsp.Header()
	hdr.Set("Content-Type", "application/json")

	rsp.WriteHeader(http.StatusOK)
	json.NewEncoder(rsp).Encode(serviceStatus{
		Status: status,
	})
}
