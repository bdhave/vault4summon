vault audit enable file file_path=vault/logs/audit/audit.log
vault secrets enable transit
vault write -f transit/keys/autounseal
vault policy write/vault/config/autounseal autounseal.hcl
