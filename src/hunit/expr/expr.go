package expr

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/instaunit/instaunit/hunit/env"
	"github.com/instaunit/instaunit/hunit/expr/runtime"

	"github.com/bww/epl/v1"
)

var ErrEndOfInput = errors.New("Unexpected end-of-input")

type ExprError struct {
	err error
	cxt interface{}
}

func newExprError(err error, cxt interface{}) ExprError {
	return ExprError{
		err: err,
		cxt: cxt,
	}
}

func (e ExprError) Unwrap() error {
	return e.err
}

func (e ExprError) Error() string {
	if env.ExprDebug {
		return e.Detail()
	} else {
		return e.Summary()
	}
}

func (e ExprError) Summary() string {
	return e.err.Error()
}

func (e ExprError) Detail() string {
	return fmt.Sprintf("%v\nContext:\n%s", e.err, spew.Sdump(e.cxt))
}

// Variables
type Variables map[string]interface{}

// Map the environment
func mapenv(v []string) map[string]string {
	env := make(map[string]string)
	for _, e := range v {
		if x := strings.Index(e, "="); x > 0 {
			env[e[:x]] = e[x+1:]
		} else {
			env[e] = ""
		}
	}
	return env
}

// Produce a context with the standard library included
func RuntimeContext(v Variables, e []string) Variables {
	c := make(Variables)
	for k, x := range v {
		c[k] = x
	}
	c["std"] = runtime.Stdlib
	c["env"] = mapenv(e)
	return c
}

// Interpolate expressions in a string
func Interpolate(s string, v Variables) (string, error) {
	return interpolate(s, "${", "}", RuntimeContext(v, os.Environ()))
}

// Interpolate every expression in one set of variables in terms of another
func InterpolateAll(a, v Variables) (Variables, error) {
	var err error
	d := make(Variables)
	for k, e := range a {
		if s, ok := e.(string); ok {
			d[k], err = Interpolate(s, v)
			if err != nil {
				return nil, err
			}
		} else {
			d[k] = e
		}
	}
	return d, nil
}

// Interpolate
func interpolate(s, pre, suf string, context interface{}) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("%v: [%s] with context: %s)", err, s, spew.Sdump(context)))
		}
	}()

	if len(pre) < 1 || len(suf) < 1 {
		return "", fmt.Errorf("Invalid variable prefix/suffix")
	}

	fp := pre[0]
	fs := suf[0]

	out := &strings.Builder{}
	var expr *strings.Builder
	var i, esc, start int

	for {
		if i >= len(s) {
			if expr != nil {
				return "", ErrEndOfInput
			} else {
				break
			}
		}

		if s[i] == '\\' {
			esc++
			if (esc % 2) == 0 {
				if expr != nil {
					expr.WriteRune('\\')
				} else {
					out.WriteRune('\\')
				}
			}
			i++
			continue
		}

		if expr != nil {
			if s[i] == fs && (esc%2) == 0 && matchAhead(s[i:], suf) {
				prg, err := epl.Compile(expr.String())
				if err != nil {
					return "", err
				}

				res, err := prg.Exec(context)
				if err != nil {
					return "", newExprError(fmt.Errorf("Could not evaluate expression: %s%v%s: %v", pre, s[start:i], suf, err), context)
				}

				switch v := res.(type) {
				case string:
					out.WriteString(v)
				default:
					out.WriteString(fmt.Sprint(v))
				}

				i += len(suf)
				start = 0
				expr = nil
			} else {
				expr.WriteByte(s[i])
				i++
			}
		} else if s[i] == fp && (esc%2) == 0 && matchAhead(s[i:], pre) {
			i += len(pre)
			start = i
			expr = &strings.Builder{}
		} else {
			out.WriteByte(s[i])
			i++
		}

		esc = 0
	}

	return out.String(), nil
}

// Match ahead
func matchAhead(s, x string) bool {
	if len(s) < len(x) {
		return false
	}
	for i := 0; i < len(x); i++ {
		if s[i] != x[i] {
			return false
		}
	}
	return true
}
