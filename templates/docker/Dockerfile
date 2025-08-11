# Template Dockerfile for generated CMP apps

FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY src ./src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /out/ctx ./src/cli/main.go

FROM python:3.11-slim
ENV CMP_ENV=production PORT=8000
WORKDIR /app
COPY --from=builder /out/ctx /usr/local/bin/ctx
COPY contexts ./contexts
COPY prompts ./prompts
COPY tools ./tools
COPY config ./config
USER 10001
EXPOSE 8000
ENTRYPOINT ["ctx"]
CMD ["serve", "--addr", ":8000"]


