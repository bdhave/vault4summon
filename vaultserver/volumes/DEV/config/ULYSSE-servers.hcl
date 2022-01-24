path "secret/DIGIT/ULYSSE/*" {
  capabilities = ["list", "read"]
}
path "secret/DIGIT/ULYSSE/prod/*" {
  capabilities = ["deny"]
}