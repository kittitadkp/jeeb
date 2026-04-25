# Architecture Overview

## Application — Clean Architecture

```
┌─────────────────────────────────────┐
│           HTTP Handlers             │  ← adapter/in
├─────────────────────────────────────┤
│            Use Cases                │  ← usecase
├─────────────────────────────────────┤
│             Domain                  │  ← domain
├─────────────────────────────────────┤
│    MongoDB │ Calendar │ Notify      │  ← adapter/out
└─────────────────────────────────────┘
```

Dependencies point inward. Domain has no external dependencies.

## Ports

| Port | Purpose |
|------|---------|
| `CalendarPort` | External calendar sync |
| `NotificationPort` | Push notifications |
| `*Repository` | Data persistence |

## Feature Modules

Each feature (workout, study, sleep, finance) is self-contained:
- Own domain models
- Own use cases
- Own repository interface

---

## Infrastructure — Kubernetes

All services run in the `jeeb` namespace on Kubernetes (Docker Desktop).

```
k8s/
├── 00-namespace.yaml
├── apply.sh
├── app/
│   ├── secrets.yaml
│   ├── backend/        configmap, deployment, service
│   ├── frontend/       deployment, service
│   ├── keycloak/       deployment, service
│   └── mongodb/        statefulset, service
├── jenkins/            deployment, pvc, rbac, service
├── nexus/              deployment, pvc, service
└── sonarqube/          deployment, pvc, service
```

## Service Ports (NodePort)

| Service | NodePort | In-cluster URL |
|---------|----------|----------------|
| Frontend | 30000 | frontend.jeeb.svc.cluster.local:80 |
| Backend | 30080 | backend.jeeb.svc.cluster.local:8080 |
| Keycloak | 30081 | keycloak.jeeb.svc.cluster.local:8080 |
| Jenkins | 30082 | jenkins.jeeb.svc.cluster.local:8080 |
| Nexus UI | 30083 | nexus.jeeb.svc.cluster.local:8081 |
| SonarQube | 30090 | sonarqube.jeeb.svc.cluster.local:9000 |
| MongoDB | 30017 | mongodb.jeeb.svc.cluster.local:27017 |
| Nexus Registry | 30050 | nexus.jeeb.svc.cluster.local:5000 |

## CI/CD Pipeline

```
git push
  → Jenkins polls SCM (every 5 min)
      → Kubernetes pod agent spins up
          → Test → SonarQube → Kaniko build → Push to Nexus
              → kubectl set image → rollout
```

- **Jenkins** — pipeline orchestration (Kubernetes plugin, pod agents)
- **Nexus** — Docker image registry (`localhost:30050/jeeb/<service>`)
- **SonarQube** — static analysis & code quality gate
- **Kaniko** — in-cluster Docker image builds (no Docker daemon)
