# Dockerfile para aplicação Go
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/ordersystem
RUN go build -o /app/ordersystem main.go wire_gen.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ordersystem .
COPY cmd/ordersystem/.env ./
EXPOSE 8080 50051 8000
CMD ["./ordersystem"]