FROM golang:1.7-alpine

# Install base packages
RUN apk update && apk upgrade && apk add \
    curl \
    git \
    bash \
    vim

RUN curl https://glide.sh/get | sh

RUN mkdir -p /go/src/github.com/vgarvardt/rklotz
RUN mkdir -p /data/bin
RUN mkdir -p /data/db
RUN mkdir -p /data/templates
RUN mkdir -p /data/static

WORKDIR /go/src/github.com/vgarvardt/rklotz

EXPOSE 8080