FROM golang:1.18-alpine AS builder

WORKDIR /wikiwoop-chain

ADD go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o wikiwoop-chain main.go

FROM alpine:latest AS runner

RUN apk --no-cache add ca-certificates jq

WORKDIR /wikiwoop-chain

COPY --from=builder /wikiwoop-chain/wikiwoop-chain ./
COPY ./docker/local/wikiwoop-chain.sh ./

# Expose json-rpc, libp2p and grpc ports
EXPOSE 8545 9632 1478

ENTRYPOINT ["./wikiwoop-chain.sh"]
