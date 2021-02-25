ui=true

storage "file" {
  path = "/vault/files"
}

listener "tcp" {
  address     = "[::]:8100"
  tls_disable = "true"
}

seal "transit" {
  address = "http://vault:8200"
  disable_renewal = "true"
  key_name = "autounseal"
  mount_path = "transit/"
  tls_skip_verify = "true"
}
