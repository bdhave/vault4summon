package common

import (
	"fmt"
	"os"
)

func CheckArgument() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("%s", "A variable ID or version flag(-v or --version) must be given as the first and only one argument!")
	}
	return nil
}

func ExitIfError(err error) {
	if err == nil {
		return
	}
	_, _ = os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}

func PrintSecret(result string) {
	_, _ = os.Stdout.Write([]byte(result))
}
