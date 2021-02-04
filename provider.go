package main

import (
	"os"
	"vault4summon/common"
	"vault4summon/vault"
)

func main() {
	common.CheckArgument()

	var result string
	var err error

	// Get the secret and key name from the argument
	argument := os.Args[1]
	switch argument {
	case "-v", "--version":
		result = vault.Version()
	default:
		result, err = vault.RetrieveSecret(argument)
	}

	if err != nil {
		common.Exit(err)

	}
	common.PrintSecret(result)
}
