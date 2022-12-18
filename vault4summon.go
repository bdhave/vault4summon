package main

import (
	_ "embed"
	"fmt"
	"os"
	"vault4summon/vault"
)

// generate version from last tag created
//
//go:generate bash -c "printf %s $(git describe --tags --abbrev=0) > version.txt"
//go:embed version.txt
var version string

func main() {
	checkArgument()

	argument := os.Args[1]

	// Get the secret and key name from the argument
	switch argument {
	case "-v", "--version":
		printVersion()
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

func printSecret(secret string) {
	_, err := os.Stdout.Write([]byte(secret))
	exitIfError(err)
}

func printVersion() {
	if len(version) == 0 {
		version = "dev"

	}
	_, err := os.Stdout.Write([]byte(version))
	exitIfError(err)
}
