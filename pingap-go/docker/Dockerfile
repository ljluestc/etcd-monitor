# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pingap ./cmd/pingap

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/pingap .
COPY --from=builder /build/configs/example.toml ./config.toml

# Create directories for logs and data
RUN mkdir -p /var/log/pingap /opt/pingap

# Expose ports
EXPOSE 80 443 3000 9090

# Run as non-root user
RUN addgroup -g 1000 pingap && \
    adduser -D -u 1000 -G pingap pingap && \
    chown -R pingap:pingap /app /var/log/pingap /opt/pingap

USER pingap

ENTRYPOINT ["./pingap"]
CMD ["-config", "/app/config.toml"]
