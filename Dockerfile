FROM golang:1.23.8-alpine

ADD pkg /root/pkg
ADD cmd /root/cmd
ADD config/config-docker.yaml /root/config.yaml

ADD go.mod /root/go.mod
ADD go.sum /root/go.sum

EXPOSE 8080

WORKDIR /root

ARG VERSION=dev
ENV VERSION=${VERSION}

RUN go build -ldflags "-X main.version=$VERSION" -o /root/main cmd/*
CMD ["/root/main", "-f", "/root/config.yaml"]
