package command

import (
	"os"
	"path/filepath"
)

func FullFileName(fileName string) string {
	rootPath := extractRoot4Vault()
	return filepath.Join(rootPath, fileName)
}

func extractRoot4Vault() string {
	const root4vault = "ROOT4VAULT"
	rootPath := os.Getenv(root4vault)
	if len(rootPath) < 1 {
		_, _ = os.Stdout.WriteString(root4vault + " environment variable is not defined, set './' as default\n")
		rootPath = "./"
	}
	_, _ = os.Stdout.WriteString("Using " + rootPath + " as " + root4vault + "\n")
	return rootPath
}
