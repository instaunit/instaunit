package hunit

import (
  "hunit/test"
)

import (
  
)

// Manages a persistent connection and stream exchange tests
type StreamMonitor struct {
  conn    
  stream  test.Stream
  cancel  chan struct{}
}

// Create a stream monitor for the provided stream
func NewStreamMonitor(s test.Stream) *StreamMonitor {
  return &StreamMonitor{s, make(chan struct{})}
}

// Run the stream monitor
func Run() {
}
