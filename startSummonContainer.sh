
docker run -it --rm \
    --env VAULT_TOKEN=SP2021 \
    --env VAULT_ADDR=http://vault:8200 \
    --link vault \
    --name summon summon