ui = true
log_level = "warn"
default_lease_ttl = "768h"

storage "file" {
  path = "/vault/files"
}

listener "tcp" {
  address     = "[::]:8200"
  tls_disable = 1
}

api_addr = "http://localhost:8200"