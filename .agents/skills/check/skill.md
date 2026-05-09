# /check skill

Run `k8s-manager check` (which executes `scripts/health_check.py`) and format the output.

## Trigger phrases
- `/check`
- "check cluster"
- "health check"
- "is cluster ok"
- "cluster status"

## Behavior

1. Run the health check:
   ```
   k8s-manager check
   ```
   Or directly:
   ```
   python3 k8s-manager/scripts/health_check.py
   ```

2. Parse and display results. For each `[FAIL]` line, immediately mention it prominently.

3. If failures exist, suggest:
   ```
   k8s-manager maintain
   ```
   to get exact fix commands.

## Output format

Present results as a brief table. Lead with failures if any. Example:

```
PASS  Pod health        all 14 pods healthy
FAIL  Vault             http://localhost:30091/v1/sys/health — Vault is sealed
PASS  Keycloak          http://localhost:30081/realms/jeeb → 200
PASS  Kong              http://localhost:30088/health → 200
PASS  Jenkins           http://localhost:30082/login → 200
FAIL  DNS: auth         not resolved inside cluster
PASS  Vault secrets     backend secrets readable

2 checks failed. Run `k8s-manager maintain` for diagnosis.
```

If all pass, one line: "All checks passed — cluster is healthy."
