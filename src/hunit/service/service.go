package service

import (
	"fmt"
	"strings"
	"sync"

	"github.com/instaunit/instaunit/hunit/service/status"
)

var (
	statusSvc  *status.Service
	statusOnce sync.Once
)

func StatusService() (*status.Service, error) {
	var err error
	statusOnce.Do(func() {
		statusSvc, err = status.New()
	})
	return statusSvc, err
}

type Backend string

const (
	REST Backend = "http"
	GRPC Backend = "grpc"
)

// A service
type Service interface {
	Start() error
	Stop() error
}

// Service config
type Config struct {
	Impl   Backend
	Addr   string
	Path   string
	Status *status.Service
}

// Parse configuration
func ParseConfig(s string) (Config, error) {
	var conf Config

	p := strings.Split(s, "=")
	if len(p) != 2 {
		return conf, fmt.Errorf("Invalid service: %v", s)
	}

	if len(p[0]) < 1 {
		return conf, fmt.Errorf("Invalid service address: %v", s)
	}
	if len(p[1]) < 1 {
		return conf, fmt.Errorf("Invalid service resource: %v", s)
	}

	var (
		backend Backend
		addr    string
		sep     = "://"
	)
	if n := strings.Index(p[0], sep); n > 0 {
		backend, addr = Backend(strings.TrimSpace(p[0][:n])), strings.TrimSpace(p[0][n+len(sep):])
	} else {
		backend, addr = REST, p[0]
	}

	conf.Impl = backend
	conf.Addr = addr
	conf.Path = p[1]

	return conf, nil
}
