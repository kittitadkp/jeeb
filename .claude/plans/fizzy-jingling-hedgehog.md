# Plan: Fix Keycloak Setup via k8s-manager

## Context

เจอ 3 ปัญหาต่อเนื่องกัน:
1. **Redirect URI invalid** — `realm-jeeb.json` ไม่มี `https://jeeb-dev.local/*` → แก้ไขแล้ว
2. **Admin password ว่างในcluster** — `keycloak-secret` ถูกสร้างโดย `helm upgrade` โดยไม่ได้ pass `values-secrets.yaml` เลยได้ค่า default (ว่าง) แทนที่จะอ่านจาก `secrets.yaml`
3. **PVC stuck Terminating** — ทุกครั้งที่ delete PVC ต้อง patch finalizer ด้วยตัวเอง

**Root cause:** ไม่เคยรัน `k8s-manager deploy data` ซึ่งเป็น command ที่ถูกต้องในการ generate `values-secrets.yaml` จาก `secrets.yaml` แล้ว pass ให้ Helm สร้าง secret ให้ถูกต้อง

## Critical Files

- `k8s-manager/internal/setup/secrets_values.go` — generate values-secrets.yaml
- `k8s-manager/internal/setup/setup.go` — deployData() passes `-f values-secrets.yaml`
- `k8s/charts/jeeb-data/templates/secrets.yaml` — Helm template สร้าง keycloak-secret
- `k8s/charts/jeeb-data/files/realm-jeeb.json` — ✅ แก้ redirect URIs แล้ว

## Steps

- [ ] 1. **แก้ Keycloak crash ระหว่าง realm import** (พบ error: "Database is already closed")
       สาเหตุ: JVM shutdown ระหว่าง import → อาจ OOMKill หรือ H2 race condition
       ตรวจสอบ:
       ```
       kubectl describe pod -n jeeb-dev -l app=keycloak | grep -A5 "OOMKilled\|Reason\|Exit Code"
       kubectl top pod -n jeeb-dev -l app=keycloak
       ```
       ถ้า OOMKilled → เพิ่ม memory limit ใน deployment (ปัจจุบัน limit: 1Gi)
       ถ้าไม่ใช่ → รอ pod restart เองแล้วดู logs ใหม่

- [ ] 2. **รัน k8s-manager deploy data** เพื่อ sync secrets อย่างถูกต้อง
       ```
       cd k8s-manager
       go run ./cmd/k8s-manager deploy data
       ```
       command นี้จะ:
       - อ่าน `env/secrets.yaml` → generate `values-secrets.yaml`
       - รัน `helm upgrade --install jeeb-data -f values-dev.yaml -f values-secrets.yaml`
       - สร้าง/อัปเดต `keycloak-secret` ด้วย password ที่ถูกต้อง
       - อัปเดต `keycloak-realm` ConfigMap (realm-jeeb.json)

- [ ] 3. **ถ้า PVC stuck อีก** patch finalizer แล้ว recreate:
       ```
       kubectl patch pvc keycloak-pvc -n jeeb-dev -p '{"metadata":{"finalizers":[]}}' --type=merge
       kubectl apply -f - <<EOF
       apiVersion: v1
       kind: PersistentVolumeClaim
       metadata:
         name: keycloak-pvc
         namespace: jeeb-dev
       spec:
         accessModes: [ReadWriteOnce]
         resources:
           requests:
             storage: 1Gi
       EOF
       kubectl rollout restart deployment/keycloak -n jeeb-dev
       ```

- [ ] 4. **สร้าง Keycloak user** หลัง Keycloak พร้อม
       ไปที่ `https://auth.jeeb-dev.local/admin/` หรือ `http://localhost:30081/admin/`
       - Login: `admin` / `K@ng_12092540`
       - สร้าง user ใน realm `jeeb`

- [ ] 5. **ทดสอบ login** ที่ `https://jeeb-dev.local/`
       ตรวจสอบว่า redirect ไป Keycloak และ login สำเร็จ

## Verification

```
# Keycloak ready
kubectl rollout status deployment/keycloak -n jeeb-dev

# Secret มี password ถูกต้อง
kubectl get secret keycloak-secret -n jeeb-dev -o jsonpath='{.data.admin-password}' | base64 -d

# Login ที่ https://jeeb-dev.local/ ไม่ขึ้น redirect_uri error
curl -Lk -o /dev/null -w "%{http_code}\n" https://jeeb-dev.local/
```
