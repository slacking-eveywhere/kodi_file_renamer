# Multi-stage Dockerfile for kodi-renamer
# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies (if any exist)
RUN go mod download || true

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build the application
# -ldflags="-s -w" strips debug information for smaller binary
# -trimpath removes file system paths for reproducible builds
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /app/kodi-renamer \
    ./cmd/kodi-renamer

# Verify the binary was created and is executable
RUN test -f /app/kodi-renamer && chmod +x /app/kodi-renamer

# Stage 2: Create minimal runtime image using distroless
FROM gcr.io/distroless/static-debian12:nonroot

# Add labels for metadata
LABEL org.opencontainers.image.title="kodi-renamer"
LABEL org.opencontainers.image.description="Command-line tool to rename movie and TV series files for Kodi"
LABEL org.opencontainers.image.authors="kodi-renamer contributors"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/yourusername/kodi-renamer"

ENV TVDB_API_KEY=
ENV TMDB_API_KEY=
ENV MOVIE_TO_RENAME_DIR=
ENV MOVIE_RENAMED_DIR=
ENV SERIE_TO_RENAME_DIR=
ENV SERIE_RENAMED_DIR=

# Copy the binary from builder stage
COPY --from=builder /app/kodi-renamer /usr/local/bin/kodi-renamer

# Set the working directory for media files
WORKDIR /media

# Use non-root user (distroless nonroot user UID 65532)
USER nonroot:nonroot

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/kodi-renamer"]
