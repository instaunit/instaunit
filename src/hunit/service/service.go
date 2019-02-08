package service

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// A service
type Service interface {
	Start() error
	Stop() error
}

// Service config
type Config struct {
	Addr     string
	Path     string
	Resource io.ReadCloser
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

	f, err := os.Open(p[1])
	if err != nil {
		return conf, err
	}

	conf.Addr = p[0]
	conf.Path = p[1]
	conf.Resource = f

	return conf, nil
}
