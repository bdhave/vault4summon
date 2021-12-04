# This section grants all access on "secret/*". Further restrictions can be
# applied to this broad policy, as shown below.

path "secret/DIGIT/ULYSSE/prod/*" {
  capabilities = ["read"]
}
