# Kubernetes Troubleshooting

## Pod stuck in Pending

**Cause:** Insufficient resources or unschedulable node.

**Solution:**
```bash
kubectl describe pod -n jeeb <pod-name>
# Look at Events section for reason

# Check node resources
kubectl describe nodes
```

---

## Pod CrashLoopBackOff

**Cause:** App crashes on startup — missing env vars, bad config, or port conflict.

**Solution:**
```bash
kubectl logs -n jeeb <pod-name> --previous
kubectl describe pod -n jeeb <pod-name>
```

---

## ImagePullBackOff

**Cause:** Cannot pull image from Nexus registry.

**Solution:**
```bash
# Verify image exists in Nexus
curl http://localhost:30083

# Check the docker secret is applied
kubectl get secret nexus-docker-secret -n jeeb

# Recreate the secret
kubectl create secret docker-registry nexus-docker-secret \
  --docker-server=localhost:30050 \
  --docker-username=admin \
  --docker-password=<password> \
  -n jeeb
```

---

## Service not reachable

**Cause:** Pod not ready, wrong NodePort, or service selector mismatch.

**Solution:**
```bash
# Check service and endpoints
kubectl get svc -n jeeb
kubectl get endpoints -n jeeb

# Verify pod labels match service selector
kubectl describe svc backend -n jeeb
kubectl get pods -n jeeb --show-labels
```

---

## PVC stuck in Pending

**Cause:** No StorageClass available or insufficient storage.

**Solution:**
```bash
kubectl get pvc -n jeeb
kubectl describe pvc <pvc-name> -n jeeb

# Docker Desktop uses hostpath — ensure default StorageClass exists
kubectl get storageclass
```

---

## kubectl: connection refused

**Cause:** Docker Desktop Kubernetes not running.

**Solution:**
- Open Docker Desktop → Settings → Kubernetes → Enable Kubernetes
- Wait for the green indicator
```bash
kubectl cluster-info
```

---

## Out of disk space

```bash
# Remove unused Docker images
docker system prune -a

# Remove unused volumes
docker volume prune
```
