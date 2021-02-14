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
	ExitIfError(err)

	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Stderr = &bytes.Buffer{}

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

		return nil, newCommandError(exitError, out, command, args)
	}

	return out, nil
}

func ExitIfError(err error) {
	if err == nil {
		return
	}
	var commandError *CommandError
	_, _ = os.Stdout.Write([]byte(err.Error()))
	if errors.As(err, &commandError) {
		os.Exit(commandError.ExitCode)
	}
	os.Exit(-1)
}

type CommandError struct {
	err       *exec.ExitError
	ExitCode  int
	command   string
	args      []string
	output    []byte
	outputErr []byte
}

func (e *CommandError) Error() string {
	description := fmt.Sprintf("When executing %s %s", e.command, strings.Join(e.args, " "))
	var stdout = ""
	if len(e.output) > 0 {
		stdout = "\n" + string(e.output)
	}
	var stderr = ""
	if len(e.outputErr) > 0 {
		stderr = "\nstderr:\n" + string(e.outputErr)
	}
	return fmt.Sprintf("%s:\n%s%s%v", description, stdout, stderr, e.err)
}

func newCommandError(err *exec.ExitError, console []byte, command string, args []string) error {
	if err == nil {
		// As a convenience, if err is nil, newCommandError returns nil.
		return nil
	}
	return &CommandError{err, err.ExitCode(), command, args, console, err.Stderr}
}
