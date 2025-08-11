# Contexis CMP Framework - Production Dockerfile (Week 9: Deployment System)

# --- Builder stage: build the Go CLI binary ---
FROM golang:1.21-alpine AS builder

WORKDIR /src

# Install git and build deps
RUN apk add --no-cache git build-base

# Cache go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY src ./src

# Build CLI
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /out/ctx ./src/cli/main.go


# --- Runtime stage: minimal Python runtime + ctx binary ---
FROM python:3.11-slim AS runtime

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    CMP_ENV=production \
    PORT=8000

WORKDIR /app

# System deps (optional, light)
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*

# Copy CLI binary
COPY --from=builder /out/ctx /usr/local/bin/ctx

# Copy minimal runtime assets (contexts, prompts, tools, config) if present
COPY contexts ./contexts
COPY prompts ./prompts
COPY tools ./tools
COPY config ./config

# Create non-root user
RUN useradd -m -u 10001 appuser \
    && chown -R appuser:appuser /app

USER appuser

EXPOSE 8000

ENTRYPOINT ["ctx"]
CMD ["serve", "--addr", ":8000"]

LABEL org.opencontainers.image.title="Contexis CMP Framework" \
      org.opencontainers.image.description="Rails-inspired framework for reproducible AI apps (CMP)." \
      org.opencontainers.image.source="https://github.com/contexis-cmp/contexis" \
      org.opencontainers.image.licenses="MIT"


