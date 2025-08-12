# Contexis CMP Framework - Production Dockerfile (Week 9: Deployment System)

# --- Builder stage: build the Go CLI binary ---
FROM golang:1.24-alpine AS builder

WORKDIR /src

# Ensure the Go toolchain can auto-manage minor versions if needed
ENV GOTOOLCHAIN=auto

# Install git and build deps
RUN apk update && apk upgrade --no-cache \
    && apk add --no-cache git build-base

# Cache go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY src ./src

# Build CLI
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /out/ctx ./src/cli/main.go


######## Distroless runtime stage (minimal attack surface) ########
FROM cgr.dev/chainguard/static:latest AS runtime

ENV CMP_ENV=production \
    PORT=8000

WORKDIR /app

# Copy CLI binary to /app/ctx and use absolute entrypoint
COPY --from=builder /out/ctx /app/ctx

# Copy runtime assets (read-only rootfs expected in k8s); ensure presence
COPY contexts ./contexts
COPY prompts ./prompts
COPY tools ./tools
COPY config ./config

# Run as non-root user (Chainguard static uses non-root by default, but set explicitly)
USER 65532:65532

EXPOSE 8000

ENTRYPOINT ["/app/ctx"]
CMD ["serve", "--addr", ":8000"]

LABEL org.opencontainers.image.title="Contexis CMP Framework" \
      org.opencontainers.image.description="Rails-inspired framework for reproducible AI apps (CMP)." \
      org.opencontainers.image.source="https://github.com/contexis-cmp/contexis" \
      org.opencontainers.image.licenses="MIT"


