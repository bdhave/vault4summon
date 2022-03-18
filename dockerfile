FROM golang:alpine as golang-base
RUN apk -U upgrade && apk cache clean
workdir temp
COPY go.mod .
RUN go mod download
WORKDIR ..

# compile vault4summon
FROM golang-base as builder
ADD . /source
WORKDIR /source
RUN go build -o target/vault4summon

# create an alpine image with bash, Hashicorp Vault & CyberArk Summon
FROM alpine:3.15 as alpine-base
RUN apk -U upgrade && \
    apk add bash libcap vault git openssl && \
    apk cache clean && \
    setcap cap_ipc_lock= /usr/sbin/vault && \
    curl -sSL https://raw.githubusercontent.com/cyberark/summon/master/install.sh | bash && \
    mkdir test

FROM alpine-base
WORKDIR test
COPY --from=builder /source/secrets.yml .

# install provider
COPY --from=builder /source/target/vault4summon  /usr/local/lib/summon/

ENTRYPOINT /usr/local/bin/summon --yaml 'hello: !var secret/hello#foo' printenv hello
