FROM golang:buster as builder
ARG app_binary=riggs

RUN go get -u github.com/grpc-ecosystem/grpc-health-probe

ENV GO111MODULE=on
WORKDIR /tmp/src/

COPY go.mod  /tmp/src/
COPY go.sum  /tmp/src/
RUN go mod download

COPY .  /tmp/src/
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${app_binary}

FROM debian:buster-slim as final


COPY --from=builder /tmp/src/${app_binary} /app/
COPY --from=builder /go/bin/grpc-health-probe /usr/local/bin/
WORKDIR /app
