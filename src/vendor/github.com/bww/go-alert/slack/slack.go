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

package slack

import (
  "fmt"
  "time"
  "bytes"
  "net/url"
  "net/http"
  "github.com/bww/go-alert"
)

const maxErrors = 5

var client = &http.Client{Timeout:time.Second * 10}

/**
 * The slack logging target
 */
type slackTarget struct {
  Token     string
  Channel   string
  Prefix    string
  Threshold alt.Level
  errors    int
}

/**
 * Create a new target
 */
func New(token, channel, prefix string, threshold alt.Level) (alt.Target, error) {
  return &slackTarget{token, channel, prefix, threshold, 0}, nil
}

/**
 * Log to slack
 */
func (t *slackTarget) Log(event *alt.Event) error {
  if t.errors > maxErrors {
    return nil // stop trying to log to this target if we produce too many errors
  }
  if event.Level <= t.Threshold {
    input := bytes.NewBuffer([]byte(fmt.Sprintf("*%v*: %v", t.Prefix, event.Message)))
    
    req, err := http.NewRequest("POST", fmt.Sprintf("https://mess.slack.com/services/hooks/slackbot?token=%v&channel=%v", url.QueryEscape(t.Token), url.QueryEscape(fmt.Sprintf("#%v", t.Channel))), input)
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
      t.errors++
      return fmt.Errorf("Could not log to Slack: %v", rsp.Status)
    }
  }
  return nil
}
