package vault

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"strings"
)

func RetrieveSecret(argument string) (string, error) {
	var err error
	var client *api.Client
	client, err = api.NewClient(nil)
	if err != nil {
		return "", err
	}

	// use KvV2 as a first try
	var isVaultEngineV2 = true
	var variableId *VariableID
	variableId, err = NewVariableID(argument)
	if err != nil {
		return "", err
	}

	var secret *api.Secret
	secret, err = getSecrets(variableId, client, true)
	if err != nil {
		return "", err
	}

	// check if metadata is present. If true, vault V2 is the right API
	if !isMetadataPresent(secret) {
		// no metadata, so fall to KvV1
		isVaultEngineV2 = false
		secret, err = getSecrets(variableId, client, false)
		if err != nil {
			return "", err
		}
	}

	return retrieveValue(secret, variableId.Key, isVaultEngineV2)
}

func isMetadataPresent(secret *api.Secret) bool {
	return secret.Data["metadata"] != nil
}

func getSecrets(variableID *VariableID, client *api.Client, isVaultEngineV2 bool) (*api.Secret, error) {
	var path = normalizePath(variableID.Path, isVaultEngineV2)
	secret, err := client.Logical().Read(path)

	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret for path %s\n", variableID.Path)
	}
	if len(secret.Warnings) > 0 {
		return nil, fmt.Errorf("%s\n", strings.Join(secret.Warnings, ". "))
	}

	return secret, nil
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

	return "", fmt.Errorf("field not found %q", field)
}

func normalizePath(path string, isVaultEngineV2 bool) string {
	if !isVaultEngineV2 {
		return path
	}
	var parts = strings.SplitN(path, "/", 2)
	return parts[0] + "/data/" + parts[1]
}
