docker build -f dockerfile -t summon .

docker run -it --rm \
    --env VAULT_TOKEN=myRootToken \
    --env VAULT_ADDR=http://vault:8200 \
    --link vault \
    --name summon summon