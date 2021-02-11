#export ROOT4VAULT=vaultserver/volumes/transit
export ROOT4VAULT=/vault

vault audit enable file file_path=$ROOT4VAULT/logs/audit/audit.log
vault secrets enable transit
vault write -f transit/keys/autounseal
vault policy write transit-policy $ROOT4VAULT/config/autounseal.hcl

