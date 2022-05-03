package vault

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

type secretID struct {
	Path string
	Key  string
}

func (i secretID) sanitize(argument string) (string, string, error) {
	argument = strings.TrimSpace(argument)
	countCross := strings.Count(argument, "#")
	if countCross > 1 {
		return "", "", fmt.Errorf("SYNTAX ERROR: secretID %q contains %d '#'. Maximum ONE '#' is allowed", argument, countCross)
	}

	path := argument
	var key string
	if countCross == 1 {
		arguments := strings.SplitN(argument, "#", 2)
		path = arguments[0]
		key = arguments[1]
		if strings.Count(key, "/") > 0 {
			return "", "", fmt.Errorf("SYNTAX ERROR: secretID %q (AWS format) with a slash after the '#'", argument)
		}
		if len(key) < 1 {
			return "", "", fmt.Errorf("SYNTAX ERROR: secretID %q (AWS format) ends with '#'", argument)
		}
	}

	parts := strings.Split(path, "/")
	if len(parts) == 1 {
		return "", "", fmt.Errorf("SYNTAX ERROR: secretID %q DOESN'T contain any slash or '#'", argument)
	}

	if len(key) == 0 {
		key = parts[len(parts)-1]
		// remove key form parts
		parts = parts[:len(parts)-1]
	}

	if len(key) < 1 {
		return "", "", fmt.Errorf("SYNTAX ERROR: secretID %q key is empty", argument)
	}

	for i := 0; i < len(parts); i++ {
		if len(parts[i]) < 1 {
			return "", "", fmt.Errorf("SYNTAX ERROR: secretID %q contains leading slash or ending slash or double slashes", argument)
		}
	}

	return strings.Join(parts, "/"), key, nil
}

func newSecretID(argument string) (*secretID, error) {
	id := &secretID{}

	path, key, err := id.sanitize(argument)
	if err != nil {
		return nil, err
	}
	id.Path = path
	id.Key = key

	return id, nil
}

/*
GetSecret calls Hashicorp Vault API to retrieve the secret associated with argument key.
*/
func GetSecret(key string) (string, error) {
	var err error
	var client *api.Client
	client, err = api.NewClient(nil)
	if err != nil {
		return "", err
	}

	// use KvV2 as a first try
	isVaultEngineV2 := true
	var id *secretID
	id, err = newSecretID(key)
	if err != nil {
		return "", err
	}

	var secret *api.Secret
	secret, err = getSecrets(id, client, true)
	if err != nil {
		return "", err
	}

	// check if metadata is present. If true, vault V2 is the right API
	if !isMetadataPresent(secret) {
		// no metadata, so fall to KvV1
		isVaultEngineV2 = false
		secret, err = getSecrets(id, client, false)
		if err != nil {
			return "", err
		}
	}

	return retrieveValue(secret, id.Key, isVaultEngineV2)
}

func isMetadataPresent(secret *api.Secret) bool {
	return secret.Data["metadata"] != nil
}

func getSecrets(id *secretID, client *api.Client, isVaultEngineV2 bool) (*api.Secret, error) {
	path := normalizePath(id.Path, isVaultEngineV2)
	secret, err := client.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret for path %s\n", id.Path)
	}
	if len(secret.Warnings) > 0 {
		return nil, fmt.Errorf("%s\n", strings.Join(secret.Warnings, ". "))
	}

	return secret, nil
}

func retrieveValue(secret *api.Secret, id string, isVaultEngineV2 bool) (string, error) {
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

	if entry, ok := data[id]; ok {
		return fmt.Sprintf("%s", entry), nil
	}

	return "", fmt.Errorf("field not found %q", id)
}

func normalizePath(path string, isVaultEngineV2 bool) string {
	if !isVaultEngineV2 {
		return path
	}
	parts := strings.SplitN(path, "/", 2)
	return parts[0] + "/data/" + parts[1]
}
