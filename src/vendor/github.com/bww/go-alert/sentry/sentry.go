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

package sentry

import (
  "github.com/bww/go-alert"
  "github.com/bww/raven-go/raven"
)

const maxErrors = 5

/**
 * The sentry logging target
 */
type sentryTarget struct {
  client    *raven.Client
  Threshold alt.Level
  errors    int
}

/**
 * Create a new target
 */
func New(dsn string, threshold alt.Level) (alt.Target, error) {
  client, err := raven.NewClient(dsn)
  if err != nil {
    return nil, err
  }
  return &sentryTarget{client, threshold, 0}, nil
}

/**
 * Log to sentry
 */
func (t *sentryTarget) Log(event *alt.Event) error {
  if t.errors > maxErrors {
    return nil // stop trying to log to this target if we produce too many errors
  }
  if event.Level <= t.Threshold {
    t.client.Capture(&raven.Event{
      Message: event.Message,
      Level: event.Level.Name(),
      Logger: event.Logger,
      Tags: event.Tags,
      Extra: event.Extra,
      Stacktrace:convertStacktrace(event.Stacktrace),
    })
  }
  return nil
}

/**
 * Convert stacktrace
 */
func convertStacktrace(stack alt.Stacktrace) raven.Stacktrace {
  if stack.Frames == nil || len(stack.Frames) < 1 {
    return raven.Stacktrace{Frames:[]raven.Frame{}}
  }
  
  frames := make([]raven.Frame, len(stack.Frames))
  for i, e := range stack.Frames {
    frames[i] = raven.Frame{
      Filename: e.Filename,
      LineNumber: e.LineNumber,
      FilePath: e.FilePath,
      Function: e.Function,
      Module: e.Module,
    }
  }
  
  return raven.Stacktrace{frames}
}
