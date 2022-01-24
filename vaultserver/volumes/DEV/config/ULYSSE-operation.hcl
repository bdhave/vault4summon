path "secret/DIGIT/ULYSSE/*" {
  capabilities = ["list"]
}
path "secret/DIGIT/ULYSSE/test/*" {
  capabilities = ["create", "read", "update", "delete"]
}
path "secret/DIGIT/ULYSSE/acc/*" {
  capabilities = ["create", "read", "update", "delete"]
}
path "secret/DIGIT/ULYSSE/prod/*" {
  capabilities = ["create", "update", "delete"]
}
