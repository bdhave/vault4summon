#Execute the following command to unwrap the secrets passed from Vault 1.
Now, start a second Vault instance which listens to port 8100. The server configuration file should define a seal stanza with parameters properly set based on the tasks you performed in Step 1.

Scenario Overview

Execute the following command to unwrap the secrets passed from Vault 1.

$ VAULT_TOKEN=<wrapping_token> vault unwrap

$ VAULT_TOKEN="s.AFqDxN5jdiDQDNuodJxsC6dm" vault unwrap

Key                  Value
---                  -----
token                s.SJj086AW7ZaobvRmNoSjhVaj
token_accessor       6lUHdlwRZfyRnHtNiz8ZB1Jx
token_duration       768h
token_renewable      true
token_policies       ["autounseal" "default"]
identity_policies    []
policies             ["autounseal" "default"]
Copy
The revealed token is the client token Vault 2 will use to connect with Vault 1.

Set VAULT_TOKEN environment variable whose value is the client token you just unwrapped.

Example:

$ export VAULT_TOKEN="s.SJj086AW7ZaobvRmNoSjhVaj"
Copy
Create a server configuration file (config-autounseal.hcl) to start a second Vault instance (Vault 2).

disable_mlock = true
ui=true

storage "file" {
  path = "/vault-2/data"
}

listener "tcp" {
  address     = "127.0.0.1:8100"
  tls_disable = "true"
}

seal "transit" {
  address = "http://127.0.0.1:8200"
  disable_renewal = "false"
  key_name = "autounseal"
  mount_path = "transit/"
  tls_skip_verify = "true"
}
Copy
Notice that the address points to the Vault server listening to port 8200 (Vault 1). The key_name and mount_path match to what you created in Step 1.

NOTE: The seal stanza does not set the token value since it's already set as VAULT_TOKEN environment variable.

Although the listener stanza disables TLS (tls_disable = "true") for this tutorial, Vault should always be used with TLS in production to provide secure communication between clients and the Vault server. It requires a certificate file and key file on each Vault host.

Start the vault server with the configuration file.

$ vault server -config=config-autounseal.hcl
Copy
Open another terminal and initialize your second Vault server (Vault 2).

$ VAULT_ADDR=http://127.0.0.1:8100 vault operator init -recovery-shares=1 \
        -recovery-threshold=1 > recovery-key.txt
Copy
By passing the VAULT_ADDR, the subsequent command gets executed against the second Vault server (http://127.0.0.1:8100).

Notice that you are setting the number of recovery key and recovery threshold because there is no unseal keys with auto-unseal. Vault 2's master key is now protected by the transit secrets engine of Vault 1. Recovery keys are used for high-privilege operations such as root token generation. Recovery keys are also used to make Vault operable if Vault has been manually sealed through the "vault operator seal" command.

Check the Vault 2 server status. It is now successfully initialized and unsealed.

$ VAULT_ADDR=http://127.0.0.1:8100 vault status

Key                      Value
---                      -----
Recovery Seal Type       shamir
Initialized              true
Sealed                   false
Total Recovery Shares    1
Threshold                1
# ...snip...
Copy
Notice that it shows Total Recovery Shares instead of Total Shares. The transit secrets engine is solely responsible for protecting the master key of Vault 2. There are some operations that still requires Shamir's keys (e.g. regenerate a root token). Therefore, Vault 2 server requires recovery keys although auto-unseal has been enabled.

»Step 3: Verify Auto-Unseal
To verify that Vault 2 gets automatically unseal, press Ctrl + C to stop the Vault 2 server where it is running.

...snip...
[INFO]  core.cluster-listener: rpc listeners successfully shut down
[INFO]  core: cluster listeners successfully shut down
[INFO]  core: vault is sealed
Copy
Note that Vault 2 is now sealed.

Press the upper-arrow key, and execute the vault server -config=config-autounseal.hcl command again to start Vault 2 and see what happens.

$ vault server -config=config-autounseal.hcl

==> Vault server configuration:

                Seal Type: transit
          Transit Address: http://127.0.0.1:8200
        Transit Key Name: autounseal
      Transit Mount Path: transit/
                      Cgo: disabled
              Listener 1: tcp (addr: "0.0.0.0:8100", cluster address: "0.0.0.0:8101", max_request_duration: "1m30s", max_request_size: "33554432", tls: "disabled")
                Log Level: info
                    Mlock: supported: true, enabled: false
                  Storage: file
                  Version: Vault v1.1.0
              Version Sha: 36aa8c8dd1936e10ebd7a4c1d412ae0e6f7900bd

==> Vault server started! Log data will stream in below:

[WARN]  no `api_addr` value specified in config or in VAULT_API_ADDR; falling back to detection if possible, but this value should be manually set
[INFO]  core: stored unseal keys supported, attempting fetch
[INFO]  core: vault is unsealed
# ...snip...
Copy
Notice that the Vault server is already unsealed. The Transit Address is set to your Vault 1 which is listening to port 8200 (http://127.0.0.1:8200).

Check the Vault 2 server status.

$ VAULT_ADDR=http://127.0.0.1:8100 vault status

Key                      Value
---                      -----
Recovery Seal Type       shamir
Initialized              true
Sealed                   false
Total Recovery Shares    1
Threshold                1
# ...snip...
Copy
Now, examine the audit log in Vault 1.

$ tail -f audit.log | jq

# ...snip...
"request": {
  "id": "a46719eb-eee0-92a4-2da6-6c7de77fd410",
  "operation": "update",
  "client_token": "hmac-sha256:ce8613487054dadb36a9d08da1f5a4bbee2fbfc1ef1ec5ebdeec696df7823e69",
  "client_token_accessor": "hmac-sha256:f3b6cb798605835e8a00bafa9e0e16fc0534b8923b31e499f2c8e694f6b69158",
  "namespace": {
    "id": "root",
    "path": ""
  },
  "path": "transit/decrypt/autounseal",
    # ...snip...
  "remote_address": "127.0.0.1",
  "wrap_ttl": 0,
  "headers": {}
},
# ...snip...
}
Copy
You should see an update request against the transit/decrypt/autounseal path. The remote_address is 127.0.0.1 in this example since Vault 1 and Vault 2 are both running locally. If the Vault 2 is running on a different host, the audit log will show the IP address of the Vault 2 host.

