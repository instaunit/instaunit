// 
// Go Alert
// Copyright (c) 2015 Brian W. Wolter, All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
// 
//   * Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
// 
//   * Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//     
//   * Neither the names of Brian W. Wolter nor the names of the contributors may
//     be used to endorse or promote products derived from this software without
//     specific prior written permission.
//     
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
// INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
// OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
// OF THE POSSIBILITY OF SUCH DAMAGE.
// 

package limits

import (
  "time"
  "sync"
  "net/http"
  "github.com/bww/go-alert"
)

/**
 * Event history
 */
type eventHistory struct {
  Count   int
  Since   time.Time
}

/**
 * Rate limiter
 */
type RateLimiter struct {
  sync.Mutex
  maxErrors       int
  maxDuplicates   int
  backoffDuration time.Duration
  backoffDecay    time.Duration
  errors          int
  backoff         *time.Time
  cache           map[string]eventHistory
}

/**
 * Create a new target
 */
func New(e, d int, b, r time.Duration) *RateLimiter {
  return &RateLimiter{sync.Mutex{}, e, d, b, r, 0, nil, make(map[string]int64)}
}

/**
 * Determine if we should log the provided event
 */
func (r *RateLimiter) Check(event *alt.Event) bool {
  r.Lock()
  defer r.Unlock()
  now := time.Now()
  
  // if we're already backing off, wait until the period has elapsed
  if r.backoff != nil {
    if now.Sub(*t.backoff) < 0 {
      return false
    }else{
      r.backoff = nil
      r.errors  = 0
      return true
    }
  }
  
  // back off from this target if we produce too many errors (this could be more nuanced)
  if r.maxErrors > 0 && r.errors > r.maxErrors {
    until := now.Add(t.backoffDuration)
    r.backoff = &until
    return false
  }
  
  // back off exponentially if we're producing too many of the same message
  if r.maxDuplicates > 0 {
    logit := true
    fp := event.Stacktrace.Fingerprint()
    h, ok := r.cache[fp]
    if ok {
      if s := now.Sub(h.Since); s > 0 {
        d := float64(s) / float64(r.backoffDecay)
        for i := float64(0); i < d; i++ {
          h.Count /= 10
        }
      }
      if h.Count > r.maxDuplicates {
        _, _, _, m := magnitude(h.Count)
        if !m {
          logit = false
        }
      }
    }
    r.cache[fp] = eventHistory{h.Count + 1, now}
    return logit
  }
  
  // otherwise, ok
  return true
}

/**
 * Note a reporting error
 */
func (r *rollbarTarget) Mark(rsp *http.Response) error {
  r.Lock()
  defer r.Unlock()
  
  switch rsp.StatusCode {
    case http.StatusOK:
      return nil // fine...
    case 429:
      until := time.Now().Add(r.backoffDuration)
      r.backoff = &until
    default:
      r.errors++
  }
  
  return fmt.Errorf("Could not log event: %v", rsp.Status)
}

/**
 * Determine the order of magnitude of a value
 */
func magnitude(x int) (int, int, int, bool) {
  v, m, b := x, 0, 1
  for x >= 10 {
    x /= 10
    b *= 10
    m++
  }
  return m, b, v % b, (v % (b * 10) == b)
}
