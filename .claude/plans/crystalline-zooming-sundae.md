# Plan: Helm + Multi-Environment (dev / UAT) Refactor

## Context
The current k8s folder uses raw manifests with hardcoded values (registry URLs, Keycloak hostnames, Vault paths, credentials). Adding a UAT environment requires duplicating everything and editing by hand. This plan migrates to Helm, splits infra from app concerns, and introduces env-specific ingress so `jeeb-dev.local` and `jeeb-uat.local` work out of the box.

---

## Target Layout

```
k8s/
  charts/
    jeeb-infra/          # Vault + Jenkins + Nexus + SonarQube (shared, namespace jeeb-infra)
      Chart.yaml
      values.yaml
      templates/
        namespace.yaml
        vault/            (statefulset, service, configmap, pvc, rbac, serviceaccount)
        jenkins/          (deployment, service, pvc, rbac)
        nexus/            (deployment, service, pvc, secret)
        sonarqube/        (deployment, service, pvc)
        ingress.yaml      (jenkins.jeeb.local, nexus.jeeb.local, sonarqube.jeeb.local, vault.jeeb.local)
    jeeb-app/            # MongoDB + Keycloak + Backend + Frontend (per environment)
      Chart.yaml
      values.yaml         (shared defaults)
      values-dev.yaml     (dev overrides)
      values-uat.yaml     (UAT overrides)
      templates/
        namespace.yaml
        secrets.yaml
        mongodb/          (statefulset, service)
        keycloak/         (deployment, service, pvc, configmap for realm)
        backend/          (deployment, service, configmap, serviceaccount, vault-agent-config)
        frontend/         (deployment, service, serviceaccount, vault-agent-config)
        ingress.yaml      (www / api / auth per env)
  nginx-ingress/          # nginx ingress controller (one-time bootstrap)
    install.sh
  apply.sh                (updated: helm upgrade --install for both charts)
  apply-dev.sh            (shortcut: helm upgrade jeeb-app with values-dev.yaml)
  apply-uat.sh            (shortcut: helm upgrade jeeb-app with values-uat.yaml)
```

---

## Namespaces

| Namespace | Contents |
|-----------|----------|
| `jeeb-infra` | Vault, Jenkins, Nexus, SonarQube |
| `jeeb-dev` | MongoDB, Keycloak, Backend, Frontend (dev) |
| `jeeb-uat` | MongoDB, Keycloak, Backend, Frontend (UAT) |

---

## Ingress Hostnames

### App (env-specific)
| Host | Routes to |
|------|-----------|
| `jeeb-dev.local` | frontend service in jeeb-dev |
| `api.jeeb-dev.local` | backend service in jeeb-dev |
| `auth.jeeb-dev.local` | keycloak service in jeeb-dev |
| `jeeb-uat.local` | frontend service in jeeb-uat |
| `api.jeeb-uat.local` | backend service in jeeb-uat |
| `auth.jeeb-uat.local` | keycloak service in jeeb-uat |

### Infra (shared)
| Host | Routes to |
|------|-----------|
| `jenkins.jeeb.local` | jenkins in jeeb-infra |
| `nexus.jeeb.local` | nexus (UI port 8081) in jeeb-infra |
| `sonarqube.jeeb.local` | sonarqube in jeeb-infra |
| `vault.jeeb.local` | vault in jeeb-infra |

NodePorts on infra services are **kept** (backward compat during migration). Ingress is additive.

---

## Key Values That Differ dev vs UAT

```yaml
# values-dev.yaml
global:
  env: dev
  namespace: jeeb-dev
  domain: jeeb-dev.local
  imageTag: latest

keycloak:
  hostnameUrl: http://auth.jeeb-dev.local

backend:
  vault:
    path: secret/data/jeeb/backend/develop

frontend:
  vault:
    path: secret/data/jeeb/frontend/develop
```

```yaml
# values-uat.yaml
global:
  env: uat
  namespace: jeeb-uat
  domain: jeeb-uat.local
  imageTag: "1.0.0"           # pinned tag for UAT

keycloak:
  hostnameUrl: http://auth.jeeb-uat.local

backend:
  vault:
    path: secret/data/jeeb/backend/uat

frontend:
  vault:
    path: secret/data/jeeb/frontend/uat
```

---

## Steps

- [x] **Step 1 — Install nginx ingress controller**
  ```bash
  helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
  helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
    --namespace ingress-nginx --create-namespace \
    --set controller.service.type=LoadBalancer
  ```
  Verify: `kubectl get svc -n ingress-nginx` shows EXTERNAL-IP as `localhost`.

- [x] **Step 2 — Scaffold jeeb-infra chart**
  - Create `k8s/charts/jeeb-infra/Chart.yaml` (name: jeeb-infra, version: 0.1.0)
  - Create `k8s/charts/jeeb-infra/values.yaml` with all current hardcoded values externalised
  - Move and templatise each infra manifest into `templates/` subdirs
  - Add `templates/ingress.yaml` with four infra hosts

- [x] **Step 3 — Scaffold jeeb-app chart**
  - Create `k8s/charts/jeeb-app/Chart.yaml` (name: jeeb-app, version: 0.1.0)
  - Create `k8s/charts/jeeb-app/values.yaml` with all shared defaults
  - Move and templatise app manifests into `templates/` subdirs
  - Replace hardcoded `jeeb` namespace with `{{ .Values.global.namespace }}`
  - Replace `host.docker.internal:30081` keycloak hostname with `{{ .Values.keycloak.hostnameUrl }}`
  - Replace Vault path `secret/data/jeeb/backend/develop` with `{{ .Values.backend.vault.path }}`
  - Replace Vault path `secret/data/jeeb/frontend/develop` with `{{ .Values.frontend.vault.path }}`
  - Add `templates/ingress.yaml` with frontend/api/auth hosts using `{{ .Values.global.domain }}`

- [x] **Step 4 — Create values-dev.yaml and values-uat.yaml**
  - Only override what differs (env, namespace, domain, imageTag, keycloak URL, vault paths)

- [x] **Step 5 — Update apply.sh + add apply-dev.sh / apply-uat.sh**
  ```bash
  # apply.sh
  helm upgrade --install jeeb-infra k8s/charts/jeeb-infra --namespace jeeb-infra --create-namespace
  helm upgrade --install jeeb-dev   k8s/charts/jeeb-app   --namespace jeeb-dev   --create-namespace -f k8s/charts/jeeb-app/values-dev.yaml
  helm upgrade --install jeeb-uat   k8s/charts/jeeb-app   --namespace jeeb-uat   --create-namespace -f k8s/charts/jeeb-app/values-uat.yaml
  ```

- [x] **Step 6 — Windows hosts file** (manual, one-time)
  Add to `C:\Windows\System32\drivers\etc\hosts`:
  ```
  127.0.0.1 jeeb-dev.local api.jeeb-dev.local auth.jeeb-dev.local
  127.0.0.1 jeeb-uat.local api.jeeb-uat.local auth.jeeb-uat.local
  127.0.0.1 jenkins.jeeb.local nexus.jeeb.local sonarqube.jeeb.local vault.jeeb.local
  ```

- [x] **Step 7 — Delete old raw-manifest resources and apply Helm charts**
  ```bash
  kubectl delete namespace jeeb          # removes old single-env resources
  bash k8s/apply.sh                      # deploys everything via Helm
  ```

- [x] **Step 8 — Verify**
  ```bash
  kubectl get pods -n jeeb-infra
  kubectl get pods -n jeeb-dev
  kubectl get pods -n jeeb-uat
  kubectl get ingress -A
  # Then open in browser:
  # http://jeeb-dev.local
  # http://jeeb-uat.local
  # http://jenkins.jeeb.local
  ```

---

## Jenkins Pipeline Impact

The Jenkinsfiles currently end with:
```groovy
sh "kubectl set image deployment/backend backend=localhost:30050/jeeb/backend:${BUILD_NUMBER} -n jeeb"
```

After this migration that command must change to:
```groovy
// backend/Jenkinsfile — deploy to dev on every merge
sh """
  helm upgrade jeeb-dev k8s/charts/jeeb-app \
    --namespace jeeb-dev \
    --reuse-values \
    --set global.imageTag=${BUILD_NUMBER}
"""
```

For UAT promotion (manual trigger or separate pipeline):
```groovy
sh """
  helm upgrade jeeb-uat k8s/charts/jeeb-app \
    --namespace jeeb-uat \
    -f k8s/charts/jeeb-app/values-uat.yaml \
    --set global.imageTag=${params.IMAGE_TAG}
"""
```

- [x] **Step 9 — Update jenkins/backend/Jenkinsfile** — replace `kubectl set image` with `helm upgrade --reuse-values --set global.imageTag=...` targeting `jeeb-dev`
- [x] **Step 10 — Update jenkins/frontend/Jenkinsfile** — same pattern for frontend

> **Image tag strategy**: dev pipeline uses `${BUILD_NUMBER}` (auto-deploy on every build); UAT uses a manually promoted pinned tag via a parameterised pipeline or manual `helm upgrade`.

---

## Files to Create/Modify

**Create (new):**
- `k8s/charts/jeeb-infra/Chart.yaml`
- `k8s/charts/jeeb-infra/values.yaml`
- `k8s/charts/jeeb-infra/templates/` (all infra templates)
- `k8s/charts/jeeb-app/Chart.yaml`
- `k8s/charts/jeeb-app/values.yaml`
- `k8s/charts/jeeb-app/values-dev.yaml`
- `k8s/charts/jeeb-app/values-uat.yaml`
- `k8s/charts/jeeb-app/templates/` (all app templates)
- `k8s/nginx-ingress/install.sh`
- `k8s/apply-dev.sh`
- `k8s/apply-uat.sh`

**Modify:**
- `k8s/apply.sh` — switch from kubectl apply to helm upgrade

**Delete (superseded by Helm templates):**
- `k8s/00-namespace.yaml` and all files under `k8s/app/`, `k8s/jenkins/`, `k8s/nexus/`, `k8s/sonarqube/`, `k8s/vault/`

---

## Verification

1. `helm list -A` — shows three releases: jeeb-infra, jeeb-dev, jeeb-uat
2. `kubectl get pods -n jeeb-dev` and `kubectl get pods -n jeeb-uat` — all Running
3. `kubectl get ingress -A` — all six app hosts + four infra hosts listed
4. Browser: `http://jeeb-dev.local` loads frontend; `http://jeeb-uat.local` loads frontend (different image tag)
5. `helm upgrade jeeb-uat k8s/charts/jeeb-app -f values-uat.yaml --set global.imageTag=1.1.0` — rolling update applies without touching dev
