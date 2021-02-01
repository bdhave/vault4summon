FROM golang:alpine as builder-base
RUN apk -U upgrade


# compile vault4summon
FROM builder-base as builder
ADD . /source
WORKDIR /source
RUN go mod download all
RUN go build -o target/vault4summon

# create an alpine image with bash, Hashicorp Vault & CyberArk Summon
FROM alpine as alpine-base
RUN apk -U upgrade
RUN apk add curl zsh bash libcap vault git
RUN setcap cap_ipc_lock= /usr/sbin/vault
RUN curl -sSL https://raw.githubusercontent.com/cyberark/summon/master/install.sh | zsh
RUN sh -c "$(curl -fsSL https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
RUN mkdir test


FROM alpine-base
WORKDIR test
COPY --from=builder /source/secrets.yml .

# install provider
RUN mkdir  /usr/local/lib/summon/
COPY --from=builder /source/target/vault4summon  /usr/local/lib/summon/
RUN chmod +x  /usr/local/lib/summon/

ENTRYPOINT /usr/local/bin/summon --yaml 'hello: !var secret/hello#foo' printenv hello
