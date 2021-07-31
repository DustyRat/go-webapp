FROM golang:alpine

RUN set -ex; \
    apk update; \
    apk add --update --no-cache openssl-dev musl-dev zlib-dev curl gcc git tzdata;

ADD . /root/
WORKDIR /root/
