package vault

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

type secretID struct {
	Path      string
	Key       string
	awsFormat bool
}

func newSecretID(argument string) (*secretID, error) {
	var id = &secretID{} //nolint:gofumpt
	path, key, awsFormat, err := id.sanitize(argument)

	if err != nil {
		return nil, err
	}

	id.Path = path
	id.Key = key
	id.awsFormat = awsFormat

	return id, nil
}

func (i secretID) sanitize(argument string) (string, string, bool, error) {
	argument = strings.TrimSpace(argument)
	countCross := strings.Count(argument, "#")
	awsFormat := countCross > 0

	if countCross > 1 {
		return "", "", awsFormat, newSecretIDError(argument, "contains %d '#'. Maximum ONE '#' is allowed", countCross)
	}

	key := ""
	path := argument

	if countCross == 1 {
		arguments := strings.SplitN(argument, "#", 2) //nolint:gomnd
		path = arguments[0]
		key = arguments[1]

		if strings.Count(key, "/") > 0 {
			return "", "", awsFormat, newSecretIDError(argument, "(AWS format) with a slash after the '#'")
		}

		if len(key) < 1 {
			return "", "", awsFormat, newSecretIDError(argument, "(AWS format) ends with '#'", argument)
		}
	}

	parts := strings.Split(path, "/")
	if len(parts) == 1 {
		return "", "", awsFormat, newSecretIDError(argument, "DOESN'T contain any slash or '#'")
	}

	if len(key) == 0 {
		key = parts[len(parts)-1]
		// remove key form parts
		parts = parts[:len(parts)-1]
	}

	if len(key) < 1 {
		return "", "", awsFormat, newSecretIDError(argument, "key is empty")
	}

	for i := 0; i < len(parts); i++ {
		if len(parts[i]) < 1 {
			return "", "", awsFormat, newSecretIDError(argument, "contains leading slash or ending slash or double slashes")
		}
	}

	return strings.Join(parts, "/"), key, awsFormat, nil
}

type secretIDError struct {
	secretID string
	msg      string
}

func (e *secretIDError) Error() string {
	return fmt.Sprintf("ERROR: secretId %q %s", e.secretID, e.msg)
}

func newSecretIDError(secretID string, msg string, args ...any) error {
	return &secretIDError{secretID, fmt.Sprintf(msg, args)}
}

/*
GetSecret calls Hashicorp Vault API to retrieve the secret associated with argument key.
*/
func GetSecret(key string) (string, error) {
	client, err := api.NewClient(nil)
	if err != nil {
		return "", err
	}

	// use KvV2 as a first try
	isVaultEngineV2 := true
	secretID, err := newSecretID(key)

	if err != nil {
		return "", err
	}

	var secret *api.Secret
	secret, err = getSecrets(secretID, client, true)

	if err != nil {
		return "", err
	}

	// check if metadata is present. If true, vault V2 is the right API
	if !isMetadataPresent(secret) {
		// no metadata, so fall to KvV1
		isVaultEngineV2 = false
		secret, err = getSecrets(secretID, client, false)

		if err != nil {
			return "", err
		}
	}

	return retrieveValue(secret, secretID.Key, isVaultEngineV2)
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

func retrieveValue(secret *api.Secret, key string, isVaultEngineV2 bool) (string, error) {
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

	if entry, ok := data[key]; ok {
		return fmt.Sprintf("%s", entry), nil
	}

	return "", fmt.Errorf("field not found %q", key)
}

func normalizePath(path string, isVaultEngineV2 bool) string {
	if !isVaultEngineV2 {
		return path
	}

	parts := strings.SplitN(path, "/", 2) //nolint:gomnd

	return parts[0] + "/data/" + parts[1]
}
