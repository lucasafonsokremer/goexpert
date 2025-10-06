FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/api .

ENV PORT=8080
ENV WEATHER_SERVICE_URL=http://service-b:8080

EXPOSE 8080
CMD ["./api"]