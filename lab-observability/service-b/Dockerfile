# Stage 1: Build
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.23-alpine AS builder

# Build arguments for architecture
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -o /api ./cmd/api

# Stage 2: Run
FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.18

WORKDIR /

COPY --from=builder /api /api

EXPOSE 8080

CMD ["/api"]
