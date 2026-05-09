# Kubernetes Layout

The repository deploys Jeeb through Helm charts rather than raw `kubectl apply` workflows.

## Charts

- `charts/jeeb-data`: MongoDB and Keycloak
- `charts/jeeb-infra`: Kong, Jenkins, Nexus, SonarQube, Vault
- `charts/jeeb-app`: main backend and main frontend
- `charts/jeeb-learning`: learning backend and learning frontend
- `charts/jeeb-obs`: Prometheus, Loki, Tempo, Grafana

## Supporting manifests

- `tls-issuer.yaml`: local wildcard certificate issuer
- `coredns-patch.yaml`: static host mappings for `.local` ingress names

## Current deployment model

- Namespaces default to `jeeb-dev`, `jeeb-infra`, and `jeeb-obs`.
- App images are pulled from the local Nexus registry at `localhost:30050`.
- Backends load Vault-rendered env files before starting.
- Frontend deployments also mount Vault-rendered env files, but the current Nginx-based images do not consume them.

## Public hosts

- `jeeb-dev.local`
- `api.jeeb-dev.local`
- `auth.jeeb-dev.local`
- `learning.jeeb-dev.local`
- `learning-api.jeeb-dev.local`
- `jenkins.jeeb.local`
- `nexus.jeeb.local`
- `sonarqube.jeeb.local`
- `vault.jeeb.local`
