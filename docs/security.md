# Security

Features (opt-in):
- API key authentication: `CMP_AUTH_ENABLED=true` enables auth checks
- Rate limiting: token-bucket per API key/tenant/IP
- RBAC: require `chat:execute` for chat endpoint
- Audit logging: JSON events written to `audit.log`

Headers:
- `X-Tenant-ID`: optional tenant scoping for requests

Enable:
```bash
export CMP_AUTH_ENABLED=true
ctx serve --addr :8000
```
