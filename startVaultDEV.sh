docker run --cap-add=IPC_LOCK -e 'VAULT_DEV_ROOT_TOKEN_ID=myRootToken' -p 127.0.0.1:8200:8200/tcp --name vault vault
