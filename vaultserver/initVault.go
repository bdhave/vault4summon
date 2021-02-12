package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const sealedVaultStatusCommandExitCode = 2

type status struct {
	Initialized bool   `json:"initialized"`
	Sealed      bool   `json:"sealed"`
	Version     string `json:"version"`
}

type initialization struct {
	Seals     []string `json:"unseal_keys_b64"`
	RootToken string   `json:"root_token"`
}

type seal struct {
	Sealed   bool   `json:"sealed"`
	Progress string `json:"progress"`
}

const filename = "initialization.json"

func main() {

	var status = getStatus()
	var initialization *initialization
	if !status.Initialized {
		initialization = doInit()
		status = getStatus()
	}

	if status.Sealed {
		status = doUnseal(initialization)
	}

	if status.Sealed {
		exitIfError(fmt.Errorf("cannot unseal vault"))
	}
}

func getStatus() *status {
	exitCodes := append(make([]int, 1), sealedVaultStatusCommandExitCode)
	jsonData := execute("vault",
		exitCodes,
		"status", "-format", "json")

	status := &status{}
	err := json.Unmarshal(jsonData, status)
	exitIfError(err)

	return status
}

func doInit() *initialization {
	jsonData := execute("vault",
		nil,
		"operator", "init", "-format", "json")
	var initialization = &initialization{}
	err := json.Unmarshal(jsonData, initialization)
	if err != nil || len(initialization.Seals) < 1 {
		exitIfError(fmt.Errorf("no seals available"))
	}

	jsonData, err = json.Marshal(initialization)
	exitIfError(err)
	//enc := json.NewEncoder(os.Stdout)

	err = ioutil.WriteFile(filename, jsonData, 0644)
	exitIfError(err)

	return initialization
}

func doUnseal(init *initialization) *status {

	if init == nil {
		dat, err := ioutil.ReadFile(filename)
		exitIfError(err)
		err = json.Unmarshal(dat, init)
		exitIfError(err)
	}

	if getStatus().Sealed && init != nil && init.Seals != nil {
		for i := 0; i < len(init.Seals); i++ {
			seal := doOneUnseal(init.Seals[i])
			if !seal.Sealed {
				break
			}
		}
	}

	return getStatus()
}

func doOneUnseal(value string) *seal {
	jsonData := execute("vault",
		nil,
		"operator", "unseal", "-format", "json", value)
	seal := &seal{}
	_ = json.Unmarshal(jsonData, seal)

	return seal
}

func execute(command string, ignoredExitCode []int, args ...string) []byte {
	var err error
	command, err = exec.LookPath(command)
	exitIfError(err)

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
						return out
					}
				}
			}
		}

		exitIfError(newCommandError(exitError))
	}

	return out
}

func exitIfError(err error) {
	if err == nil {
		return
	}

	var exitCode int = 1
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		exitCode = exitError.ExitCode()
		_ = sliceOutputStream(exitError.Stderr)
		_, _ = os.Stderr.Write([]byte(err.Error()))
	}
	os.Exit(exitCode)

}
func sliceOutputStream(out []byte) []string {
	if out == nil {
		return nil
	}
	return strings.Split(string(out), "\n")
}

type CommandError struct {
	err *exec.ExitError
}

func (e *CommandError) Error() string {
	exitCode := e.err.ExitCode()
	if exitCode == 1 {
		// Local errors such as incorrect flags, failed validations, or wrong numbers of arguments
	}
	if exitCode == 2 {
		// Any remote errors such as API failures, bad TLS, or incorrect API parameters
	}

	return ": " + e.err.Error()

}

func newCommandError(err *exec.ExitError) error {
	if err == nil {
		// As a convenience, if err is nil, NewSyscallError returns nil.
		return nil
	}
	return &CommandError{err}
}

/*

type WrappedError struct {
    Context string
    Err     error
}

func (w *WrappedError) Error() string {
    return fmt.Sprintf("%s: %v", w.Context, w.Err)
}

func Wrap(err error, info string) *WrappedError {
    return &WrappedError{
        Context: info,
        Err:     err,
    }
}

*/
