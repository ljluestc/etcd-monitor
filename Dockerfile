# Build stage
FROM golang:1.19-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o etcd-monitor cmd/etcd-monitor/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /build/etcd-monitor .

# Expose ports
EXPOSE 8080

# Run the application
ENTRYPOINT ["./etcd-monitor"]
CMD ["--endpoints=localhost:2379", "--api-port=8080"]
