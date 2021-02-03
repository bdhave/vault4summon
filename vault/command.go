package vault

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"strings"
	"vault4summon/common"
)

func Version() string {
	return "0.1"
}

func RetrieveSecret(argument string) (string, error) {
	client, err := configureClient()
	if err != nil {
		return "", err
	}

	// use KvV2 as a first try
	var isVaultEngineV2 = true
	var variableId = common.NewVariableID(argument)
	secret, err := getSecrets(variableId, client, true)
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

func configureClient() (*api.Client, error) {
	//todo use default config
	// todo remove default values, may be load more values needed dor a hardened vault deployment
	const EnvVaultAddress = "VAULT_ADDR"
	const EnvVaultToken = "VAULT_TOKEN"
	var vaultUrl = common.GetEnv(EnvVaultAddress, "http://localhost:8200")
	var vaultToken = common.GetEnv(EnvVaultToken, "myRootToken")

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

func getSecrets(variableID *common.VariableID, client *api.Client, isVaultEngineV2 bool) (*api.Secret, error) {

	var path = normalizePath(variableID.Path, isVaultEngineV2)
	secret, err := client.Logical().Read(path)

	if err != nil {
		return secret, err
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
	if len(parts) < 2 {
		common.PrintAndExit(fmt.Errorf("%d", "variableID path  %q MUST contains at least one '/' .", path))
	}

	return parts[0] + "/data/" + parts[1]
}
