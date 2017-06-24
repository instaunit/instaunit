package trace

import (
  "io"
  "os"
  "fmt"
  "time"
  "math"
  "strings"
)

import (
  "github.com/bww/go-util/debug"
)

type aggregate func([]time.Duration)(time.Duration)

var displayUnit time.Duration
var groupByName aggregate
var colorize    bool

func init() {
  if os.Getenv("GOUTIL_TRACE_COLORIZE_OUTPUT") != "" && os.Getenv("TERM") != "" {
    colorize = true
  }
  switch strings.ToLower(os.Getenv("GOUTIL_TRACE_GROUP_SPANS_BY")) {
    case "none":        // nothing
    case "avg", "mean": groupByName = mean
    case "max":         groupByName = max
    case "sum":         groupByName = sum
    default:            groupByName = sum
  }
  switch strings.ToLower(os.Getenv("GOUTIL_TRACE_DURATION_UNITS")) {
    case  "s":        displayUnit = time.Second
    case "ms":        displayUnit = time.Millisecond
    case "us", "μs":  displayUnit = time.Microsecond
    default:          displayUnit = time.Nanosecond
  }
}

// An individual span
type Span struct {
  Name      string
  Started   time.Time
  Duration  time.Duration
  Aggregate int
  Spans     []*Span
}

// Begin a sub-span
func (s *Span) Start(n string) *Span {
  var c *Span
  if s != nil {
    c = &Span{n, time.Now(), 0, 0, nil}
    s.Spans = append(s.Spans, c)
  }
  return c
}

// Finish a span
func (s *Span) Finish() {
  if s != nil {
    s.Duration = time.Since(s.Started)
  }
}

// A trace, which manages a set of related spans
type Trace struct {
  Name  string
  Spans []*Span
  warn  time.Duration
}

// Create a trace
func New(n string) *Trace {
  if debug.TRACE {
    return &Trace{Name:n}
  }else{
    return nil
  }
}

// Set the warning threshold
func (t *Trace) Warn(d time.Duration) *Trace {
  if t != nil {
    t.warn = d
  }
  return t
}

// Begin a new span
func (t *Trace) Start(n string) *Span {
  var s *Span
  if t != nil {
    s = &Span{n, time.Now(), 0, 0, nil}
    t.Spans = append(t.Spans, s)
  }
  return s
}

// Finish a trace
func (t *Trace) Finish() {
  if t != nil {
    t.Write(os.Stdout)
  }
}

// Write a trace to the specified writer
func (t *Trace) Write(w io.Writer) (int, error) {
  if t != nil {
    return fmt.Fprint(w, t.format(true, t.Name, t.Spans, 0, 0, "  "))
  }else{
    return 0, nil
  }
}

// Write a trace to the specified writer
func (t *Trace) format(root bool, name string, spans []*Span, depth, rem int, indent string) string {
  var s string
  
  // group by name, use the position of the first occurrance
  if groupByName != nil {
    spans = group(groupByName, spans)
  }
  
  // compute the trace duration
  var et, lt time.Time
  var sd time.Duration
  var si int
  for i, e := range spans {
    if i == 0 || e.Started.Before(et) {
      et = e.Started
    }
    if a := e.Started.Add(e.Duration); a.After(lt) {
      lt = a
    }
    if e.Duration > sd {
      sd = e.Duration
      si = i
    }
  }
  
  if root {
    if td := lt.Sub(et); td > 0 {
      s = fmt.Sprintf("%v (%v in %d spans; longest: #%d @ %s)\n", name, td, len(spans), si + 1, formatDuration(sd))
    }else{
      s = fmt.Sprintf("%v (no closed spans)\n", name)
    }
  }
  
  if l := len(spans); l > 0 {
    nd := int(math.Log10(float64(l + 1))) + 1
    nf := fmt.Sprintf("%%%dd", nd)
    
    var dm int
    var ds []string
    for _, e := range spans {
      var d string
      if e.Duration > 0 {
        d = formatDuration(e.Duration)
      }else {
        d = "(open)"
      }
      ds = append(ds, d)
      if l := len(d); l > dm {
        dm = l
      }
    }
    
    df := fmt.Sprintf("%%%ds", dm)
    for i, e := range spans {
      s += indent
      if depth > 0 {
        last := i + 1 == len(spans)
        for j := 1; j < depth; j++ {
          if last && rem < 1 {
            s += "      "
          }else{
            s += " │    "
          }
        }
        if last {
          s += " └─── "
        }else{
          s += " ├─── "
        }
      }
      warn := t.warn > 0 && e.Duration > t.warn
      if colorize && warn {
        s += string([]byte("\x1b[031m"))
      }
      s += fmt.Sprintf("#"+ nf +" "+ df +" ", i + 1, ds[i])
      s += e.Name
      if e.Aggregate > 1 {
        s += fmt.Sprintf(" (⨉%d)", e.Aggregate)
      }
      if warn {
        s += fmt.Sprintf(" (%v over threshold)", e.Duration - t.warn)
      }
      s += "\n"
      if colorize && warn {
        s += string([]byte("\x1b[000m"))
      }
      if l := len(e.Spans); l > 0 {
        s += t.format(false, e.Name, e.Spans, depth + 1, rem + (len(spans) - i - 1), indent + strings.Repeat(" ", nd - 1))
      }
    }
  }
  
  return s
}

// Group spans using the specified aggregate function
func group(a aggregate, s []*Span) []*Span {
  base := make([]*Span, len(s))
  copy(base, s)
  
  for i := 0; i < len(base); i++ {
    b := base[i]
    m := []time.Duration{b.Duration}
    u := []*Span{}
    for j := i + 1; j < len(base); {
      if c := base[j]; c.Name == b.Name {
        m = append(m, c.Duration)
        u = append(u, c.Spans...)
        for k := j + 1; k < len(base); k++ { base[k-1] = base[k] }
        base = base[:len(base)-1]
      }else{
        j++
      }
    }
    if len(m) > 1 {
      base[i] = &Span{Name:b.Name, Started:b.Started, Duration:a(m), Aggregate:len(m), Spans:u}
    }
  }
  
  return base
}

// Format a duration
func formatDuration(d time.Duration) string {
  if displayUnit == time.Nanosecond {
    return d.String()
  }else{
    return fmt.Sprintf("%f", float64(d) / float64(displayUnit)) + unitSuffix(displayUnit)
  }
}

// Obtain the unit suffix
func unitSuffix(u time.Duration) string {
  switch u {
    case time.Second: return "s"
    case time.Millisecond: return "ms"
    case time.Microsecond: return "μs"
    case time.Nanosecond: return "ns"
    default: return "?s"
  }
}

// Mean aggregate function
func mean(d []time.Duration) time.Duration {
  var t time.Duration
  for _, e := range d {
    t += e
  }
  return t / time.Duration(len(d))
}

// Max aggregate function
func max(d []time.Duration) time.Duration {
  var m time.Duration
  for _, e := range d {
    if e > m {
      m = e
    }
  }
  return m
}

// Sum aggregate function
func sum(d []time.Duration) time.Duration {
  var s time.Duration
  for _, e := range d {
    s += e
  }
  return s
}
