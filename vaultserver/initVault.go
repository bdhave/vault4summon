package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type status struct {
	initialized bool
	sealed      bool
	version     string
}

type initialization struct {
	Seals       []string
	RootToken   string
	Initialized bool
	Sealed      bool
}

type seal struct {
	sealed   bool
	progress int64
}

const filename = "/vault/config/initialization.json"

func main() {

	var status = getStatus()
	var initialization *initialization
	if !status.initialized {
		initialization = doInit()
		status = getStatus()
	}

	if status.sealed {
		status = doUnseal(initialization)
	}

	if status.sealed {
		exitIfError(fmt.Errorf("cannot useal vault"), nil)
	}
}

func getStatus() *status {
	results := execute("vault", "status")

	var initialized, _ = strconv.ParseBool(parseResult(results, "Initialized", " "))
	var sealed, _ = strconv.ParseBool(parseResult(results, "Sealed", " "))

	status := &status{
		initialized,
		sealed,
		parseResult(results, "Version", " "),
	}

	return status
}

func doInit() *initialization {
	results := execute("vault", "init")

	var seals []string
	for i := 1; i < len(results); i++ {
		if strings.Contains(results[i], "Unseal Key") {
			seals = append(seals, strings.Split(results[i], ":")[1])
		}
	}
	initialized, _ := strconv.ParseBool(parseResult(results, "Initialized", ":"))
	sealed, _ := strconv.ParseBool(parseResult(results, "Sealed", ":"))

	initialization := &initialization{
		seals,
		parseResult(results, "Root Token", ":"),
		initialized,
		sealed,
	}
	if initialization.Initialized && initialization.Sealed && len(initialization.Seals) < 1 {
		exitIfError(fmt.Errorf("no seals available"), nil)
	}

	var jsonData, err = json.Marshal(initialization)
	exitIfError(err, nil)
	err = ioutil.WriteFile(filename, jsonData, 0644)
	exitIfError(err, nil)

	return initialization
}

func doUnseal(init *initialization) *status {

	if init == nil {
		init = &initialization{}
		dat, err := ioutil.ReadFile(filename)
		exitIfError(err, nil)
		err = json.Unmarshal(dat, init)
		exitIfError(err, nil)
	}

	if init.Sealed {
		for i := 0; i < len(init.Seals); i++ {
			seal := doOneUnseal(init.Seals[i])
			if !seal.sealed {
				break
			}
		}
	}

	return getStatus()
}

func doOneUnseal(value string) *seal {
	results := execute("vault", "init", "unseal", value)
	var sealed, _ = strconv.ParseBool(parseResult(results, "Sealed", ":"))
	var progress, _ = strconv.ParseInt(parseResult(results, "Progress", ":"), 10, 32)
	seal := &seal{
		sealed,
		progress,
	}

	return seal
}

func exitIfError(err error, errOutput []byte) {
	if err == nil {
		return
	}
	sliceOutputStream(errOutput)
	_, _ = os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}

func execute(command string, args ...string) []string {
	var err error
	command, err = exec.LookPath(command)
	exitIfError(err, nil)

	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Stderr = &bytes.Buffer{}

	var out []byte
	out, err = cmd.Output()
	exitIfError(err, out)

	return sliceOutputStream(out)
}

func sliceOutputStream(out []byte) []string {
	if out == nil {
		return nil
	}
	return strings.Split(string(out), "\n")
}

func parseResult(array []string, key string, separator string) string {
	for i := 0; i < len(array); i++ {
		if strings.Contains(array[i], key) {
			var parts = strings.Split(array[i], separator)
			return strings.TrimSpace(parts[1])
		}
	}
	return ""
}
