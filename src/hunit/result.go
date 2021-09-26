package hunit

import (
	"time"

	"github.com/instaunit/instaunit/hunit/route"
)

// A test result
type Result struct {
	Name    string        `json:"name"`
	Success bool          `json:"success"`
	Skipped bool          `json:"skipped"`
	Route   *route.Route  `json:"route,omitempty"`
	Errors  []string      `json:"errors,omitempty"`
	Reqdata []byte        `json:"request_data,omitempty"`
	Rspdata []byte        `json:"response_data,omitempty"`
	Status  int           `json:"status,omitempty"`
	Context Context       `json:"context"`
	Runtime time.Duration `json:"duration"`
}

// Assert equality. If the values are not equal an error is added to the result.
func (r *Result) AssertEqual(e, a interface{}, m string, x ...interface{}) bool {
	err := assertEqual(e, a, m, x...)
	if err != nil {
		r.Error(err)
		return false
	}
	return true
}

// Add an error to the result. The result is marked as unsuccessful and
// the result is returned so calls can be chained.
func (r *Result) Error(e error) *Result {
	r.Success = false
	r.Errors = append(r.Errors, e.Error())
	return r
}

// Route states
type StatusStats struct {
	Count   int           // request count
	Runtime time.Duration // total request runtime
}

func (s StatusStats) AvgRuntime() time.Duration {
	return s.Runtime / time.Duration(s.Count)
}

// Route states
type RouteStats struct {
	Route    *route.Route
	Requests int                 // request count
	Statuses map[int]StatusStats // result status counts
}

// Result set stats
type Stats struct {
	Routes []RouteStats // distinct routes
}

func NewStats(v []*Result) Stats {
	s := make(map[string]RouteStats)
	var rids []string

	for _, e := range v {
		if r := e.Route; r != nil {
			var t RouteStats
			var ok bool
			if t, ok = s[r.Id]; !ok {
				t.Route = r
				t.Statuses = make(map[int]StatusStats)
				rids = append(rids, r.Id)
			}
			x := t.Statuses[e.Status]
			x.Count++
			x.Runtime += e.Runtime
			t.Requests++
			t.Statuses[e.Status] = x
			s[r.Id] = t
		}
	}

	d := make([]RouteStats, 0, len(s))
	for _, e := range rids {
		d = append(d, s[e])
	}

	return Stats{
		Routes: d,
	}
}
