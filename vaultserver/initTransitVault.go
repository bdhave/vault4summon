package main

import (
	"vaultserver/command"
)

func main() {
	command.Setup(command.VaultAddress)

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
	command.ExitIfError(err)
}
