package main

import (
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Pluralize
func plural(v int, s, p string) string {
	if v == 1 {
		return s
	} else {
		return p
	}
}

// Return the first non-empty string from those provided
func coalesce(v ...string) string {
	for _, e := range v {
		if e != "" {
			return e
		}
	}
	return ""
}

// String to bool
func strToBool(s string, d ...bool) bool {
	if s == "" {
		if len(d) > 0 {
			return d[0]
		} else {
			return false
		}
	}
	return strings.EqualFold(s, "t") || strings.EqualFold(s, "true") || strings.EqualFold(s, "y") || strings.EqualFold(s, "yes")
}

// String to int
func strToInt(s string, d ...int) int {
	if s == "" {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}
	return v
}

// String to duration
func strToDuration(s string, d ...time.Duration) time.Duration {
	if s == "" {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}
	return v
}

// Dump environment pairs
func dumpEnv(w io.Writer, env map[string]string) {
	wk := 0
	for k, _ := range env {
		if l := len(k); l < 40 && l > wk {
			wk = l
		}
	}
	f := fmt.Sprintf("        %%%ds = %%s\n", wk)
	for k, v := range env {
		fmt.Fprintf(w, f, k, v)
	}
}

// Disambiguate a filename
func disambigFile(base, ext string, counts map[string]int) string {
	stem := base[:len(base)-len(path.Ext(base))]
	n, ok := counts[stem]
	if ok && n > 0 {
		stem = fmt.Sprintf("%v-%d", stem, n)
	}
	counts[stem] = n + 1
	return stem + ext
}

type colorWriter struct {
	io.WriteCloser
	attrs []color.Attribute
}

func (w colorWriter) Write(p []byte) (int, error) {
	return color.New(w.attrs...).Fprint(w.WriteCloser, string(p))
}
