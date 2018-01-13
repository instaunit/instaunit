package hunit

import (
  "fmt"
  "time"
  "sync"
  "io/ioutil"
  
  "hunit/test"
)

import (
  "github.com/bww/go-util/text"
  "github.com/bww/go-util/debug"
  "github.com/gorilla/websocket"
)

// A future result
type FutureResult interface {
  Finish(time.Time)(*Result, error)
}

// Manages a persistent connection and stream exchange tests
type StreamMonitor struct {
  sync.Mutex
  url     string
  context Context
  conn    *websocket.Conn
  stream  test.Stream
  cancel  chan struct{}
  valid   bool
  result  *Result
}

// Create a stream monitor for the provided stream
func NewStreamMonitor(url string, context Context, conn *websocket.Conn, stream test.Stream) *StreamMonitor {
  return &StreamMonitor{sync.Mutex{}, url, context, conn, stream, nil, false, nil}
}

// Run the stream monitor
func (m *StreamMonitor) Run(result *Result) error {
  m.Lock()
  defer m.Unlock()
  if m.valid || m.cancel != nil {
    return fmt.Errorf("Already started")
  }
  if m.conn == nil {
    return fmt.Errorf("Connection is nil")
  }
  m.valid = true
  m.cancel = make(chan struct{})
  go m.run(m.conn, m.stream, m.cancel, result)
  return nil
}

// Actually run the stream monitor
func (m *StreamMonitor) run(conn *websocket.Conn, stream test.Stream, cancel chan struct{}, result *Result) {
  outer:
  for _, e := range stream {
    if e.Output != nil {
      w, err := conn.NextWriter(websocket.TextMessage)
      if err != nil {
        result.Error(err)
        break outer
      }
      b := []byte(*e.Output)
      if debug.VERBOSE {
        fmt.Println()
        fmt.Println("---->", m.url)
        fmt.Println(text.Indent(string(b), "      > "))
      }
      for len(b) > 0 {
        n, err := w.Write(b)
        if err != nil {
          result.Error(err)
          break outer
        }
        b = b[n:]
      }
    }
    
    if e.Input != nil {
      t, r, err := conn.NextReader()
      if err != nil {
        result.Error(err)
        break outer
      }
      if t != websocket.TextMessage {
        result.Error(fmt.Errorf("Unsupported message type: %v", t))
        break outer
      }
      d, err := ioutil.ReadAll(r)
      if err != nil {
        result.Error(err)
        break outer
      }
      if debug.VERBOSE {
        fmt.Println()
        fmt.Println("---->", m.url)
        fmt.Println(text.Indent(string(d), "      < "))
      }
      result.AssertEqual(*e.Input, string(d), "Websocket messages do not match")
    }
    
    m.Lock()
    valid := m.valid
    m.Unlock()
    if !valid {
      break outer
    }
  }
  m.Lock()
  m.result = result
  close(cancel)
  m.Unlock()
}

// Finish
func (m *StreamMonitor) Finish(deadline time.Time) (*Result, error) {
  
  m.Lock()
  v := m.valid
  m.valid = false
  c := m.conn
  m.conn = nil
  x := m.cancel
  m.cancel = nil
  m.Unlock()
  if v {
    if c == nil {
      return nil, fmt.Errorf("Monitor is valid but connection is nil")
    }
    if x == nil {
      return nil, fmt.Errorf("Monitor is valid but canceler nil")
    }
    c.SetWriteDeadline(deadline)
    c.SetReadDeadline(deadline)
    <-x // wait for it to finish...
    c.Close()
  }
  
  m.Lock()
  r := m.result
  m.Unlock()
  if r == nil {
    return nil, fmt.Errorf("No result produced")
  }
  
  return r, nil
}
