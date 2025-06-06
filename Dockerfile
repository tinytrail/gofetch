# Build stage
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates for dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o fetch-mcp-server ./cmd/server

# Runtime stage - use distroless for minimal attack surface
FROM gcr.io/distroless/static-debian12:nonroot

# Copy CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from builder stage
COPY --from=builder /app/fetch-mcp-server /usr/local/bin/fetch-mcp-server

# Use non-root user (already set in distroless:nonroot)
# USER 65532:65532


# Set minimal environment
ENV PATH="/usr/local/bin:${PATH}"

# Health check is not applicable for stdio-based servers
# No HEALTHCHECK directive

# Run as non-root
ENTRYPOINT ["/usr/local/bin/fetch-mcp-server"]

# Default arguments can be overridden
CMD []