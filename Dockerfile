FROM golang:1.24.0-alpine

ADD pkg /root/pkg
ADD cmd /root/cmd
ARG CONFIG_FILE=config-docker.yaml
ADD config/${CONFIG_FILE} /root/config.yaml

ADD go.mod /root/go.mod
ADD go.sum /root/go.sum

EXPOSE 8080

WORKDIR /root

ARG VERSION=dev
ENV VERSION=${VERSION}

RUN go build -ldflags "-X main.version=$VERSION" -o /root/main cmd/*
CMD ["/root/main", "-f", "/root/config.yaml"]
