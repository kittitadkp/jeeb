---
name: Jeeb infra audit baseline
description: Full k8s audit completed 2026-05-03; known findings, accepted risks, and security posture baseline
type: project
---

Comprehensive audit of all manifests in k8s/ completed on 2026-05-03.

**Why:** Establishing baseline for ongoing security posture tracking across jeeb-dev, jeeb-infra, jeeb-obs namespaces.

**How to apply:** Use as reference for future audits to identify new regressions vs. known-accepted risks.

## Key findings at audit time

### Critical / accepted risks (personal dev project)
- Jenkins ServiceAccount bound to `cluster-admin` ClusterRoleBinding — intentional for local dev CI/CD (kubectl set image etc.), but should be scoped down if moving toward production.
- Vault TLS disabled (`tls_disable = true`) — acceptable on Docker Desktop localhost but blocks any remote access security.
- All Vault traffic in-cluster is HTTP, not HTTPS.
- MongoDB password `jeeb123` and Keycloak password `admin123` committed in `k8s/app/secrets.yaml` and `k8s/charts/jeeb-app/values.yaml` — known risk, dev-only credentials.
- Nexus admin credential `K@ng_12092540` is encoded (not encrypted) in two dockerconfigjson secrets committed to git. Same password reused for Jenkins admin and Nexus PAT.
- SonarQube token `squ_4cfe749941a6aea97dd29a3679dbd11caeae139b` committed plaintext in `k8s/charts/jeeb-infra/values.yaml`.
- Grafana admin password `admin123` in `k8s/charts/jeeb-obs/values.yaml`.

### Architecture decisions
- Vault agent sidecar pattern used for all app services (backend, frontend, learning) — deliberate, good pattern.
- `IPC_LOCK` capability added to Vault and vault-agent containers — required to prevent secret swap to disk, correct.
- Old raw manifests in `k8s/app/`, `k8s/jenkins/`, `k8s/nexus/`, `k8s/sonarqube/`, `k8s/vault/` all target the obsolete `jeeb` namespace (not jeeb-dev/jeeb-infra). These appear superseded by Helm charts.
- Helm chart is the authoritative deployment path via `apply.sh`.

### Security contexts missing (broad)
- All deployments lack `runAsNonRoot: true` and `seccompProfile` at pod level.
- MongoDB StatefulSet runs as root (no securityContext at all).
- Keycloak Deployment has no securityContext.
- Nexus Deployment has no securityContext.
- SonarQube uses privileged init container (required for vm.max_map_count).
- Promtail DaemonSet has no securityContext (mounts host paths).

### Reliability gaps
- No livenessProbe on Jenkins (only readiness).
- No livenessProbe on Nexus.
- No livenessProbe on SonarQube.
- No livenessProbe on Grafana (only readiness).
- No livenessProbe on Loki (only readiness).
- No livenessProbe on Tempo (only readiness).
- MongoDB StatefulSet has no liveness or readiness probes at all.
- No NetworkPolicies anywhere in the cluster.
- No PodDisruptionBudgets (acceptable for single-node dev).

### Image tagging
- `sonatype/nexus3:latest` used in both raw manifest and Helm chart — unpinned.
- `busybox` unpinned in SonarQube init container.
- `mongo:7` is major-tag-only (not digest pinned) — acceptable for dev.
- App images use `:latest` pulled from local Nexus — acceptable for dev CI.

### Configuration bugs
- Raw manifests in `k8s/app/` and `k8s/jenkins/` etc. target `namespace: jeeb` which no longer exists (split into jeeb-dev / jeeb-infra). These manifests are stale/broken if applied directly.
- `k8s/app/backend/configmap.yaml` contains `MONGO_URI` with plaintext credentials — superseded by Vault in Helm chart but left in place.
- `k8s/00-namespace.yaml` creates `jeeb` namespace — obsolete.
- Vault VAULT_ADDR in old `k8s/vault/statefulset.yaml` still points to `vault.jeeb.svc.cluster.local` instead of `vault.jeeb-infra.svc.cluster.local`.

### Jenkins CASC security
- `useScriptSecurity: false` for Job DSL — allows arbitrary Groovy execution in Jenkins without sandbox. Acceptable for local but risky if Jenkins ever exposed.
- CSP header disabled via JAVA_OPTS: `-Dhudson.model.DirectoryBrowserSupport.CSP=` — enables XSS vectors in build artifacts.
