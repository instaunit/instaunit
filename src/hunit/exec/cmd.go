package exec

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const newline = '\n'

// A writer that indents lines with a prefix string
type prefixWriter struct {
	writer io.Writer
	prefix string
	nl     bool
}

func (s prefixWriter) Write(p []byte) (int, error) {
	var i, x, n int
	var r rune

	if s.nl {
		z, err := s.writer.Write([]byte(s.prefix))
		if err != nil {
			return z, err
		}
		s.nl = false
	}

	c := string(p)
	for i, r = range c {
		if s.nl {
			z, err := s.writer.Write([]byte(s.prefix))
			if err != nil {
				return z, err
			}
			s.nl = false
			// don't account for the prefix, this presumably would confuse callers if: n > len(p)
		}
		if r == newline {
			z, err := s.writer.Write([]byte(c[x : i+1]))
			if err != nil {
				return z, err
			}
			x = i + 1
			n += z
			s.nl = true
		}
	}

	if x < i {
		z, err := s.writer.Write([]byte(c[x:]))
		if err != nil {
			return z, err
		}
		n += z
		if c[len(c)-1] == newline {
			s.nl = true
		}
	}

	return n, nil
}

func (s prefixWriter) Close() error {
	return nil // do nothing; this is intended to be used with os.Stdout
}

// Create a prefix writer
func NewPrefixWriter(w io.Writer, p string) prefixWriter {
	return prefixWriter{w, p, true}
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
	cancel  context.CancelFunc
	redir   io.WriteCloser
	linger  time.Duration
}

func (p *Process) Redirect(out io.WriteCloser) error {
	p.Lock()
	defer p.Unlock()
	if p.cmd == nil {
		return fmt.Errorf("No process")
	}
	if p.redir != nil {
		return fmt.Errorf("Output already redirected")
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
		_, err := io.Copy(out, stdout)
		if err != nil && err != io.EOF {
			fmt.Printf("io: %s\n", err)
		}
	}()
	go func() {
		_, err := io.Copy(out, stderr)
		if err != nil && err != io.EOF {
			fmt.Printf("io: %s\n", err)
		}
	}()

	p.redir = out
	return nil
}

func (p Process) Linger() time.Duration {
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
	if p.cancel != nil {
		p.cancel()
	} else {
		return fmt.Errorf("Process is not cancelable")
	}
	if p.redir != nil {
		p.redir.Close()
	}
	p.cmd = nil // mark the process as cancelled
	p.redir = nil
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
func (c Command) Start(out io.WriteCloser) (*Process, error) {
	cxt, cancel := context.WithCancel(context.Background())
	cmd, err := c.cmd(cxt)
	if err != nil {
		return nil, err
	}

	proc := &Process{sync.Mutex{}, c.Command, cmd, cxt, cancel, nil, c.Linger}

	err = proc.Redirect(out)
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return proc, nil
}
