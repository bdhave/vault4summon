package main

const EnvVaultAddress = "VAULT_ADDR"
const EnvVaultAgentAddr = "VAULT_AGENT_ADDR"
const EnvVaultCACert = "VAULT_CACERT"
const EnvVaultCAPath = "VAULT_CAPATH"
const EnvVaultClientCert = "VAULT_CLIENT_CERT"
const EnvVaultClientKey = "VAULT_CLIENT_KEY"
const EnvVaultClientTimeout = "VAULT_CLIENT_TIMEOUT"
const EnvVaultSkipVerify = "VAULT_SKIP_VERIFY"
const EnvVaultNamespace = "VAULT_NAMESPACE"
const EnvVaultTLSServerName = "VAULT_TLS_SERVER_NAME"
const EnvVaultWrapTTL = "VAULT_WRAP_TTL"
const EnvVaultMaxRetries = "VAULT_MAX_RETRIES"
const EnvVaultToken = "VAULT_TOKEN"
const EnvVaultMFA = "VAULT_MFA"
const EnvRateLimit = "VAULT_RATE_LIMIT"

func Version() string {
	return "0.1"
}

// todo remove default values, may be load more values needed dor a hardened vault deplyment
var vaultUrl = getEnv(EnvVaultAddress, "http://localhost:8200")
var vaultToken = getEnv(EnvVaultToken, "myRootToken")

func RetrieveSecret(argument string) (string, error) {
	client, err := configureClient(vaultUrl, vaultToken)
	if err != nil {
		return "", err
	}

	// use KvV2 as a first try
	var isVaultEngineV2 = true
	key, secret, err := getSecrets(argument, client, true)
	if err != nil {
		return "", err
	}

	// check if metadata is present. If true, vault V2 is the right API

	metadataRaw := secret.Data["metadata"]
	if metadataRaw == nil {
		// no metadata, so fall to KvV1
		isVaultEngineV2 = false
		key, secret, err = getSecrets(argument, client, isVaultEngineV2)
		if err != nil {
			return "", err
		}
	}

	return retrieveValue(secret, key, isVaultEngineV2)
}
