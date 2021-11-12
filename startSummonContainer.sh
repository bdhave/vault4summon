docker build -f dockerfile -t summon .

docker run -it --rm \
    --env VAULT_TOKEN=00000000-0000-0000-0000-000000000000 \
    --env VAULT_ADDR=http://vault:8200 \
    --link vault \
    --name summon summon