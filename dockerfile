FROM golang:alpine as builder-base
RUN apk -U upgrade
ADD . /source
WORKDIR /source
RUN go mod download all


# compile vault4summon
FROM builder-base as builder
#ADD . /source
WORKDIR /source
RUN go build -o target/vault4summon
RUN go build -o target/test4summon test/test4summon.go

# create an alpine image with bash, Hashicorp Vault & CyberArk Summon
FROM alpine as alpine-base
RUN apk -U upgrade
RUN apk add curl zsh bash libcap vault git
RUN setcap cap_ipc_lock= /usr/sbin/vault
RUN curl -sSL https://raw.githubusercontent.com/cyberark/summon/master/install.sh | zsh
RUN sh -c "$(curl -fsSL https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"


FROM alpine-base
RUN mkdir test
WORKDIR test
COPY --from=builder /source/secrets.yml .
COPY --from=builder /source/target/test4summon .
COPY --from=builder /source/target/vault4summon .
RUN chmod +x .

# install provider
RUN mkdir  /usr/local/lib/summon/
COPY --from=builder /source/target/vault4summon  /usr/local/lib/summon/
RUN chmod +x  /usr/local/lib/summon/

#ENTRYPOINT /usr/local/bin/summon -D ENV=dev ./test4summon

#/usr/local/bin/summon -D ENV=dev ./test4summon