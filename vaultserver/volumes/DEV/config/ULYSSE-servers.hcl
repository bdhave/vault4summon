path "secret/DIGIT/ULYSSE/*" {
  capabilities = ["read"]
}
path "secret/DIGIT/ULYSSE/prod/*" {
  capabilities = ["deny"]
}