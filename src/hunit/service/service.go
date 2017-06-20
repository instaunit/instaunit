package service

import (
  "os"
  "io"
  "fmt"
  "strings"
)

// A service
type Service interface {
  StartService()(int, error)
  StopService()(error)
}

// Service config
type Config struct {
  Addr      string
  Path      string
  Resource  io.ReadCloser
}

// Parse configuration
func ParseConfig(s string) (Config, error) {
  var conf Config
  var srcPath, port string 
  
  p := strings.Split(s, "=")
  if len(p) == 2 {
    if len(p[0]) < 1 {
      return conf, fmt.Errorf("Invalid service address: %v", s)
    }
    if len(p[1]) < 1 {
      return conf, fmt.Errorf("Invalid service resource: %v", s)
    }
    port = p[0]
    srcPath = p[1]
  } else if len(p) == 1 {
    // bind the service to port 0, a random free port from 1024 to 65535 will be selected
    port = ":0"
    srcPath = p[0]
  } else {
    return conf, fmt.Errorf("Invalid service: %v", s)
  }

  f, err := os.Open(srcPath)
  if err != nil {
    return conf, err
  }

  conf.Addr = port
  conf.Path = srcPath
  conf.Resource = f
  
  return conf, nil
}
