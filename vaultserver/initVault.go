package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

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
		exitIfError(fmt.Errorf("cannot unseal vault"), nil)
	}
}

func getStatus() *status {
	exitCodes := append(make([]int, 3), 0, 1, 2)
	jsonData := execute("vault",
		exitCodes,
		"status", "-format", "json")

	status := &status{}
	err := json.Unmarshal(jsonData, status)
	exitIfError(err, nil)

	return status
}

func doInit() *initialization {
	jsonData := execute("vault",
		nil,
		"operator", "init", "-format", "json")
	var initialization = &initialization{}
	err := json.Unmarshal(jsonData, initialization)
	if err != nil || len(initialization.Seals) < 1 {
		exitIfError(fmt.Errorf("no seals available"), nil)
	}

	jsonData, err = json.Marshal(initialization)
	exitIfError(err, nil)
	//enc := json.NewEncoder(os.Stdout)

	err = ioutil.WriteFile(filename, jsonData, 0644)
	exitIfError(err, nil)

	return initialization
}

func doUnseal(init *initialization) *status {

	if init == nil {
		dat, err := ioutil.ReadFile(filename)
		exitIfError(err, nil)
		err = json.Unmarshal(dat, init)
		exitIfError(err, nil)
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
	exitIfError(err, nil)

	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Stderr = &bytes.Buffer{}

	out, err := cmd.Output()
	if err != nil {
		exitError := err.(*exec.ExitError)
		exitCode := exitError.ExitCode()
		if ignoredExitCode != nil {
			for i := 0; i < len(ignoredExitCode); i++ {
				if ignoredExitCode[i] == exitCode {
					return out
				}
			}
		}
		if exitCode == 1 {

		}
		if exitCode == 2 {

		}
		/*
			Local errors such as incorrect flags, failed validations, or wrong numbers of arguments return an exit code of 1.
			Any remote errors such as API failures, bad TLS, or incorrect API parameters return an exit status of 2
		*/
		exitIfError(err, exitError.Stderr)
	}

	return out
}

func exitIfError(err error, errOutput []byte) {
	if err == nil {
		return
	}
	sliceOutputStream(errOutput)
	_, _ = os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}

func sliceOutputStream(out []byte) []string {
	if out == nil {
		return nil
	}
	return strings.Split(string(out), "\n")
}
