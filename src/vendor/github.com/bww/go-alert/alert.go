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

package alt

import (
  "fmt"
  "log"
)

var config Config
var targets []Target
var queue chan *Event

const (
  LEVEL_FATAL Level = iota
  LEVEL_ERROR
  LEVEL_WARNING
  LEVEL_INFO
  LEVEL_DEBUG
)

var levelNames = []string{
  "fatal", "error", "warning", "info", "debug",
}

/**
 * A logging level
 */
type Level int

/**
 * Obtain the level name
 */
func (l Level) Name() string {
  if int(l) >= 0 && int(l) < len(levelNames) {
    return levelNames[int(l)]
  }else{
    return "<unknown>"
  }
}

/**
 * Stringify
 */
func (l Level) String() string {
  return l.Name()
}

const (
  TAG_ENVIRON   = "env"
  TAG_HOSTNAME  = "host"
  TAG_ERROR     = "error"
  TAG_PLATFORM  = "platform"
  TAG_COMPONENT = "component"
  TAG_CONTEXT   = "context"
)

/**
 * A logging target
 */
type Target interface {
  Log(*Event)(error)
}

/**
 * A logging event
 */
type Event struct {
  Level       Level
  Message     string
  Logger      string
  Stacktrace  Stacktrace
  Tags        map[string]interface{}
  Extra       map[string]interface{}
  Display     string
}

/**
 * Create an event
 */
func NewEvent(level Level, m string, tags, extra map[string]interface{}, stack Stacktrace) *Event {
  display := fmt.Sprintf("[%v] %v\n", level, m)
  
  if config.Verbose {
    if tags != nil && len(tags) > 0 {
      var t string
      var i int
      for k, v := range tags {
        if i > 0 {
          t += fmt.Sprintf(", %s = %v", k, v)
        }else{
          t += fmt.Sprintf("%s = %v", k, v)
        }
        i++
      }
      display += fmt.Sprintf("\t# %s\n", t)
    }
  }
  
  if config.Tags != nil && len(config.Tags) > 0 {
    if tags == nil {
      tags = make(map[string]interface{})
    }
    for k, v := range config.Tags {
      tags[k] = v
    }
  }
  
  if config.Verbose {
    if extra != nil {
      var w, l int
      for k, v := range extra {
        if v != nil {
          if l = len(k); l > w { w = l }
        }
      }
      for k, v := range extra {
        if v != nil {
          display += fmt.Sprintf(fmt.Sprintf("\t+ %%%ds: %%v\n", w), k, v)
        }
      }
    }
  }
  
  return &Event{Level:level, Message:m, Logger:config.Name, Tags:tags, Extra:extra, Stacktrace:stack, Display:display}
}

/**
 * Alert configuration
 */
type Config struct {
  Debug       bool
  SentryDSN   string
  Name        string
  Tags        map[string]interface{}  // tags sent with every event
  Backlog     int
  Verbose     bool
  Targets     []Target
}

/**
 * Init
 */
func Init(c Config) {
  if c.Name == "" {
    c.Name = "main"
  }
  
  if c.Backlog > 0 {
    queue = make(chan *Event, c.Backlog)
  }else{
    queue = make(chan *Event, 256)
  }
  
  if c.Targets != nil {
    targets = make([]Target, len(c.Targets))
    for i, e := range c.Targets {
      targets[i] = e
    }
  }
  
  go run(queue)
  
  config = c
}

/**
 * Log for debugging
 */
func Debugf(f string, a ...interface{}) {
  if config.Debug {
    log.Printf(f, a...)
  }
}

/**
 * Log for debugging
 */
func Debug(m string) {
  if config.Debug {
    log.Print(m)
  }
}

/**
 * Log an informative message to sentry
 */
func Infof(f string, a ...interface{}) {
  Info(fmt.Sprintf(f, a...), nil, nil)
}

/**
 * Log an informative message to sentry
 */
func Info(m string, tags, extra map[string]interface{}) {
  Enqueue(NewEvent(LEVEL_INFO, m, tags, extra, generateStacktrace()))
}

/**
 * Log a warning to sentry
 */
func Warnf(f string, a ...interface{}) {
  Warn(fmt.Sprintf(f, a...), nil, nil)
}

/**
 * Log a warning to sentry
 */
func Warn(m string, tags, extra map[string]interface{}) {
  Enqueue(NewEvent(LEVEL_WARNING, m, tags, extra, generateStacktrace()))
}

/**
 * Log an error to sentry
 */
func Errorf(f string, a ...interface{}) {
  Error(fmt.Sprintf(f, a...), nil, nil)
}

/**
 * Log an error to sentry
 */
func Error(m string, tags, extra map[string]interface{}) {
  Enqueue(NewEvent(LEVEL_ERROR, m, tags, extra, generateStacktrace()))
}

/**
 * Log a fatal error to sentry synchronously
 */
func Fatalf(f string, a ...interface{}) {
  Fatal(fmt.Sprintf(f, a...), nil, nil)
}

/**
 * Log a fatal error to sentry synchronously
 */
func Fatal(m string, tags, extra map[string]interface{}) {
  capture(NewEvent(LEVEL_FATAL, m, tags, extra, generateStacktrace()))
}

/**
 * Enqueue an event
 */
func Enqueue(e *Event) {
  log.Print(e.Display)
  if queue != nil {
    queue <- e
  }
}

/**
 * Handle sentry
 */
func run(q <-chan *Event) {
  for e := range q {
    capture(e)
  }
}

/**
 * Capture an event
 */
func capture(e *Event) {
  if targets != nil {
    for _, t := range targets {
      err := t.Log(e)
      if err != nil {
        log.Printf("[alt] Log to target {%T} failed: %v", t, err)
      }
    }
  }
}
