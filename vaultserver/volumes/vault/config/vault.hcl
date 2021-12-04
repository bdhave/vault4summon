ui=true
log_level = "warn"
default_lease_ttl = "768h"


storage "file" {
  path = "/vault/files"
}

listener "tcp" {
  address     = "[::]:8100"
  tls_disable = "true"
}

api_addr = "http://localhost:8100"

seal "transit" {
#  address = "http://vault:8200"
  address = "http://localhost:8200"
  disable_renewal = "true"
  key_name = "autounseal"
  mount_path = "transit/"
  tls_skip_verify = "true"
}
