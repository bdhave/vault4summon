# This section grants all access on "secret/*". Further restrictions can be
# applied to this broad policy, as shown below.
path "secret/DIGIT/ULYSSE/*" {
  capabilities = ["list"]
}

path "secret/DIGIT/ULYSSE/dev/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

path "secret/DIGIT/ULYSSE/test/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

path "secret/DIGIT/ULYSSE/acc/*" {
  capabilities = ["read"]
}

path "secret/DIGIT/ULYSSE/prod/*" {
  capabilities = ["list"]
}
