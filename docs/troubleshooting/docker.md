# Docker and Cluster Troubleshooting

## Docker Desktop cluster reset

After a reset, expect to rerun:

```powershell
cd k8s-manager
go run ./cmd/k8s-manager setup
```

`coredns-patch.yaml` may also need to be refreshed because it contains a concrete ingress controller ClusterIP.

## Images do not pull

- Verify Nexus is reachable on `localhost:30050` from the cluster.
- Re-run Jenkins pipelines so the expected tags exist.
- Check `nexus-pull-secret` in the target namespace.

## Old docs mention Docker Compose

Ignore them. The current repository does not include a maintained root `docker-compose.yml`.
