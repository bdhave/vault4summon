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
const VaultAddress = "http://localhost:8200"

const auditFilename = "audit/audit.log"
const tokenFileName = "token.json"

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
	Token           string `json:"token"`
	Accessor        string `json:"accessor"`
	WrappedAccessor string `json:"wrapped_accessor"`
	TTL             int    `json:"ttl"`
}

type token struct {
	Duration  int       `json:"lease_duration"`
	Renewable bool      `json:"renewable"`
	Tokens    tokenInfo `json:"wrap_info"`
}

func Setup(address string) {
	const vaultAddr = "VAULT_ADDR"
	if len(os.Getenv(vaultAddr)) == 0 {
		_, _ = os.Stdout.WriteString(vaultAddr + " environment variable is not defined, set 'http://localhost:8200' as default\n")
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

	err = os.Setenv("VAULT_TOKEN", initialization.RootToken)
	if err != nil {
		return nil, nil, err
	}

	jsonData, err = json.Marshal(initialization)
	if err != nil {
		return nil, initialization, err
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

	_, err = CreateToken(FullFileName(tokenFileName))
	return err
}

func CreateToken(fullFileName string) (*token, error) {
	jsonData, err := common.Execute("vault",
		nil,
		"token", "create", "-policy=\"autounseal\"", "-wrap-ttl=120", "-format", "json")
	if err = wrapCommandError(err); err != nil {
		return nil, err
	}
	token := &token{}
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

func wrapCommandError(err error) error {
	if err == nil {
		return nil
	}
	var commandError *common.CommandError
	if errors.As(err, &commandError) {
		return newVaultError(commandError)
	}
	return err
}

type vaultError struct {
	err              common.CommandError
	vaultDescription string
}

func (v vaultError) Error() string {
	return fmt.Sprintf("Vault ERROR:\n%s", v.err.Error)
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

func FullFileName(fileName string) string {
	const root4vault = "ROOT4VAULT"
	rootPath := os.Getenv(root4vault)
	if len(rootPath) < 1 {
		_, _ = os.Stdout.WriteString(root4vault + " environment variable is not defined, set './' as default\n")
		rootPath = "./"
	}
	return filepath.Join(rootPath, fileName)
}
