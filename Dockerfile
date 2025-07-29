# ===== BUILD STAGE =====
# Use lightweight Alpine-based Go image for building
FROM golang:1.24.3-alpine AS builder

# Set working directory inside container
WORKDIR /app

# Install system dependencies:
# - git: Required for Go module dependencies
# - ca-certificates: For SSL certificate verification
RUN apk add --no-cache git ca-certificates

# Copy dependency definition files first (optimizes Docker layer caching)
COPY go.mod go.sum ./

# Download all dependencies (cached unless go.mod/go.sum change)
RUN go mod download

# Copy application source code into container
COPY . .

# Compile statically linked binary with optimizations:
# - CGO_ENABLED=0: Disables CGO for fully static binary
# - GOOS=linux: Targets Linux OS
# - ldflags="-w -s": Strips debug information (reduces binary size)
# - ./cmd/main.go: Entry point of the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/main.go

# ===== RUNTIME STAGE =====
# Use minimal Alpine base image for runtime
FROM alpine:3.19

# Set working directory for application
WORKDIR /app

# Create non-root user for security (reduces attack surface)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy artifacts from builder stage:
# 1. Compiled application binary
COPY --from=builder /app/main .
# 2. Environment configuration file
COPY --from=builder /app/.env .
# 3. Database migration scripts
COPY --from=builder /app/migrations ./migrations
# 4. SSL certificates (required for HTTPS/API communication)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Ensure proper file ownership for security
RUN chown -R appuser:appgroup /app

# Switch to non-root user context
USER appuser

# ===== METADATA =====
# Add informational labels (useful for maintenance and auditing)
LABEL maintainer="georgekmutonya@gmail.com"
LABEL version="1.0"
LABEL description="Savannah Go Backend"
# Link to source code repository (OCI image standard)
LABEL org.opencontainers.image.source="https://github.com/Mutonya/savannah-go-backend"

# ===== HEALTH MONITORING =====
# Configure container health checks:
# --interval=30s: Check every 30 seconds
# --timeout=3s: Fail if no response in 3 seconds
# Command: Verify /health endpoint availability silently
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Expose application port (informational only - doesn't publish ports)
EXPOSE 8080

# ===== RUNTIME COMMAND =====
# Entrypoint command to execute the application
CMD ["./main"]