// Package server exposes the HTTP runtime for the Contexis CMP framework.
//
// It wires health/readiness/version/metrics endpoints and the chat API
// (POST /api/v1/chat), integrating context resolution, memory search,
// prompt rendering, optional model inference, and optional security controls
// such as authentication, RBAC, rate limiting, prompt-injection and PII policies.
package server
