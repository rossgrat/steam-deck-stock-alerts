FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o steam-deck-stock-alerts .

FROM alpine:latest

WORKDIR /app

RUN mkdir -p /var/log/steam-deck-alerts /data

COPY --from=builder /build/steam-deck-stock-alerts .

CMD ["./steam-deck-stock-alerts", "start", "--config", "/app/config.yaml"]
