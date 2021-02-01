package main

import (
	"fmt"
	"os"
)

func main() {
	ok, err := checkArgument()
	var result string

	if ok {
		// Get the secret and key name from the argument
		argument := os.Args[1]
		switch argument {
		case "-v", "--version":
			result = Version()
		default:
			result, err = RetrieveSecret(argument)
		}
	}

	if err != nil {
		printAndExit(err)

	}
	_, _ = os.Stdout.Write([]byte(result))
}

func checkArgument() (bool, error) {
	if len(os.Args) != 2 {
		return false, fmt.Errorf("%s", "A variable ID or version flag must be given as the first and only argument!")
	}
	return true, nil
}

func printAndExit(err error) {
	_, _ = os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}
