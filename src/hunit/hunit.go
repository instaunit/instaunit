package hunit

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/instaunit/instaunit/hunit/expr"
	"github.com/instaunit/instaunit/hunit/protodyn"
	"github.com/instaunit/instaunit/hunit/runtime"
	"github.com/instaunit/instaunit/hunit/testcase"
)

// Run a test suite
func RunSuite(suite *testcase.Suite, context runtime.Context) ([]*Result, error) {
	var futures []FutureResult
	results := make([]*Result, 0)
	globals := dupVars(suite.Globals)

	// this is weird, but yes, we're evaulating globals in terms of themselves
	globals, err := expr.InterpolateAll(globals, globals)
	if err != nil {
		return nil, fmt.Errorf("Could not evaluate global: %w", err)
	}

	// load protos for the test suite
	svcreg := protodyn.NewServiceRegistry()
	if len(suite.Protos) > 0 {
		for _, e := range suite.Protos {
			err := svcreg.LoadFileDescriptorSetFromPath(path.Join(context.Root, e))
			if err != nil {
				return nil, fmt.Errorf("Could not load Protobuf file descriptor set: %w", err)
			}
		}
		context = context.WithService(svcreg)
	}

	precond := true
	for _, f := range suite.Frames() {
		e := f.Case // just unpack the case for now
		if !precond {
			results = append(results, &Result{Name: fmt.Sprintf("%v %v (dependency failed)\n", e.Request.Method, e.Request.URL), Skipped: true, Case: e})
			continue
		}

		n := e.Repeat
		if n < 1 {
			n = 1
		}

		p := e.Concurrent
		if p < 1 {
			p = 1
		}

		sem := make(chan struct{}, p)
		var wg sync.WaitGroup
		var lock sync.Mutex
		var errs []error

		for i := 0; i < n; i++ {
			lock.Lock()
			nerr := len(errs)
			lock.Unlock()
			if nerr > 0 {
				break
			}

			wg.Add(1)
			sem <- struct{}{}
			go func() {
				defer func() {
					<-sem
					wg.Done()
				}()

				lock.Lock()
				g := dupVars(globals)
				lock.Unlock()

				fvars, err := expr.InterpolateAll(f.Vars, g)
				if err != nil {
					if e.Require {
						precond = false
					}
					errs = append(errs, err)
					return // we're not locked here, so we can return early
				}

				r, f, v, err := RunTest(suite, e, context.WithVars(g, fvars))
				lock.Lock()
				if v != nil && e.Id != "" {
					globals[e.Id] = v
				}
				if err != nil {
					if e.Require {
						precond = false
					}
					errs = append(errs, err)
				} else {
					if e.Require {
						precond = precond && r.Success
					}
					if r != nil {
						results = append(results, r)
					}
					if f != nil {
						futures = append(futures, f)
					}
				}
				lock.Unlock()
			}()
		}

		wg.Wait()

		if len(errs) > 0 {
			return nil, errs[0]
		}
	}

	if len(futures) > 0 {
		d := time.Now() // no grace period by default
		if p := context.Config.Net.StreamIOGracePeriod; p > 0 {
			d = d.Add(p)
		}
		for _, e := range futures {
			r, err := e.Finish(d)
			if err != nil {
				return nil, err
			}
			results = append(results, r)
		}
	}

	return results, nil
}

func formatUnknownName(c testcase.Case) string {
	sb := &strings.Builder{}
	sb.WriteString("<unknown>")
	return formatCaseName(c, sb)
}

func formatRESTName(c testcase.Case, method, url string) string {
	sb := &strings.Builder{}
	sb.WriteString(method)
	sb.WriteString(" ")
	sb.WriteString(url)
	if v := c.Response.Status; v != 0 {
		sb.WriteString(fmt.Sprintf(" (expect: %d/%s)", v, http.StatusText(v)))
	}
	return formatCaseName(c, sb)
}

func formatGRPCName(c testcase.Case) string {
	sb := &strings.Builder{}
	sb.WriteString(c.RPC.Service)
	sb.WriteString(": ")
	sb.WriteString(c.RPC.Method)
	return formatCaseName(c, sb)
}

func formatCaseName(c testcase.Case, sb *strings.Builder) string {
	sb.WriteString(" @ ")
	sb.WriteString(path.Base(c.Source.File))
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(c.Source.Line))
	sb.WriteString("\n")
	return sb.String()
}

// Flatten a header to a one-to-one key-to-value map
func flattenHeader(header http.Header) map[string]string {
	f := make(map[string]string)
	for k, v := range header {
		if len(v) > 0 {
			f[k] = v[0]
		} else {
			f[k] = ""
		}
	}
	return f
}

// Flatten cookies to a one-to-one key-to-value map
func flattenCookies(cookies []*http.Cookie) map[string]string {
	f := make(map[string]string)
	for _, v := range cookies {
		f[v.Name] = v.Value
	}
	return f
}

// Duplicate a variable map
func dupVars(m expr.Variables) expr.Variables {
	d := make(expr.Variables)
	for k, v := range m {
		d[k] = v
	}
	return d
}

// Join paths
func joinPath(a, b string) string {
	if a == "" {
		return b
	}

	var i, j int
	for i = len(a) - 1; i >= 0; i-- {
		if a[i] != '/' {
			break
		}
	}

	for j = 0; j < len(b); j++ {
		if b[j] != '/' {
			break
		}
	}

	return a[:i+1] + "/" + b[j:]
}
