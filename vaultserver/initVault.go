package main

import (
	"vaultserver/command"
)

const defaultVaultAddress = "http://localhost:8100"

func main() {
	command.SetupWithToken(defaultVaultAddress)

	/*	_, err := command.UnWrap()
		command.ExitIfError(err)

		status, err := command.GetStatus()
		command.ExitIfError(err)
	*/
	var _ *command.Initialization
	var fullFileName = command.FullFileName(command.InitializationFilename)
	_, _, err := command.InitializeVault(fullFileName, true)
	command.ExitIfError(err)
}
