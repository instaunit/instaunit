package service

import (
  "io"
)

// Service config
type Config struct {
  Addr      string
  Resource  io.Reader
}

// A service
type Service interface {
  StartService()(error)
  StopService()(error)
}
