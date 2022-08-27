FROM golang:1.19-buster AS build-env

ENV GO111MODULE=on

RUN mkdir /api

WORKDIR /api

COPY go.mod . 

COPY go.sum . 

RUN go clean --modcache

RUN go mod download

COPY . . 

RUN go mod download

ENTRYPOINT go build  && ./api