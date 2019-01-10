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

// Copy the environment in this process as a map
func Environ() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		if n := strings.Index(e, "="); n < 0 {
			k := strings.TrimSpace(e[:n])
			env[k] = strings.TrimSpace(e[n+1:])
		}
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

func (p *Process) Redirect(dst io.WriteCloser) error {
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
		_, err := io.Copy(dst, stdout)
		if err != nil && err != io.EOF {
			fmt.Printf("io: %s\n", err)
		}
	}()
	go func() {
		_, err := io.Copy(dst, stderr)
		if err != nil && err != io.EOF {
			fmt.Printf("io: %s\n", err)
		}
	}()

	p.redir = dst
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
	Linger      time.Duration     `yaml:"linger"`
}

// Create a command that inherits its environment from this process
func NewCommand(d, c string) Command {
	return Command{d, c, Environ(), 0}
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
		cmd = exec.CommandContext(cxt, bash, "-c", c.Command)
	} else {
		cmd = exec.Command(bash, "-c", c.Command)
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
func (c Command) Start(dst io.WriteCloser) (*Process, error) {
	cxt, cancel := context.WithCancel(context.Background())
	cmd, err := c.cmd(cxt)
	if err != nil {
		return nil, err
	}

	proc := &Process{sync.Mutex{}, c.Command, cmd, cxt, cancel, nil, c.Linger}

	err = proc.Redirect(dst)
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return proc, nil
}
