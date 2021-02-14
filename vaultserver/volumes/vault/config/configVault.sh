#Execute the following command to unwrap the secrets passed from Vault 1.

VAULT_TOKEN=<wrapping_token> vault unwrap

export VAULT_TOKEN=<wrapping_token> vault unwrap
#export VAULT_TOKEN="s.SJj086AW7ZaobvRmNoSjhVaj"

vault server -config=config-autounseal.hcl

VAULT_ADDR=http://127.0.0.1:8100 vault operator init -recovery-shares=1 \
        -recovery-threshold=1 > recovery-key.txt


vault server -config=config-autounseal.hcl