FROM golang:1.19.0-alpine3.16

RUN apk update && \
    apk upgrade && \
    apk add git

RUN go install github.com/cespare/reflex
ENV CGO_ENABLED=0

WORKDIR /go/src/api
COPY ./api ./

RUN go mod download