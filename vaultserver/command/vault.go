package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"vaultserver/common"
)

const InitializationFilename = "config/initialization.json"
const auditFilename = "audit/audit.log"
const TokenFileName = "token.json"

const vaultAddr = "VAULT_ADDR"
const vaultToken = "VAULT_TOKEN"

type status struct {
	Type         string `json:"type"`
	Initialized  bool   `json:"initialized"`
	Sealed       bool   `json:"sealed"`
	Version      string `json:"version"`
	Nonce        string `json:"nonce"`
	RecoverySeal bool   `json:"recovery_seal"`
	StorageType  string `json:"storage_type"`
	HaEnabled    bool   `json:"ha_enabled"`
}

type Initialization struct {
	Seals        []string `json:"unseal_keys_b64"`
	RecoveryKeys []string `json:"recovery_keys_b64"`
	RootToken    string   `json:"root_token"`
}

type seal struct {
	Progress     int    `json:"progress"`
	Type         string `json:"type"`
	Initialized  bool   `json:"initialized"`
	Sealed       bool   `json:"sealed"`
	Version      string `json:"version"`
	Nonce        string `json:"nonce"`
	RecoverySeal bool   `json:"recovery_seal"`
	StorageType  string `json:"storage_type"`
	HaEnabled    bool   `json:"ha_enabled"`
}

type tokenInfo struct {
	Duration  int          `json:"lease_duration"`
	Renewable bool         `json:"renewable"`
	Token     tokenWrapped `json:"wrap_info"`
}

type tokenWrapped struct {
	Token           string `json:"token"`
	Accessor        string `json:"accessor"`
	WrappedAccessor string `json:"wrapped_accessor"`
	TTL             int    `json:"ttl"`
}

func SetupWithToken(address string) {
	Setup(address)

	if len(os.Args) > 2 {
		ExitIfError(errors.New("this application accepts ONE and ONLY ONE argument, the token"))
	}
	if len(os.Getenv(vaultToken)) == 0 && len(os.Args) != 2 {
		ExitIfError(fmt.Errorf(
			"%s environment variable is not defined, you MUST gives the token as argument or define %s environment variable",
			vaultToken, vaultToken))
	}
	if len(os.Getenv(vaultToken)) == 0 {
		_ = os.Setenv(vaultToken, os.Args[1])
	}

	if len(os.Getenv(vaultToken)) == 0 {
		ExitIfError(fmt.Errorf(
			"%s environment variable is not defined, you MUST gives the token as argument or define %s environment variable",
			vaultToken, vaultToken))
	}
}

func Setup(address string) {
	if len(os.Getenv(vaultAddr)) == 0 {

		_, _ = fmt.Fprintf(os.Stderr, "%s environment variable is not defined, set '%s' as default\n", vaultAddr, address)
		_ = os.Setenv(vaultAddr, address)
	}
}

func GetStatus() (*status, error) {
	const sealedVaultStatusCommandExitCode int = 2

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

func InitializeTransit(fullFileName string) (*status, *Initialization, error) {
	status, initialization, err := initialize(fullFileName)
	if err != nil {
		return status, initialization, err
	}

	if len(initialization.Seals) < 1 {
		return nil, initialization, fmt.Errorf("no seals available")
	}
	status, err = Unseal(initialization, fullFileName)
	err = EnableAudit(filepath.Join(os.Getenv("ROOT4VAULT"), auditFilename))
	ExitIfError(err)
	err = enableTransit()
	ExitIfError(err)
	status, err = GetStatus()
	ExitIfError(err)
	return status, initialization, err
}

func InitializeVault(fullFileName string) (*status, *Initialization, error) {
	status, initialization, err := initialize(fullFileName)
	if err != nil {
		return status, initialization, err
	}
	return status, initialization, nil
}

func initialize(fullFileName string) (*status, *Initialization, error) {
	var jsonData, err = common.Execute("vault",
		nil,
		"operator", "init", "-format", "json")
	if err = wrapCommandError(err); err != nil {
		return nil, nil, err
	}
	var initialization = &Initialization{}
	err = json.Unmarshal(jsonData, initialization)
	if err != nil {
		return nil, nil, err
	}

	err = os.Setenv(vaultToken, initialization.RootToken)
	if err != nil {
		return nil, nil, err
	}

	err = ioutil.WriteFile(fullFileName, jsonData, 0644)
	if err != nil {
		return nil, initialization, err
	}
	status, err := GetStatus()
	if err != nil {
		return nil, initialization, err
	}
	return status, initialization, err
}

func EnableAudit(fullFileName string) error {
	var _, err = common.Execute("vault",
		nil,
		"audit", "enable", "file", "file_path="+fullFileName)
	if err = wrapCommandError(err); err != nil {
		return err
	}
	return err
}

func enableTransit() error {
	var _, err = common.Execute("vault",
		nil,
		"secrets", "enable", "transit")
	if err = wrapCommandError(err); err != nil {
		return err
	}
	_, err = common.Execute("vault",
		nil,
		"write", "-f", "transit/keys/autounseal")
	if err = wrapCommandError(err); err != nil {
		return err
	}
	_, err = common.Execute("vault",
		nil,
		"policy", "write", "transit-policy", FullFileName("config/autounseal-policy.hcl"))

	if err = wrapCommandError(err); err != nil {
		return err
	}

	return nil
}

func CreateToken(fullFileName string) (string, *tokenInfo, error) {
	token, err := createToken(fullFileName)
	ExitIfError(err)
	return token.Token.Token, token, nil
}

func createToken(fullFileName string) (*tokenInfo, error) {
	jsonData, err := common.Execute("vault",
		nil,
		"token", "create", "-policy=\"autounseal\"", "-format", "json")
	if err = wrapCommandError(err); err != nil {
		return nil, err
	}
	token := &tokenInfo{}
	err = json.Unmarshal(jsonData, token)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(fullFileName, jsonData, 0644)
	if err != nil {
		return nil, err
	}
	return token, err
}

func UnWrap() (string, error) {
	jsonData, err := common.Execute("vault",
		nil,
		"unwrap", "-format", "json")
	if err = wrapCommandError(err); err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func Unseal(initialization *Initialization, fullFileName string) (*status, error) {
	status, err := GetStatus()
	if err != nil || !status.Sealed {
		return status, err
	}

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

	if initialization != nil {
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
	ExitIfError(err)

	if status.Sealed {
		ExitIfError(fmt.Errorf("cannot unseal vault"))
	}
	return status, nil
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
