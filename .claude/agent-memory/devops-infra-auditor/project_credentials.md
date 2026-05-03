---
name: Credential baseline
description: Known hardcoded credentials in the repo, their locations, and risk acceptance status as of 2026-05-03
type: project
---

Credentials found committed to git as of the 2026-05-03 audit. All are dev-only, personal project.

**Why:** Track so future audits can detect new credential leaks vs. known-accepted ones.

**How to apply:** Flag any NEW credentials discovered that are not in this list. Credentials here are known-accepted risks for local Docker Desktop dev.

| Credential | Value (masked) | Location | Risk Status |
|---|---|---|---|
| MongoDB password | jeeb123 | k8s/app/secrets.yaml:9, k8s/charts/jeeb-app/values.yaml:18 | Accepted (dev) |
| MongoDB URI | full URI with creds | k8s/app/backend/configmap.yaml:12 | Accepted (superseded by Vault) |
| Keycloak admin password | admin123 | k8s/app/secrets.yaml:20, k8s/charts/jeeb-app/values.yaml:28 | Accepted (dev) |
| Nexus admin password | K@ng_12092540 | k8s/charts/jeeb-infra/values.yaml:34, also encoded in dockerconfigjson | Accepted (dev) |
| Jenkins admin password | K@ng_12092540 | k8s/charts/jeeb-infra/values.yaml:34 | Accepted (dev), same as Nexus - password reuse |
| Nexus dockerconfigjson | admin:K@ng_... base64 | k8s/nexus/secret.yaml:8, k8s/charts/jeeb-infra/templates/nexus/secret.yaml:8 | Accepted (dev) |
| SonarQube token | squ_4cfe... | k8s/charts/jeeb-infra/values.yaml:39 | SHOULD ROTATE - tokens have no expiry by default |
| Grafana admin password | admin123 | k8s/charts/jeeb-obs/values.yaml:27 | Accepted (dev) |
| Keycloak RS256 public key | full PEM | k8s/charts/jeeb-infra/values.yaml:10-18 | Not a secret (public key) |
