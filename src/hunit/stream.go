package hunit

import (
  "fmt"
  "time"
  "sync"
  "io/ioutil"
  
  "hunit/test"
)

import (
  "github.com/gorilla/websocket"
)

// A future result
type FutureResult interface {
  Finish(time.Time)(*Result, error)
}

// Manages a persistent connection and stream exchange tests
type StreamMonitor struct {
  sync.Mutex
  conn    *websocket.Conn
  stream  test.Stream
  cancel  chan struct{}
  valid   bool
  result  *Result
}

// Create a stream monitor for the provided stream
func NewStreamMonitor(conn *websocket.Conn, stream test.Stream) *StreamMonitor {
  return &StreamMonitor{sync.Mutex{}, conn, stream, nil, false, nil}
}

// Run the stream monitor
func (m *StreamMonitor) Run(result *Result) error {
  m.Lock()
  defer m.Unlock()
  if m.valid || m.cancel != nil {
    return fmt.Errorf("Already started")
  }
  m.valid = true
  m.cancel = make(chan struct{})
  go m.run(m.conn, m.stream, result)
  return nil
}

// Actually run the stream monitor
func (m *StreamMonitor) run(conn *websocket.Conn, stream test.Stream, result *Result) {
  outer:
  for _, e := range stream {
    if e.Output != nil {
      w, err := conn.NextWriter(websocket.TextMessage)
      if err != nil {
        result.Error(err)
        break outer
      }
      b := []byte(*e.Output)
      fmt.Println(">>> >>> >>>", string(b))
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
      fmt.Println("<<< <<< <<<", string(d))
      // COMPARE!
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
  if m.cancel != nil {
    close(m.cancel)
  }
  m.Unlock()
}

// Finish
func (m *StreamMonitor) Finish(deadline time.Time) (*Result, error) {
  var v bool
  
  m.Lock()
  v = m.valid
  m.Unlock()
  if v {
    return nil, fmt.Errorf("Never started")
  }
  
  m.Lock()
  c := m.conn
  x := m.cancel
  v  = m.valid
  m.valid = false
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
  v  = r == nil
  m.cancel = nil
  m.Unlock()
  if v {
    return nil, fmt.Errorf("No result produced")
  }
  
  return r, nil
}
