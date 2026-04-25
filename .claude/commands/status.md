---
description: Show status of all pods and services in the jeeb namespace
---

Run the following and summarize the health of all services:

```bash
kubectl get pods -n jeeb
kubectl get svc -n jeeb
```

Report which pods are Running/Pending/CrashLoopBackOff and which services are reachable. Flag anything that needs attention.
