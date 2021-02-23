disable_mlock = true
ui=true

storage "file" {
  path = "/vault/files"
}

listener "tcp" {
  address     = "[::]:8200"
  tls_disable = 1
}