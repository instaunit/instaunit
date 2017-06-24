package trace

import (
  "time"
  "testing"
  // "github.com/stretchr/testify/assert"
)

import (
  "github.com/bww/go-util/debug"
)

/**
 * Test no-op trace handling, when debug.TRACE is false
 */
func TestNopTrace(t *testing.T) {
  r := New("Hello").Warn(time.Millisecond)
  r.Start("Sub-operation").Finish()
  r.Start("Another operation").Finish()
  r.Finish()
}

/**
 * Test trace
 */
func TestTrace(t *testing.T) {
  debug.TRACE = true
  r := New("Hello").Warn(time.Millisecond)
  r.Start("Sub-operation").Finish()
  r.Start("Another operation").Finish()
  r.Start("Another operation").Finish()
  r.Start("Another operation").Finish()
  r.Start("Open operation")
  s := r.Start("Enjoy this one as well")
  u := s.Start("Sub-op")
  <- time.After(time.Millisecond * 1)
  u.Finish()
  d := u.Start("Sub-sub-op!")
  <- time.After(time.Millisecond * 1)
  d.Finish()
  u = s.Start("Sub-op again")
  <- time.After(time.Millisecond * 1)
  u.Finish()
  d = u.Start("Sub-sub-op again!")
  <- time.After(time.Millisecond * 1)
  d.Finish()
  u.Start("Sub-sub-op again!").Finish()
  u.Start("Sub-sub-op again!").Finish()
  s.Finish()
  r.Finish()
  // assert.Equal(t, true, ResemblesUUID("ACE24573-5BD5-4C5F-B143-5E9E17F18BDB"))
}
