package legacy

import (
	"bytes"
	"os"
	"os/exec"
)

type Runtime struct {
	ScriptPath string
}

func New(scriptPath string) Runtime {
	return Runtime{ScriptPath: scriptPath}
}

func (r Runtime) Exec(args []string) int {
	command := exec.Command("/usr/local/bin/python3", append([]string{r.ScriptPath}, args...)...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = os.Environ()

	if err := command.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}

func (r Runtime) Output(args []string) (string, int, error) {
	command := exec.Command("/usr/local/bin/python3", append([]string{r.ScriptPath}, args...)...)
	command.Env = os.Environ()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return stdout.String(), exitErr.ExitCode(), nil
		}
		return stdout.String(), 1, err
	}
	return stdout.String(), 0, nil
}
