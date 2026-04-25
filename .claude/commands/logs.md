---
description: Tail logs for a jeeb service — usage: /logs <service>
---

Show the last 50 lines of logs for the service: $ARGUMENTS

```bash
kubectl logs -n jeeb deployment/$ARGUMENTS --tail=50
```

If the pod is crashing, also run:
```bash
kubectl logs -n jeeb deployment/$ARGUMENTS --previous --tail=50
```

Summarize any errors or warnings found.
