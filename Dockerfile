# =========================================================
# Dockerfile — Tu Tien Discord Bot v0.1
# Multi-stage build: builder + minimal runtime image.
# =========================================================

# --- Stage 1: Build ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git for private module support (if needed in future)
RUN apk add --no-cache git ca-certificates

# Cache Go module downloads separately from source
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/bot ./cmd/bot

# --- Stage 2: Runtime ---
FROM alpine:3.19

# ca-certificates needed for TLS connections to MongoDB Atlas and Discord
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/bot /app/bot

# Security: run as non-root user
RUN addgroup -S botgroup && adduser -S botuser -G botgroup
USER botuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/bot"]
