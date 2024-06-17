package hunit

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/instaunit/instaunit/hunit/runtime"
	"github.com/instaunit/instaunit/hunit/test"

	"github.com/bww/go-util/v1/debug"
	"github.com/bww/go-util/v1/text"
	"github.com/gorilla/websocket"
)

// A future result
type FutureResult interface {
	Finish(time.Time) (*Result, error)
}

// Manages a persistent connection and stream exchange tests
type StreamMonitor struct {
	sync.Mutex
	url      string
	context  runtime.Context
	conn     *websocket.Conn
	messages []test.MessageExchange
	finish   chan struct{}
	valid    bool
	result   *Result
}

// Create a stream monitor for the provided stream
func NewStreamMonitor(url string, context runtime.Context, conn *websocket.Conn, messages []test.MessageExchange) *StreamMonitor {
	return &StreamMonitor{sync.Mutex{}, url, context, conn, messages, nil, false, nil}
}

// Run the stream monitor
func (m *StreamMonitor) Run(result *Result) error {
	m.Lock()
	defer m.Unlock()
	if m.valid || m.finish != nil {
		return fmt.Errorf("Already started")
	}
	if m.conn == nil {
		return fmt.Errorf("Connection is nil")
	}
	m.valid = true
	m.finish = make(chan struct{})
	go m.run(m.conn, m.messages, m.finish, result)
	return nil
}

// Actually run the stream monitor
func (m *StreamMonitor) run(conn *websocket.Conn, messages []test.MessageExchange, finish chan struct{}, result *Result) {
outer:
	for _, e := range messages {
		if e.Wait > 0 {
			<-time.After(e.Wait)
		}

		if e.Output != nil {
			d, err := m.context.Interpolate(*e.Output)
			if err != nil {
				result.Error(err)
				break outer
			}
			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				result.Error(err)
				break outer
			}
			if debug.VERBOSE {
				fmt.Println()
				fmt.Println("---->", m.url)
				fmt.Println(text.Indent(d, "      > "))
			}
			for len(d) > 0 {
				n, err := w.Write([]byte(d))
				if err != nil {
					result.Error(err)
					break outer
				}
				d = d[n:]
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
			x, err := m.context.Interpolate(*e.Input)
			if err != nil {
				result.Error(err)
				break outer
			}
			result.AssertEqual(x, string(d), "Websocket messages do not match")
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
	close(finish)
	m.Unlock()
}

// Finish
func (m *StreamMonitor) Finish(deadline time.Time) (*Result, error) {

	m.Lock()
	v := m.valid
	m.valid = false
	c := m.conn
	m.conn = nil
	x := m.finish
	m.finish = nil
	m.Unlock()
	if v {
		if c == nil {
			return nil, fmt.Errorf("Monitor is valid but connection is nil")
		}
		if x == nil {
			return nil, fmt.Errorf("Monitor is valid but finish channel is nil")
		}
		if !deadline.IsZero() {
			c.SetWriteDeadline(deadline)
			c.SetReadDeadline(deadline)
		}
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
