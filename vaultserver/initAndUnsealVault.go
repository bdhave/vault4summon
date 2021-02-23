package main

import (
	"fmt"
	"os"
	"path/filepath"
	"vaultserver/command"
)

const filename = "Initialization.json"

func main() {
	var fullFileName = filepath.Join(os.Getenv("ROOT4VAULT"), filename)
	const vaultAddr = "VAULT_ADDR"
	if len(os.Getenv(vaultAddr)) == 0 {
		_, _ = os.Stdout.WriteString("VAULT_ADDR environment variable is not defined, set 'http://localhost:8200' as default\n")
		_ = os.Setenv(vaultAddr, "http://localhost:8200")
	}

	var status, err = command.GetStatus()
	command.ExitIfError(err)

	var initialization *command.Initialization
	if !status.Initialized {
		status, initialization, err = command.InitializeTransit(fullFileName)
		command.ExitIfError(err)
	}

	if status.Sealed {
		status, err = command.Unseal(initialization, fullFileName)
		command.ExitIfError(err)
	}

	if status.Sealed {
		command.ExitIfError(fmt.Errorf("cannot unseal vault"))
	}
}
