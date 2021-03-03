package main

import (
	"vaultserver/command"
)

const defaultVaultAddress = "http://localhost:8100"

func main() {
	command.SetupWithToken(defaultVaultAddress)

	_, err := command.UnWrap()
	command.ExitIfError(err)

	status, err := command.GetStatus()
	command.ExitIfError(err)

	var initialization *command.Initialization
	var fullFileName = command.FullFileName(command.InitializationFilename)
	if !status.Initialized {
		status, initialization, err = command.InitializeVault(fullFileName)
	} else if status.Sealed {
		_, err := command.Unseal(initialization, fullFileName)
		command.ExitIfError(err)
	}
	command.ExitIfError(err)
}
