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

package rollbar

import (
  "fmt"
  "time"
  "sync"
  "bytes"
  "net/http"
  "encoding/json"
  "github.com/bww/go-alert"
)

const maxErrors = 5

var client = &http.Client{Timeout:time.Second * 10}

/**
 * A rollbar event
 */
type rollbarEvent struct {
  AccessToken   string              `json:"access_token"`
  Payload       rollbarPayload      `json:"data"`
}

/**
 * A payload
 */
type rollbarPayload struct {
  Environment   string              `json:"environment"`
  Body          rollbarPayloadBody  `json:"body"`
}

/**
 * Another preposterously nested structure with only one element
 */
type rollbarPayloadBody struct {
  Trace         rollbarTrace        `json:"trace"`
  Level         string              `json:"level,omitempty"`
  Timestamp     int64               `json:"timestamp,omitempty"`
  Platform      string              `json:"platform,omitempty"`
  Language      string              `json:"language,omitempty"`
  Framework     string              `json:"framework,omitempty"`
  Context       string              `json:"context,omitempty"`
  Request       *rollbarRequest     `json:"request,omitempty"`
  Server        *rollbarServer      `json:"server,omitempty"`
  Fingerprint   string              `json:"fingerprint,omitempty"`
  Title         string              `json:"title,omitempty"`
  UUID          string              `json:"uuid,omitempty"`
  Notifier      rollbarAgent        `json:"notifier,omitempty"`
}

/**
 * A payload
 */
type rollbarTrace struct {
  Frames        []rollbarFrame      `json:"frames"`
  Exception     rollbarException    `json:"exception"`
}

/**
 * A trace frame
 */
type rollbarFrame struct {
  Filename      string              `json:"filename"`
  LineNumber    int                 `json:"lineno"`
  Function      string              `json:"method"`
}

/**
 * A trace frame
 */
type rollbarException struct {
  Type          string              `json:"class"`
  Message       string              `json:"message,omitempty"`
  Detail        string              `json:"description,omitempty"`
}

/**
 * A request description
 */
type rollbarRequest struct {
  Method        string              `json:"method"`
  URL           string              `json:"url"`
  Headers       map[string]string   `json:"headers"`
  Params        map[string]string   `json:"params"`
  UserIP        string              `json:"user_ip"`
}

/**
 * A request description
 */
type rollbarServer struct {
  Host          string              `json:"host"`
  Root          string              `json:"root"`
}

/**
 * Agent
 */
type rollbarAgent struct {
  Name          string              `json:"name"`
  Version       string              `json:"version"`
}

/**
 * The rollbar logging target
 */
type rollbarTarget struct {
  sync.Mutex
  Token     string
  Threshold alt.Level
  errors    int
}

/**
 * Create a new target
 */
func New(token string, threshold alt.Level) (alt.Target, error) {
  return &rollbarTarget{sync.Mutex{}, token, threshold, 0}, nil
}

/**
 * Log to slack
 */
func (t *rollbarTarget) Log(event *alt.Event) error {
  t.Lock()
  ecap := t.errors - maxErrors
  t.Unlock()
  if ecap > 0 {
    return nil
  }
  if event.Level <= t.Threshold {
    var request *rollbarRequest
    var server *rollbarServer
    
    frames := make([]rollbarFrame, 0)
    if event.Stacktrace.Frames != nil {
      for _, e := range event.Stacktrace.Frames {
        var module string
        if e.Module != "" {
          module = e.Module +"."
        }
        frames = append(frames, rollbarFrame{
          Filename: e.FilePath,
          LineNumber: e.LineNumber,
          Function: module + e.Function,
        })
      }
    }
    
    exception := rollbarException{
      Type: strtag(event.Tags, alt.TAG_ERROR, event.Level.Name()),
      Message: event.Message,
    }
    
    e := rollbarEvent{
      AccessToken: t.Token,
      Payload: rollbarPayload{
        Environment: strtag(event.Tags, alt.TAG_ENVIRON, ""),
        Body: rollbarPayloadBody{
          Trace: rollbarTrace{frames, exception},
          Level: event.Level.Name(),
          Platform: strtag(event.Tags, alt.TAG_PLATFORM, ""),
          Language: "Go",
          Framework: strtag(event.Tags, alt.TAG_COMPONENT, ""),
          Context: event.Logger,
          Request: request,
          Server: server,
          Fingerprint: event.Stacktrace.Fingerprint(),
          Notifier: rollbarAgent{
            Name: "go-alert",
            Version: "1.2",
          },
        },
      },
    }
    
    data, err := json.Marshal(e)
    if err != nil {
      return err
    }
    
    req, err := http.NewRequest("POST", "https://api.rollbar.com/api/1/item/", bytes.NewBuffer(data))
    if err != nil {
      return err
    }
    
    rsp, err := client.Do(req)
    if rsp != nil {
      defer rsp.Body.Close()
    }
    if err != nil {
      return err
    }
    
    if rsp.StatusCode != http.StatusOK {
      t.Lock()
      t.errors++
      t.Unlock()
      return fmt.Errorf("Could not log to Rollbar: %v", rsp.Status)
    }
  }
  return nil
}

/**
 * Obtain a tag or an empty string
 */
func strtag(tags map[string]interface{}, name, dflt string) string {
  if tags == nil {
    return dflt
  }
  
  t, ok := tags[name]
  if !ok {
    return dflt
  }
  
  s, ok := t.(string)
  if !ok {
    return fmt.Sprintf("%v", t)
  }
  
  return s
}
