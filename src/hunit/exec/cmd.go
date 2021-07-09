package exec

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const newline = '\n'

// A writer that discards output
type discardWriter struct{}

func (s discardWriter) Write(p []byte) (int, error) { return len(p), nil }
func (s discardWriter) Close() error                { return nil }

func NewDiscardWriter() discardWriter {
	return discardWriter{}
}

// A writer that indents lines with a prefix string
type prefixWriter struct {
	writer io.Writer
	buffer *bytes.Buffer
	prefix string
}

func (s prefixWriter) Write(p []byte) (int, error) {
	n, t := 0, len(p)
	for {
		r, w := utf8.DecodeRune(p)
		if r == utf8.RuneError {
			if w == 0 {
				break
			} else {
				return n, fmt.Errorf("Invalid UTF-8")
			}
		}

		n += w
		p = p[w:]

		s.buffer.WriteRune(r)
		if r == newline {
			_, err := s.flush()
			if err != nil {
				return n, err
			}
		}
	}
	return t, nil
}

func (s prefixWriter) flush() (int, error) {
	l := s.buffer.Len()
	v := make([]byte, len(s.prefix)+l)
	copy(v, []byte(s.prefix))
	copy(v[len(s.prefix):], s.buffer.Bytes())
	_, err := s.writer.Write(v)
	if err != nil {
		return 0, err
	}
	s.buffer.Reset()
	return l, nil
}

func (s prefixWriter) Close() error {
	if s.buffer.Len() > 0 {
		_, err := s.flush()
		if err != nil {
			return err
		}
	}
	return nil // do nothing; this is intended to be used with os.Stdout
}

// Create a prefix writer
func NewPrefixWriter(w io.Writer, p string) prefixWriter {
	return prefixWriter{w, &bytes.Buffer{}, p}
}

// Copy the environment in this process as a map and merge it with the
// provided set, which takes prescidence
func Environ(sup map[string]string) map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		if n := strings.Index(e, "="); n > 0 {
			k := strings.TrimSpace(e[:n])
			env[k] = strings.TrimSpace(e[n+1:])
		}
	}
	for k, v := range sup {
		env[k] = v
	}
	return env
}

// A running process
type Process struct {
	sync.Mutex
	cmdline string
	cmd     *exec.Cmd
	context context.Context
	wout    io.WriteCloser
	werr    io.WriteCloser
	linger  time.Duration
	exited  bool
	status  *os.ProcessState
	monitor chan struct{}
}

func (p *Process) Start(wout, werr io.WriteCloser) error {
	p.Lock()
	defer p.Unlock()
	if p.cmd == nil {
		return fmt.Errorf("No process")
	}
	if p.wout != nil {
		return fmt.Errorf("Process already started")
	}

	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		_, err := io.Copy(wout, stdout)
		if err != nil && err != io.EOF {
			fmt.Printf("io: %s\n", err)
		}
	}()
	go func() {
		_, err := io.Copy(werr, stderr)
		if err != nil && err != io.EOF {
			fmt.Printf("io: %s\n", err)
		}
	}()

	p.wout = wout
	p.werr = werr

	err = p.cmd.Start()
	if err != nil {
		return err
	}

	go func(cmd *exec.Cmd) {
		cmd.Wait()
		p.Lock()
		defer p.Unlock()
		p.exited = true
		p.status = cmd.ProcessState
		if p.monitor != nil {
			p.monitor <- struct{}{}
		}
	}(p.cmd)

	return nil
}

func (p *Process) Running() bool {
	p.Lock()
	defer p.Unlock()
	return !p.exited
}

func (p *Process) Monitor() *os.ProcessState {
	p.Lock()
	m := p.monitor
	p.Unlock()

	if m != nil {
		<-m
	} else {
		return p.status
	}

	p.Lock()
	defer p.Unlock()

	if p.monitor != nil {
		close(p.monitor)
		p.monitor = nil
	}

	return p.status
}

func (p *Process) Linger() time.Duration {
	return p.linger
}

func (p *Process) Kill() error {
	p.Lock()
	defer p.Unlock()
	if p.cmd == nil {
		return nil // already cancelled or not running
	}
	if p.linger > 0 {
		<-time.After(p.linger)
	}
	if !p.exited {
		if p.cmd.Process != nil {
			err := p.cmd.Process.Kill()
			if err != nil {
				return fmt.Errorf("Could not kill process: %v", err)
			}
		} else {
			return fmt.Errorf("No process")
		}
	}
	if p.wout != nil {
		p.wout.Close()
		p.wout = nil
	}
	if p.werr != nil {
		p.werr.Close()
		p.werr = nil
	}
	p.cmd = nil // mark the process as cancelled
	p.exited = true
	return nil
}

func (p *Process) String() string {
	return p.cmdline
}

// A command
type Command struct {
	Display     string            `yaml:"display"`
	Command     string            `yaml:"run"`
	Environment map[string]string `yaml:"environment"`
	Wait        time.Duration     `yaml:"wait"`
	Linger      time.Duration     `yaml:"linger"`
}

// Create a command that inherits its environment from this process
func NewCommand(d, c string) Command {
	return Command{d, c, Environ(nil), 0, 0}
}

// Prepare a command but do not execute it
func (c Command) cmd(cxt context.Context) (*exec.Cmd, error) {
	var bash string
	if v := os.Getenv("BASH"); v != "" {
		bash = v
	} else if v, err := exec.LookPath("bash"); err == nil {
		bash = v
	} else {
		return nil, fmt.Errorf("You must have the Bash Shell in your path somewhere (or set $BASH in your environment)")
	}

	var cmd *exec.Cmd
	if cxt != nil {
		cmd = exec.CommandContext(cxt, bash, "-e", "-o", "pipefail", "-c", c.Command)
	} else {
		cmd = exec.Command(bash, "-e", "-o", "pipefail", "-c", c.Command)
	}
	if len(c.Environment) > 0 {
		env := os.Environ()
		for k, v := range c.Environment {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = env
	}

	return cmd, nil
}

// Execute a command
func (c Command) Exec() (string, error) {
	cmd, err := c.cmd(nil)
	if err != nil {
		return "", err
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), err
}

// Start a command
func (c Command) Start(wout, werr io.WriteCloser) (*Process, error) {
	cxt := context.Background()
	cmd, err := c.cmd(cxt)
	if err != nil {
		return nil, err
	}

	proc := &Process{sync.Mutex{}, c.Command, cmd, cxt, nil, nil, c.Linger, false, nil, make(chan struct{}, 1)}
	err = proc.Start(wout, werr)
	if err != nil {
		return nil, err
	}
	return proc, nil
}
