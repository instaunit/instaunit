package exec

import (
	"fmt"
	"os"
	"os/exec"
)

// A command
type Command struct {
	Display     string            `yaml:"display"`
	Command     string            `yaml:"run"`
	Environment map[string]string `yaml:"environment"`
}

// Execute a command
func (c Command) Exec() (string, error) {
	var bash string
	if v := os.Getenv("BASH"); v != "" {
		bash = v
	} else if v, err := exec.LookPath("bash"); err == nil {
		bash = v
	} else {
		return "", fmt.Errorf("You must have the Bash Shell in your path somewhere (or set $BASH in your environment)")
	}

	cmd := exec.Command(bash, "-c", c.Command)
	if len(c.Environment) > 0 {
		env := os.Environ()
		for k, v := range c.Environment {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = env
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), err
}
