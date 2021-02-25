package command

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"vaultserver/common"
)

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
