package main

import (
	"fmt"
	"os"

	"vault4summon/vault"
)

func version() string {
	return "0.9.1"
}

func main() {
	checkArgument()

	argument := os.Args[1]

	// Get the secret and key name from the argument
	switch argument {
	case "-v", "--version":
		printSecret(version())
	default:
		result, err := vault.GetSecret(argument)

		exitIfError(err)
		printSecret(result)
	}
}

func checkArgument() {
	if len(os.Args) != 2 { //nolint:gomnd
		exitIfError(fmt.Errorf("%s", "ERROR: a variable ID or version flag(-v or --version) must be given as the only one argument!")) //nolint:lll
	}
}

func exitIfError(err error) {
	if err == nil {
		return
	}

	_, _ = os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}

func printSecret(result string) {
	_, err := os.Stdout.Write([]byte(result))
	exitIfError(err)
}
