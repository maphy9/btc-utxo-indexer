FROM golang:1.20-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/maphy9/btc-utxo-indexer
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/btc-utxo-indexer /go/src/github.com/maphy9/btc-utxo-indexer


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/btc-utxo-indexer /usr/local/bin/btc-utxo-indexer
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["btc-utxo-indexer"]
