# syntax=docker/dockerfile:1

FROM golang:1.18 AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY config.yaml ./
COPY *.go ./
RUN go build -o /go-short
EXPOSE 8000

# DEPLOY
FROM ubuntu:20.04

WORKDIR /

COPY config.yaml /config.yaml
COPY --from=build /go-short /go-short

EXPOSE 8000

ENTRYPOINT ["/go-short"]
