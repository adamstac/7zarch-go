# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN make build

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates p7zip

# Create non-root user
RUN addgroup -g 1000 -S 7zarch && \
    adduser -u 1000 -S 7zarch -G 7zarch

# Copy binary from builder
COPY --from=builder /app/7zarch-go /usr/local/bin/7zarch-go

# Set up working directory
WORKDIR /data

# Switch to non-root user
USER 7zarch

# Set entrypoint
ENTRYPOINT ["7zarch-go"]

# Default command
CMD ["--help"]