package main

import (
	"fmt"
	"vaultserver/command"
)

const defaultVaultAddressTransit = "http://localhost:8200"
const tokenFileName = "config/token.json"

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

	initialization, err = command.ReadInitialization(initialization, fullFileName)
	command.ExitIfError(err)

	token, err := command.CreateToken(command.FullFileName(tokenFileName))
	command.ExitIfError(err)
	//fmt.Printf("tokenInfo: %v", tokenInfo)
	_, _ = fmt.Printf("%s", token)

	command.ExitIfError(err)
}
