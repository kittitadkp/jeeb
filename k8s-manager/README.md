# k8s-manager

`k8s-manager` is the operational CLI for bootstrapping and maintaining the local Jeeb Kubernetes stack on Docker Desktop.

## Prerequisites

- Docker Desktop with Kubernetes enabled
- `kubectl` and `helm` on `PATH`
- Go 1.23+

## Bootstrap inputs

```powershell
Copy-Item env/secrets.yaml.example env/secrets.yaml
Copy-Item env/config.yaml.example env/config.yaml   # optional
go run ./cmd/k8s-manager validate
```

`env/secrets.yaml` is required. `env/config.yaml` only overrides namespaces, NodePorts, hostnames, or Vault paths.

## Main commands

```text
setup
deploy [infra] [data] [app] [learning] [obs]
seed
check
maintain
trust-cert
validate
kong-key
patch-jenkins-creds
redeploy-jenkins
namespace
status
restart <deployment>
logs <deployment>
rancher
```

Example:

```powershell
go run ./cmd/k8s-manager trust-cert
```

## `setup` flow

`setup` currently runs these steps in order:

1. Pre-flight checks
2. Remove stale files
3. Generate `values-secrets.yaml`
4. Install ingress-nginx
5. Install cert-manager and Rancher
6. Deploy `jeeb-data`
7. Wait for Keycloak
8. Fetch Kong RS256 key from Keycloak
9. Deploy `jeeb-infra`
10. Wait for Kong
11. Wait for Vault
12. Initialize Vault
13. Store unseal keys in a Kubernetes secret
14. Unseal Vault
15. Configure Vault
16. Initialize the Nexus Docker registry
17. Patch CoreDNS for `.local` hosts
18. Wait for CoreDNS rollout
19. Verify in-cluster DNS resolution
20. Seed Jenkins jobs

After `setup`, the CLI instructs you to run Jenkins pipelines, publish images to Nexus, and then deploy `app` and `learning`.

## Service endpoints

| Service | NodePort |
|---|---:|
| Frontend | 30000 |
| Backend | 30080 |
| Keycloak | 30081 |
| Jenkins | 30082 |
| Nexus UI | 30083 |
| Learning backend | 30086 |
| Learning frontend | 30087 |
| Kong | 30088 |
| SonarQube | 30090 |
| Vault | 30091 |
| Grafana | 30092 |
| Prometheus | 30093 |
| MongoDB | 30017 |
| Nexus registry | 30050 |

## Notes

- `deploy` re-runs Helm upgrades without reinitializing Vault.
- `seed` creates the Jenkins seed job from `jenkins/jobs/seed.groovy`.
- `check` is the fast health gate; `maintain` prints diagnosis and fix commands.
- `trust-cert` imports `jeeb-dev-tls` into the Windows root trust store. The default `current-user` scope does not require admin rights.
- `coredns-patch.yaml` contains a concrete ingress ClusterIP and may need regeneration after a cluster reset.
