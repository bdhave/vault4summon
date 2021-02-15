package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"vaultserver/common"
)

const sealedVaultStatusCommandExitCode int = 2

type status struct {
	Initialized bool   `json:"initialized"`
	Sealed      bool   `json:"sealed"`
	Version     string `json:"version"`
}

type Initialization struct {
	Seals     []string `json:"unseal_keys_b64"`
	RootToken string   `json:"root_token"`
}

type seal struct {
	Sealed   bool   `json:"sealed"`
	Progress string `json:"progress"`
}

func GetStatus() (*status, error) {
	jsonData, err := common.Execute("vault",
		[]int{sealedVaultStatusCommandExitCode},
		"status", "-format", "json")
	if err != nil {
		return nil, err
	}
	status := &status{}
	err = json.Unmarshal(jsonData, status)
	return status, err
}

func DoInitialization(fullFileName string) (*Initialization, error) {
	var jsonData, err = common.Execute("vault",
		nil,
		"operator", "init", "-format", "json")
	if err = wrapCommandError(err); err != nil {
		return nil, err
	}
	var initialization = &Initialization{}
	err = json.Unmarshal(jsonData, initialization)
	if err != nil {
		return nil, err
	}
	if len(initialization.Seals) < 1 {
		return nil, fmt.Errorf("no seals available")
	}

	jsonData, err = json.Marshal(initialization)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(fullFileName, jsonData, 0644)
	if err != nil {
		return nil, err
	}

	return initialization, err
}

func DoUnseal(initialization *Initialization, fullFileName string) (*status, error) {
	if initialization == nil {
		dat, err := ioutil.ReadFile(fullFileName)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dat, initialization)
		if err != nil {
			return nil, err
		}
	}

	status, err := GetStatus()
	if err != nil {
		return nil, err
	}

	if status.Sealed && initialization != nil && initialization.Seals != nil {
		for i := 0; i < len(initialization.Seals); i++ {
			seal, err := doOneUnseal(initialization.Seals[i])
			if err = wrapCommandError(err); err != nil {
				return nil, err
			}

			if !seal.Sealed {
				break
			}
		}
	}

	status, err = GetStatus()
	return status, err
}

func doOneUnseal(value string) (*seal, error) {
	jsonData, err := common.Execute("vault",
		nil,
		"operator", "unseal", "-format", "json", value)
	if err = wrapCommandError(err); err != nil {
		return nil, err
	}

	seal := &seal{}
	err = json.Unmarshal(jsonData, seal)

	return seal, err
}

func wrapCommandError(err error) error {
	if err == nil {
		return nil
	}
	var commandError *common.CommandError
	if errors.As(err, &commandError) {
		err = newVaultError(commandError)
	}
	return err
}

type vaultError struct {
	err              common.CommandError
	vaultDescription string
}

func (v vaultError) Error() string {
	panic("implement me")
}

func newVaultError(err *common.CommandError) error {
	if err == nil {
		// As a convenience, if err is nil, newCommandError returns nil.
		return nil
	}
	vaultDescription := ""
	if err.ExitCode == 1 {
		vaultDescription = "\n\texitCode was 1, VAULT description: Local errors such as incorrect flags, failed validations, or wrong numbers of arguments"
	} else if err.ExitCode == 2 {
		vaultDescription = "\n\texitCode was 2, VAULT description: Any remote errors such as API failures, bad TLS, or incorrect API parameters"
	} else if err.ExitCode != 0 {
		vaultDescription = fmt.Sprintf("\n\tExitCode was %v", err.ExitCode)
	}
	return &vaultError{*err, vaultDescription}
}

func ExitIfError(err error) {
	if err == nil {
		return
	}
	var commandError *common.CommandError
	if errors.As(err, &commandError) {
		_, _ = os.Stdout.Write([]byte(err.Error()))
		os.Exit(commandError.ExitCode)
	}
	_, _ = os.Stdout.Write([]byte(err.Error()))
	os.Exit(-1)
}
