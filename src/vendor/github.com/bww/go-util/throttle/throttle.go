package util

import (
  "fmt"
  "time"
  "sync"
)

/**
 * Throttle modes
 */
const (
  // allow events to proceed as quickly as possible until limits are reached
  THROTTLE_BURST = iota
  // spread events out evenly over a window
  THROTTLE_METER
)

/**
 * A throttle manages the execution frequency of events. A throttle allows
 * at most N events per time window W, after which it waits until the current
 * window period elapses before allowing further events.
 */
type Throttle struct {
  sync.RWMutex
  events      int
  window      time.Duration
  wstart      time.Time
  wcount      int
  mode        int
}

/**
 * Create a new throttle
 */
func NewThrottle(window time.Duration, events int, mode int) *Throttle {
  return &Throttle{window: window, events: events, mode: mode}
}

/**
 * Obtain the throttle duration
 */
func (t *Throttle) Window() time.Duration {
  return t.window
}

/**
 * Obtain the number of events permitted per window
 */
func (t *Throttle) Events() int {
  return t.events
}

/**
 * Obtain the throttle mode
 */
func (t *Throttle) Mode() int {
  return t.mode
}

/**
 * Enter an event, possibly waiting until the next event can be performed
 */
func (t *Throttle) Enter(capacity int) time.Duration {
  t.Lock()
  defer t.Unlock()
  
  w, d := t.enter(capacity)
  if w != nil {
    <- w
  }
  
  return d
}

/**
 * Enter an event. If a wait is required a non-nil time channel is returned
 * which the caller MUST wait on.
 */
func (t *Throttle) Waiter(capacity int) (<-chan time.Time, time.Duration) {
  t.Lock()
  var n chan time.Time
  
  w, d := t.enter(capacity)
  if w != nil {
    n = make(chan time.Time)
    go func(t *Throttle, w <-chan time.Time, n chan<- time.Time) {
      defer t.Unlock()
      now := <- w
      n <- now
    }(t, w, n)
  }else{
    defer t.Unlock()
  }
  
  return n, d
}

/**
 * Enter an event
 */
func (t *Throttle) enter(capacity int) (<-chan time.Time, time.Duration) {
  
  delay := time.Duration(0)
  now := time.Now()
  
  // copy local
  wstart := t.wstart
  window := t.window
  events := t.events * capacity // events for variable capacity
  
  // if the window has elapsed, clear everything
  if time.Since(wstart) > window {
    t.wstart = now
    t.wcount = 0
  }
  
  // increment the event count
  t.wcount++
  wcount := t.wcount
  
  // if the number of concurrent events exceeds the limit, provide a waitier channel
  if wcount > events {
    if delay = window - time.Since(wstart); delay > 0 {
      return time.After(delay), delay
    }
  }else if t.mode == THROTTLE_METER {
    if delay = window / time.Duration(events); delay > 0 {
      return time.After(delay), delay
    }
  }
  
  return nil, 0
}

/**
 * String description
 */
func (t *Throttle) String() string {
  return fmt.Sprintf("<%v events per %v>", t.events, t.window)
}

/**
 * A throttle that limits the rate of access to elements in a channel.
 */
type ChannelThrottle struct {
  Throttle
  queue <-chan interface{}
}

/**
 * Create a new channel throttle
 */
func NewChannelThrottle(window time.Duration, events int, mode int, queue <-chan interface{}) *ChannelThrottle {
  return &ChannelThrottle{Throttle{window: window, events: events, mode: mode}, queue}
}

/**
 * Create a new channel throttle based on the specified throttle
 */
func NewThrottleWithChannel(t *Throttle, queue <-chan interface{}) *ChannelThrottle {
  return &ChannelThrottle{Throttle{window: t.window, events: t.events, mode: t.mode}, queue}
}

/**
 * Obtain the next element from the channel.
 */
func (t *ChannelThrottle) Next(capacity int) (interface{}, bool) {
  e, ok := <- t.queue
  if !ok {
    return nil, false
  }
  t.Enter(capacity)
  return e, true
}
