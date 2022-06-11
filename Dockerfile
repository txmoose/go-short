# syntax=docker/dockerfile:1

FROM golang:1.18 AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY *.go ./
RUN GOOS=linux GOARCH=arm64 go build -o /go-short 
RUN chmod +x /go-short
EXPOSE 8000

# DEPLOY
FROM multiarch/ubuntu-core:arm64-bionic

WORKDIR /

COPY --from=build /go-short /go-short
ENV GS_SLUG_LENGTH=4

EXPOSE 8000

ENTRYPOINT ["/go-short"]
