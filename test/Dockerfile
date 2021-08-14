FROM golang:alpine

RUN set -ex; \
    apk update; \
    apk add --update --no-cache openssl-dev musl-dev zlib-dev curl gcc git tzdata;

RUN go install github.com/jstemmer/go-junit-report@latest
RUN go install github.com/boumenot/gocover-cobertura@latest

ADD . /root/
WORKDIR /root/
