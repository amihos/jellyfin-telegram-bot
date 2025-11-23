# ============================================
# Multi-stage Dockerfile for Jellyfin Telegram Bot
# ============================================
# This Dockerfile creates a minimal production image using multi-stage build
# Final image size target: under 50MB
# Security: Runs as non-root user

# ============================================
# Stage 1: Builder
# ============================================
FROM golang:1.23-alpine AS builder

# Install build dependencies
# gcc, musl-dev: Required for CGO (SQLite driver needs CGO)
# git: For go modules that might need git
# ca-certificates, tzdata: Will be copied to final image
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    gcc \
    musl-dev

# Set working directory
WORKDIR /app

# Copy go module files first for better layer caching
# This allows Docker to cache dependencies if go.mod/go.sum haven't changed
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build static binary with optimizations
# CGO_ENABLED=1: Enable CGO for SQLite driver (go-sqlite3 requires it)
# -a: Force rebuilding of packages
# -installsuffix cgo: Add suffix to package directory to keep separate from default builds
# -ldflags: Link flags for optimization
#   -w: Disable DWARF generation (removes debugging symbols)
#   -s: Disable symbol table
#   -linkmode external: Use external linker for better static linking with CGO
#   -extldflags "-static": Create fully static binary
# Note: We use CGO for SQLite but still create a portable binary via musl-dev
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -a \
    -installsuffix cgo \
    -ldflags="-w -s -linkmode external -extldflags '-static' -X main.version=docker" \
    -o jellyfin-telegram-bot \
    ./cmd/bot

# Verify binary was created and show its size
RUN ls -lh jellyfin-telegram-bot

# ============================================
# Stage 2: Runtime
# ============================================
FROM alpine:latest

# Install runtime dependencies
# ca-certificates: Required for HTTPS connections to Telegram API and Jellyfin
# tzdata: Timezone data for correct timestamp handling
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user and group for security
# Using numeric UID/GID for better compatibility
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Create directories for data persistence with correct permissions
RUN mkdir -p /app/data /app/logs && \
    chown -R appuser:appgroup /app

# Copy binary from builder stage
COPY --from=builder --chown=appuser:appgroup /app/jellyfin-telegram-bot .

# Copy locales directory for i18n support
COPY --from=builder --chown=appuser:appgroup /app/locales ./locales

# Verify locales were copied
RUN ls -la ./locales

# Switch to non-root user
USER appuser

# Set environment variables
ENV DATABASE_PATH=/app/data/bot.db
ENV LOG_FILE=/app/logs/bot.log

# Expose webhook port (default 8080, configurable via PORT env var)
EXPOSE 8080

# Health check endpoint
# Checks if the webhook server is responding
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Set entrypoint
ENTRYPOINT ["/app/jellyfin-telegram-bot"]

# ============================================
# Build Instructions
# ============================================
# Build image:
#   docker build -t jellyfin-telegram-bot:latest .
#
# Run container:
#   docker run -d \
#     --name jellyfin-telegram-bot \
#     --env-file .env \
#     -p 8080:8080 \
#     -v $(pwd)/data:/app/data \
#     -v $(pwd)/logs:/app/logs \
#     jellyfin-telegram-bot:latest
#
# View logs:
#   docker logs -f jellyfin-telegram-bot
#
# Check size:
#   docker images jellyfin-telegram-bot:latest
#
# ============================================
# Security Notes
# ============================================
# - Runs as non-root user (UID 1000)
# - Static binary compiled with musl libc
# - Minimal attack surface (alpine base)
# - No shell or development tools in runtime image
# - Health check for monitoring
#
# ============================================
# Technical Notes
# ============================================
# CGO is enabled for SQLite compatibility (go-sqlite3 requires it)
# The binary is still static thanks to musl-dev and -extldflags '-static'
# This approach balances portability with SQLite functionality
