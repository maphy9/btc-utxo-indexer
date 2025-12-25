FROM golang:1.24-alpine AS buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/maphy9/btc-utxo-indexer

COPY go.mod go.sum ./
RUN go mod download

# COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/btc-utxo-indexer /go/src/github.com/maphy9/btc-utxo-indexer


FROM alpine:3.23

COPY --from=buildbase /usr/local/bin/btc-utxo-indexer /usr/local/bin/btc-utxo-indexer
RUN apk add --no-cache ca-certificates

ENV KV_VIPER_FILE=/app/config/config.yaml

ENTRYPOINT ["btc-utxo-indexer"]