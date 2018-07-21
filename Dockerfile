FROM golang:alpine AS go-build-env

WORKDIR /go/src/github.com/lon9/kawaiibot
ADD . /go/src/github.com/lon9/kawaiibot

RUN apk add --no-cache git
RUN go get && go build -o /usr/bin/kawaiibot

FROM alpine

WORKDIR /go/src/github.com/lon9/kawaiibot
RUN apk add --no-cache ca-certificates
COPY --from=go-build-env /usr/bin/kawaiibot /usr/bin/kawaiibot