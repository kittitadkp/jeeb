---
description: Deploy or redeploy a jeeb service — usage: /deploy <service>
---

Redeploy the service: $ARGUMENTS

Steps:
1. Check current pod status: `kubectl get pods -n jeeb -l app=$ARGUMENTS`
2. Restart the deployment: `kubectl rollout restart deployment/$ARGUMENTS -n jeeb`
3. Watch rollout: `kubectl rollout status deployment/$ARGUMENTS -n jeeb --timeout=120s`
4. Show final pod status

If rollout fails, fetch logs and describe the pod to diagnose the issue.
