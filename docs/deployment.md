# Deployment

Build:
```bash
ctx build --image contexis-cmp/contexis --tag latest
```

Run locally with Docker:
```bash
ctx deploy --target docker --image contexis-cmp/contexis:latest --ports 8000:8000 --detach
```

Kubernetes (Helm chart available under `charts/contexis-app/`):
- Prepare image registry and set values
- Apply manifests or use Helm to install into your namespace

Argo Rollouts manifests are under `deploy/`.
