package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func GetStatus() *status {
	var exitCodes = []int{sealedVaultStatusCommandExitCode}
	jsonData, err := common.Execute("vault",
		exitCodes,
		"status", "-format", "json")
	//todo check err

	status := &status{}
	err = json.Unmarshal(jsonData, status)
	common.ExitIfError(err)
	return status
}

func DoInitialization(fullFileName string) *Initialization {
	var jsonData, err = common.Execute("vault",
		nil,
		"operator", "init", "-format", "json")

	//todo check err
	var initialization = &Initialization{}
	err = json.Unmarshal(jsonData, initialization)
	if err != nil || len(initialization.Seals) < 1 {
		common.ExitIfError(fmt.Errorf("no seals available"))
	}

	jsonData, err = json.Marshal(initialization)
	common.ExitIfError(err)
	//enc := json.NewEncoder(os.Stdout)

	err = ioutil.WriteFile(fullFileName, jsonData, 0644)
	common.ExitIfError(err)

	return initialization
}

func DoUnseal(initialization *Initialization, fullFileName string) *status {
	if initialization == nil {
		dat, err := ioutil.ReadFile(fullFileName)
		common.ExitIfError(err)
		err = json.Unmarshal(dat, initialization)
		common.ExitIfError(err)
	}

	if GetStatus().Sealed && initialization != nil && initialization.Seals != nil {
		for i := 0; i < len(initialization.Seals); i++ {
			seal := doOneUnseal(initialization.Seals[i])
			if !seal.Sealed {
				break
			}
		}
	}

	return GetStatus()
}

func doOneUnseal(value string) *seal {
	jsonData, err := common.Execute("vault",
		nil,
		"operator", "unseal", "-format", "json", value)
	//todo check err
	if err != nil {

	}
	seal := &seal{}
	_ = json.Unmarshal(jsonData, seal)

	return seal
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
