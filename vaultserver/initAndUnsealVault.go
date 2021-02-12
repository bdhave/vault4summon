package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const sealedVaultStatusCommandExitCode int = 2

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
	var fullFileName = filepath.Join(os.Getenv("ROOT4VAULT"), filename)
	const vaultAddr = "VAULT_ADDR"
	if len(os.Getenv(vaultAddr)) == 0 {
		_, _ = os.Stdout.WriteString("VAULT_ADDR environment variable is not defined, set 'http://localhost:8200' as default\n")
		_ = os.Setenv(vaultAddr, "http://localhost:8200")
	}

	var status = getStatus()

	var initialization *initialization
	if !status.Initialized {
		initialization = doInit(fullFileName)
		status = getStatus()
	}

	if status.Sealed {
		if initialization == nil {
			dat, err := ioutil.ReadFile(fullFileName)
			exitIfError(err)
			err = json.Unmarshal(dat, initialization)
			exitIfError(err)
		}

		status = doUnseal(initialization)
	}

	if status.Sealed {
		exitIfError(fmt.Errorf("cannot unseal vault"))
	}
}

func getStatus() *status {
	var exitCodes = []int{sealedVaultStatusCommandExitCode}
	jsonData := execute("vault",
		exitCodes,
		"status", "-format", "json")

	status := &status{}
	err := json.Unmarshal(jsonData, status)
	exitIfError(err)
	return status
}

func doInit(fullFileName string) *initialization {
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

	err = ioutil.WriteFile(fullFileName, jsonData, 0644)
	exitIfError(err)

	return initialization
}

func doUnseal(init *initialization) *status {
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

		exitIfError(newCommandError(exitError, out, command, args))
	}

	return out
}

func exitIfError(err error) {
	if err == nil {
		return
	}
	var commandError *CommandError
	_, _ = os.Stdout.Write([]byte(err.Error()))
	if errors.As(err, &commandError) {
		os.Exit(commandError.exitCode)
	}
	os.Exit(-1)
}

type CommandError struct {
	err       *exec.ExitError
	exitCode  int
	command   string
	args      []string
	output    []byte
	outputErr []byte
}

func (e *CommandError) Error() string {
	exitCodeDescription := ""
	if e.exitCode == 1 {
		exitCodeDescription = "\n\texitCode was 1, VAULT description: Local errors such as incorrect flags, failed validations, or wrong numbers of arguments"
	} else if e.exitCode == 2 {
		exitCodeDescription = "\n\texitCode was 2, VAULT description: Any remote errors such as API failures, bad TLS, or incorrect API parameters"
	} else if e.exitCode != 0 {
		exitCodeDescription = fmt.Sprintf("\n\tExitCode was %v", e.exitCode)
	}

	args := strings.Join(e.args, " ")
	description := fmt.Sprintf("When executing VAULT %s %s %s", e.command, args, exitCodeDescription)
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
