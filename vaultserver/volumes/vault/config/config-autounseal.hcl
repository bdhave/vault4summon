# https://learn.hashicorp.com/tutorials/vault/autounseal-transit
disable_mlock = true
ui=true

storage "file" {
  path = "/vault/data"
}

listener "tcp" {
  address     = "[::]:8100"
  tls_disable = "true"
}

seal "transit" {
  address = "http://127.0.0.1:8200"
  disable_renewal = "false"
  key_name = "autounseal"
  mount_path = "transit/"
  tls_skip_verify = "true"
}
