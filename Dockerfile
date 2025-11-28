# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod
COPY go.mod ./

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy static files
COPY --from=builder /app/web ./web

# Expose port
EXPOSE 8080

# Run the server
CMD ["./server"]
