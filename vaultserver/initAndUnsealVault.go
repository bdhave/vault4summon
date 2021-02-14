package main

import (
	"fmt"
	"os"
	"path/filepath"
	"vaultserver/command"
	"vaultserver/common"
)

const filename = "Initialization.json"

func main() {
	var fullFileName = filepath.Join(os.Getenv("ROOT4VAULT"), filename)
	const vaultAddr = "VAULT_ADDR"
	if len(os.Getenv(vaultAddr)) == 0 {
		_, _ = os.Stdout.WriteString("VAULT_ADDR environment variable is not defined, set 'http://localhost:8200' as default\n")
		_ = os.Setenv(vaultAddr, "http://localhost:8200")
	}

	var status = command.GetStatus()

	var initialization *command.Initialization
	if !status.Initialized {
		initialization = command.DoInitialization(fullFileName)
		status = command.GetStatus()
	}

	if status.Sealed {
		status = command.DoUnseal(initialization, fullFileName)
	}

	if status.Sealed {
		common.ExitIfError(fmt.Errorf("cannot unseal vault"))
	}
}
