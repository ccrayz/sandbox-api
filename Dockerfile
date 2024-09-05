FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ccrayz-cli .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/ccrayz-cli .

ENTRYPOINT ["./ccrayz-cli", "server"]