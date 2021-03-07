package command

import (
	"os"
	"path/filepath"
)

func FullFileName(fileName string) string {
	const root4vault = "ROOT4VAULT"
	rootPath := os.Getenv(root4vault)
	if len(rootPath) < 1 {
		_, _ = os.Stdout.WriteString(root4vault + " environment variable is not defined, set './' as default\n")
		rootPath = "./"
	}
	return filepath.Join(rootPath, fileName)
}
