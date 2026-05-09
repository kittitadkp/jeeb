# New Cluster Setup

## Goal

Bring up the full local Jeeb stack on Docker Desktop Kubernetes.

## Steps

1. Enable Kubernetes in Docker Desktop.
2. Prepare secrets:

```powershell
cd k8s-manager
Copy-Item env/secrets.yaml.example env/secrets.yaml
go run ./cmd/k8s-manager validate
```

3. Run bootstrap:

```powershell
go run ./cmd/k8s-manager setup
```

4. Open Jenkins at `http://localhost:30082` and run the four service pipelines.
5. Deploy app workloads:

```powershell
go run ./cmd/k8s-manager deploy app learning
```

6. Verify health:

```powershell
go run ./cmd/k8s-manager check
```

## Hosts and DNS

The stack expects `.local` hostnames such as `jeeb-dev.local` and `auth.jeeb-dev.local`. `setup` patches CoreDNS, but the patch hardcodes the ingress controller ClusterIP. Re-run or refresh it after a cluster reset if those names stop resolving.

## First troubleshooting checks

- `kubectl get pods -A`
- `go run ./cmd/k8s-manager check`
- `go run ./cmd/k8s-manager maintain`
