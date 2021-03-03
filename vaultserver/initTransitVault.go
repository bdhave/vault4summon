package main

import (
	"fmt"
	"os"
	"vaultserver/command"
)

const defaultVaultAddressTransit = "http://localhost:8200"

func main() {
	command.Setup(defaultVaultAddressTransit)

	var status, err = command.GetStatus()
	command.ExitIfError(err)

	var initialization *command.Initialization
	var fullFileName = command.FullFileName(command.InitializationFilename)
	if !status.Initialized {
		status, initialization, err = command.InitializeTransit(fullFileName)
	} else if status.Sealed {
		_, err := command.Unseal(initialization, fullFileName)
		command.ExitIfError(err)
	}
	token, _, _ := command.CreateToken(command.FullFileName(command.TokenFileName))
	_, _ = fmt.Fprintf(os.Stdout, "%s", token)

	command.ExitIfError(err)
}
