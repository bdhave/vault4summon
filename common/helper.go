package common

import (
	"fmt"
	"os"
)

// todo remove
func GetEnv(env string, defaultValue string) string {
	var value, ok = os.LookupEnv(env)

	if !ok {
		value = defaultValue
	}
	return value
}

func CheckArgument() {
	if len(os.Args) != 2 {
		PrintAndExit(fmt.Errorf("%s", "A variable ID or version flag must be given as the first and only argument!"))
	}
}

func PrintAndExit(err error) {
	_, _ = os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}

func PrintSecret(result string) {
	_, _ = os.Stdout.Write([]byte(result))
}
