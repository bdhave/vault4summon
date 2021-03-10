package common

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Execute(command string, ignoredExitCode []int, args ...string) ([]byte, error) {
	var err error
	command, err = exec.LookPath(command)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	out, err := cmd.Output()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode := exitError.ExitCode()
			if ignoredExitCode != nil {
				for i := 0; i < len(ignoredExitCode); i++ {
					if ignoredExitCode[i] == exitCode {
						return out, nil
					}
				}
			}
		}

		return nil, newCommandError(exitError, out, stderr.Bytes(), command, args)
	}

	return out, nil
}

type CommandError struct {
	err      *exec.ExitError
	ExitCode int
	command  string
	args     []string
	stdout   []byte
	stderr   []byte
}

func (e *CommandError) Error() string {
	description := fmt.Sprintf("When executing %s %s", e.command, strings.Join(e.args, " "))
	var stdout = ""
	if len(e.stdout) > 0 {
		stdout = "\n" + string(e.stdout)
	}
	var stderr = ""
	if len(e.stderr) > 0 {
		stderr = "\nstderr:\n" + string(e.stderr)
	}
	return fmt.Sprintf("%s:\nstdout\n%s\nstdout\n%s\n%v", description, stdout, stderr, e.err)
}

func newCommandError(err *exec.ExitError, console []byte, stderr []byte, command string, args []string) error {
	if err == nil {
		// As a convenience, if err is nil, newCommandError returns nil.
		return nil
	}
	return &CommandError{err, err.ExitCode(), command, args, console, stderr}
}
