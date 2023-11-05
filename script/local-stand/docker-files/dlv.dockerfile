FROM golang:1.20.7-alpine3.18

ENV CGO_ENABLED 0

RUN apk add --no-cache git bash

WORKDIR /go/src/

RUN go install github.com/go-delve/delve/cmd/dlv@v1.21.0
RUN mkdir "/.cache" && chmod -R 0777 /.cache
