ui=true

storage "file" {
  path = "/vault/files"
}

listener "tcp" {
  address     = "[::]:8100"
  tls_disable = 1
}

seal "transit" {
  address = "http://vault:8200"
  disable_renewal = "false"
  key_name = "autounseal"
  mount_path = "transit/"
  tls_skip_verify = "true"
}
