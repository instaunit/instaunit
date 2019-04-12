package debug

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"
)

var (
	DEBUG   bool
	VERBOSE bool
	TRACE   bool
)

var previous time.Time

const threshold = time.Second * 3

var sourceRoots []string

func init() {
	sourceRoots = make([]string, 0)
	for _, e := range strings.Split(os.Getenv("GOPATH"), ":") {
		e = strings.TrimSpace(e)
		if e != "" {
			sourceRoots = append(sourceRoots, e)
		}
	}
	// init defaults from the environment
	DEBUG = istrue(os.Getenv("DEBUG"), os.Getenv("GOUTIL_DEBUG"))
	VERBOSE = istrue(os.Getenv("VERBOSE"), os.Getenv("GOUTIL_VERBOSE"))
	TRACE = istrue(os.Getenv("TRACE"), os.Getenv("GOUTIL_TRACE"))
}

func relativeSourcePath(p string) string {
	for _, e := range sourceRoots {
		if strings.HasPrefix(p, e) {
			if len(p) > len(e) {
				return p[len(e)+1:]
			} else {
				return "/"
			}
		}
	}
	return p
}

type Frame struct {
	File string
	Path string
	Line int
	Name string
	Func *runtime.Func
}

func (f Frame) String() string {
	return fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Name)
}

func Stacktrace() []Frame {
	pc := make([]uintptr, 64)
	n := runtime.Callers(2, pc)
	t := make([]Frame, n)
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])
		t[i] = Frame{relativeSourcePath(file), file, line, f.Name(), f}
	}
	return t
}

func CurrentContext() string {
	pc := make([]uintptr, 2)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return fmt.Sprintf("%s:%d %s", relativeSourcePath(file), line, f.Name())
}

func DumpRoutines() {
	data := make([]byte, 1<<20)
	n := runtime.Stack(data, true)
	io.Copy(os.Stderr, bytes.NewReader(data[:n]))
}

func DumpRoutinesOnInterrupt() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, os.Kill)
	go func() {
		for e := range sig {
			if time.Since(previous) < threshold {
				log.Printf("\nSecond signal, exiting...\n")
				os.Exit(0) // just exit, it's the second in a series
			}

			log.Printf("\nReceived a signal, dumping stack...\n")
			DumpRoutines()

			previous = time.Now()
			if e == os.Kill {
				os.Exit(0)
			}
		}
	}()
}

func istrue(v ...string) bool {
	for _, e := range v {
		if strings.EqualFold("true", e) || strings.EqualFold("t", e) || strings.EqualFold("yes", e) || strings.EqualFold("y", e) {
			return true
		}
	}
	return false
}
