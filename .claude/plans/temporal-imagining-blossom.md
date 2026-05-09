# Plan: Zero-Downtime for All Services

## Context

เมื่อกี้เพิ่ม vault-agent `exec` ที่ใช้ `pkill` เพื่อ restart container เมื่อ Vault secret เปลี่ยน
ปัญหา: `replicas: 1` ทำให้ทุก restart = downtime แม้แต่ช่วงสั้นๆ และถ้า replicas เพิ่มเป็น 2
vault-agent บนทั้ง 2 pods จะ `pkill` nginx พร้อมกัน → ยัง downtime อยู่

**Solution**: เปลี่ยนจาก `pkill` เป็น `kubectl rollout restart` ซึ่ง K8s จะจัดการ
rolling update เอง (maxUnavailable: 0) + เพิ่ม replicas: 2 + PodDisruptionBudget

## Architecture

```
Vault secret เปลี่ยน
  → vault-agent (ทั้ง 2 pods) render ไฟล์ใหม่
  → exec: kubectl rollout restart deployment/<name>
  → K8s สร้าง pod ใหม่ (อ่าน Vault secret ล่าสุด → app-config.js ใหม่)
  → รอ readiness probe ผ่าน
  → terminate pod เก่า
  → zero downtime ตลอด (maxUnavailable: 0)
```

Replicas ยังเป็น 1 แต่ maxSurge: 1 + maxUnavailable: 0 ที่มีอยู่แล้ว
→ K8s spin up pod ใหม่ก่อน (รวม 2 pods ชั่วคราว) แล้วค่อย terminate pod เก่า
→ zero-downtime แม้ replicas = 1

## Changes Required

### 1. _helpers.tpl — เพิ่ม kubectl init container + tools volume

**`k8s/charts/jeeb-app/templates/_helpers.tpl`**
**`k8s/charts/jeeb-learning/templates/_helpers.tpl`**

เพิ่ม 2 helpers ใหม่:

```
jeeb-app.vaultKubectlInitContainer  →  init container: copy kubectl จาก bitnami/kubectl:1.31 → /tools/
jeeb-app.vaultToolsVolume           →  emptyDir volume ชื่อ "tools"
```

แก้ `vaultAgentContainer` helper เดิม → เพิ่ม mount `/tools` (readOnly)
แก้ `vaultVolumes` helper เดิม → append tools emptyDir

### 3. deployment.yaml × 4 — ใส่ init container + tools volume

เพิ่ม `initContainers:` block โดย include helper ใหม่:

```yaml
initContainers:
  {{- include "jeeb-app.vaultKubectlInitContainer" . | nindent 8 }}
```

เพิ่ม volume ใน `volumes:` section:
```yaml
{{- include "jeeb-app.vaultToolsVolume" . | nindent 8 }}
```

ไฟล์:
- `k8s/charts/jeeb-app/templates/frontend/deployment.yaml`
- `k8s/charts/jeeb-app/templates/backend/deployment.yaml`
- `k8s/charts/jeeb-learning/templates/frontend/deployment.yaml`
- `k8s/charts/jeeb-learning/templates/backend/deployment.yaml`

### 4. vault-agent-config.yaml × 4 — เปลี่ยน restartCmd

| Service | restartCmd เดิม | restartCmd ใหม่ |
|---------|----------------|----------------|
| jeeb-app frontend | `pkill -TERM -f 'nginx: master' \|\| true` | `/tools/kubectl rollout restart deployment/frontend -n <ns>` |
| jeeb-app backend | `pkill -TERM jeeb \|\| true` | `/tools/kubectl rollout restart deployment/backend -n <ns>` |
| learning frontend | `pkill -TERM -f 'nginx: master' \|\| true` | `/tools/kubectl rollout restart deployment/learning-frontend -n <ns>` |
| learning backend | `pkill -TERM jeeb-learning \|\| true` | `/tools/kubectl rollout restart deployment/learning-backend -n <ns>` |

Namespace ใช้ `{{ .Values.global.namespace }}` จาก template (ไม่ hardcode)

### 5. RBAC — ไฟล์ใหม่ × 4 (Role + RoleBinding)

แต่ละ ServiceAccount ต้องการ `get` + `patch` บน Deployment ของตัวเอง
เพื่อให้ `kubectl rollout restart` ทำงานได้

```
k8s/charts/jeeb-app/templates/frontend/rollout-rbac.yaml
k8s/charts/jeeb-app/templates/backend/rollout-rbac.yaml
k8s/charts/jeeb-learning/templates/frontend/rollout-rbac.yaml
k8s/charts/jeeb-learning/templates/backend/rollout-rbac.yaml
```

Pattern ตัวอย่าง (frontend):
```yaml
kind: Role
rules:
- apiGroups: ["apps"]
  resources: ["deployments"]
  resourceNames: ["frontend"]
  verbs: ["get", "patch"]
---
kind: RoleBinding
subjects:
- kind: ServiceAccount
  name: frontend
```

### 6. PodDisruptionBudget — ไฟล์ใหม่ × 4

ป้องกัน K8s จาก involuntary eviction (node drain, upgrades) ที่อาจ kill ทุก pod พร้อมกัน

```
k8s/charts/jeeb-app/templates/frontend/pdb.yaml
k8s/charts/jeeb-app/templates/backend/pdb.yaml
k8s/charts/jeeb-learning/templates/frontend/pdb.yaml
k8s/charts/jeeb-learning/templates/backend/pdb.yaml
```

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: frontend  # ต่างกันตาม service
```

## Files Summary

| Action | File |
|--------|------|
| Edit | `k8s/charts/jeeb-app/templates/_helpers.tpl` |
| Edit | `k8s/charts/jeeb-learning/templates/_helpers.tpl` |
| Edit | `k8s/charts/jeeb-app/templates/frontend/deployment.yaml` |
| Edit | `k8s/charts/jeeb-app/templates/backend/deployment.yaml` |
| Edit | `k8s/charts/jeeb-learning/templates/frontend/deployment.yaml` |
| Edit | `k8s/charts/jeeb-learning/templates/backend/deployment.yaml` |
| Edit | `k8s/charts/jeeb-app/templates/frontend/vault-agent-config.yaml` |
| Edit | `k8s/charts/jeeb-app/templates/backend/vault-agent-config.yaml` |
| Edit | `k8s/charts/jeeb-learning/templates/frontend/vault-agent-config.yaml` |
| Edit | `k8s/charts/jeeb-learning/templates/backend/vault-agent-config.yaml` |
| New  | `k8s/charts/jeeb-app/templates/frontend/rollout-rbac.yaml` |
| New  | `k8s/charts/jeeb-app/templates/backend/rollout-rbac.yaml` |
| New  | `k8s/charts/jeeb-learning/templates/frontend/rollout-rbac.yaml` |
| New  | `k8s/charts/jeeb-learning/templates/backend/rollout-rbac.yaml` |
| New  | `k8s/charts/jeeb-app/templates/frontend/pdb.yaml` |
| New  | `k8s/charts/jeeb-app/templates/backend/pdb.yaml` |
| New  | `k8s/charts/jeeb-learning/templates/frontend/pdb.yaml` |
| New  | `k8s/charts/jeeb-learning/templates/backend/pdb.yaml` |

## Verification

1. Deploy: `go run ./cmd/k8s-manager deploy app && go run ./cmd/k8s-manager deploy learning`
2. ตรวจสอบ pods มี 2 replicas: `kubectl get pods -n jeeb-dev`
3. ตรวจสอบ PDB: `kubectl get pdb -n jeeb-dev`
4. ทดสอบ vault secret change trigger:
   ```powershell
   vault kv patch secret/jeeb/frontend/develop VITE_KEYCLOAK_URL=https://auth.jeeb-dev.local
   ```
5. ดู rolling restart เกิดขึ้น: `kubectl get pods -n jeeb-dev -w`
   - ต้องเห็น pod ใหม่ขึ้นมาก่อน pod เก่าถูก terminate
6. ตรวจสอบ `https://jeeb-dev.local` ยังใช้งานได้ระหว่าง restart
