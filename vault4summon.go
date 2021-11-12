package main

import (
	"fmt"
	"os"
	"vault4summon/vault"
)

func main() {
	checkArgument()

	var result string
	var err error
	// Get the secret and key name from the argument
	argument := os.Args[1]
	switch argument {
	case "-v", "--version":
		result = version()
	default:
		result, err = vault.RetrieveSecret(argument)
	}
	exitIfError(err)
	printSecret(result)
}

func version() string {
	return "0.4"
}

func printSecret(result string) {
	_, err := os.Stdout.Write([]byte(result))
	exitIfError(err)
}

func checkArgument() {
	if len(os.Args) != 2 {
		exitIfError(fmt.Errorf("%s", "ERROR: a variable ID or version flag(-v or --version) must be given as the first and only one argument!"))
	}
}

func exitIfError(err error) {
	if err == nil {
		return
	}
	_, _ = os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}
