package main

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"os"
	"regexp"
	"strings"
)

func getEnv(env string, defaultValue string) string {
	var value, ok = os.LookupEnv(env)

	if !ok {
		value = defaultValue
	}
	return value
}

func configureClient(vaultUrl string, vaultToken string) (*api.Client, error) {
	config := &api.Config{
		Address: vaultUrl,
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(vaultToken)
	return client, err
}

func getSecrets(argument string, client *api.Client, isVaultEngineV2 bool) (string, *api.Secret, error) {
	path, key := extractKey(argument, isVaultEngineV2)

	secret, err := client.Logical().Read(path)

	if err != nil {
		return "", nil, err
	}

	if secret == nil {
		return "", nil, fmt.Errorf("no secret for path %s\n", path)
	}
	if len(secret.Warnings) > 0 {
		return "", nil, fmt.Errorf("%s\n", strings.Join(secret.Warnings, ". "))
	}

	return key, secret, nil
}

func retrieveValue(secret *api.Secret, field string, isVaultEngineV2 bool) (string, error) {

	data := secret.Data

	if isVaultEngineV2 && data != nil {
		data = nil
		dataRaw := secret.Data["data"]
		if dataRaw != nil {
			data = dataRaw.(map[string]interface{})
		}
	}

	if data == nil {
		return "", fmt.Errorf("no data")
	}

	if entry, ok := data[field]; ok {
		return fmt.Sprintf("%s", entry), nil
	}

	return "field not found", fmt.Errorf("field not found %s", field)
}

func extractKey(argument string, isVaultEngineV2 bool) (string, string) {

	var arguments = strings.SplitN(strings.TrimSpace(argument), "#", 2)
	var isAwsSyntax = false
	var path = arguments[0]
	var key string
	if len(arguments) == 2 {
		key = arguments[1]
		isAwsSyntax = true
	}

	var regexpCompiled = regexp.MustCompile(`/`)
	split := regexpCompiled.Split(path, -1)

	length := len(split)
	if !isAwsSyntax {
		length -= 1
	}
	path = ""
	for i := 0; i < length; i++ {
		path += split[i] + "/"
		if isVaultEngineV2 && i == 0 {
			path += "data/"
		}
	}
	if !isAwsSyntax {
		key = split[length]
	}

	return path, key
}
