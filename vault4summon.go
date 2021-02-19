package main

import (
	"os"
	"vault4summon/common"
	"vault4summon/vault"
)

func main() {
	err := common.CheckArgument()
	common.ExitIfError(err)

	var result string

	// Get the secret and key name from the argument
	argument := os.Args[1]
	switch argument {
	case "-v", "--version":
		result = vault.Version()
	default:
		result, err = vault.RetrieveSecret(argument)
	}
	common.ExitIfError(err)
	common.PrintSecret(result)
}
