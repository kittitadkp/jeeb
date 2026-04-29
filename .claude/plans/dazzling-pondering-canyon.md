# Plan: Observability Stack — Logs, Metrics, Traces

## Context

The jeeb project has solid CI/CD and multi-env infrastructure but **zero observability**. Logs
are ephemeral (pod stdout only), there are no metrics, and there is no distributed tracing.
When a backend request fails, there is no way to correlate what happened across services, view
trends over time, or identify slow queries. This plan adds the full LGTM stack
(Loki + Grafana + Tempo + Prometheus) as a new Helm chart, instruments the Go backend with
OpenTelemetry and Prometheus, and wires everything together so logs, metrics, and traces are
linked by `trace_id` in a single Grafana UI.

---

## Target Architecture

```
Browser / Grafana UI  →  grafana.jeeb.local
                         │
              ┌──────────┼────────────┐
         Loki (logs) Prometheus   Tempo (traces)
              │       (metrics)       │
         Promtail    ← scrape     OTLP push
         DaemonSet    backend       from
         (all pods)   /metrics    backend
```

All observability services live in a new `jeeb-obs` namespace managed by a new
`jeeb-obs` Helm chart in `k8s/charts/jeeb-obs/`.

---

## New Files

```
k8s/charts/jeeb-obs/
  Chart.yaml
  values.yaml
  templates/
    namespace.yaml
    prometheus/
      deployment.yaml      # Prometheus server
      service.yaml
      configmap.yaml       # scrape_configs targeting jeeb-dev + jeeb-uat backend pods
      pvc.yaml
    loki/
      deployment.yaml      # Loki single-binary
      service.yaml
      configmap.yaml       # local filesystem storage
      pvc.yaml
    promtail/
      daemonset.yaml       # scrapes /var/log/pods on every node
      serviceaccount.yaml
      clusterrole.yaml
      clusterrolebinding.yaml
      configmap.yaml       # pipeline_stages: add namespace/pod/container labels
    tempo/
      deployment.yaml      # Tempo single-binary
      service.yaml         # ports: 4317 (OTLP gRPC), 4318 (OTLP HTTP), 3200 (query)
      configmap.yaml
      pvc.yaml
    grafana/
      deployment.yaml
      service.yaml         # NodePort 30092
      pvc.yaml
      configmap-datasources.yaml   # auto-provision Prometheus + Loki + Tempo datasources
      configmap-dashboards.yaml    # provisioning dir config
      configmap-dashboard-backend.yaml  # pre-built backend HTTP dashboard JSON
    ingress.yaml           # grafana.jeeb.local

k8s/apply.sh              # add: helm upgrade --install jeeb-obs ...
```

### Modified Files

```
backend/go.mod / go.sum                       # add OTel + Prometheus deps
backend/cmd/api/main.go                       # init tracer + metrics server
backend/internal/adapter/in/http/router.go   # add otelhttp + promhttp middleware
backend/internal/adapter/in/http/middleware/
  tracing.go   (new)   # init OTLP tracer provider
  metrics.go   (new)   # Prometheus registry + HTTP metrics middleware

k8s/charts/jeeb-app/templates/backend/
  configmap.yaml        # add OTEL_EXPORTER_OTLP_ENDPOINT
  deployment.yaml       # add prometheus.io/scrape annotations

k8s/charts/jeeb-infra/templates/ingress.yaml   # add grafana.jeeb.local rule
C:\Windows\System32\drivers\etc\hosts          # add grafana.jeeb.local (manual)
CoreDNS configmap                              # add grafana.jeeb.local → nginx ingress IP
```

---

## Steps

### Step 1 — jeeb-obs Helm chart scaffold

Create `k8s/charts/jeeb-obs/Chart.yaml` (name: jeeb-obs, version: 0.1.0).

Create `k8s/charts/jeeb-obs/values.yaml` with:
```yaml
namespace: jeeb-obs

prometheus:
  image: prom/prometheus:v2.52.0
  retention: 15d
  storageSize: 5Gi
  nodePort: 30093

loki:
  image: grafana/loki:3.0.0
  storageSize: 5Gi

tempo:
  image: grafana/tempo:2.4.2
  storageSize: 5Gi

promtail:
  image: grafana/promtail:3.0.0

grafana:
  image: grafana/grafana:11.0.0
  storageSize: 1Gi
  nodePort: 30092
  adminPassword: admin123
```

- [x] **Step 1a** — `templates/namespace.yaml`
- [x] **Step 1b** — Prometheus: deployment, service, pvc, configmap (scrape_configs)

  Scrape targets in prometheus configmap:
  ```yaml
  scrape_configs:
    - job_name: jeeb-backend-dev
      static_configs:
        - targets: ['backend.jeeb-dev.svc.cluster.local:8080']
    - job_name: jeeb-backend-uat
      static_configs:
        - targets: ['backend.jeeb-uat.svc.cluster.local:8080']
    - job_name: kubernetes-pods          # annotation-based scraping
      kubernetes_sd_configs:
        - role: pod
      relabel_configs: ...               # filter by prometheus.io/scrape=true
  ```

- [x] **Step 1c** — Loki: deployment, service, pvc, configmap (local filesystem, single binary)
- [x] **Step 1d** — Promtail: daemonset, serviceaccount, clusterrole + binding, configmap

  Promtail config pipes `/var/log/pods/**/*.log` to Loki, labels each stream with
  `namespace`, `pod`, `container`, `app`.

- [x] **Step 1e** — Tempo: deployment, service, pvc, configmap

  Tempo config enables OTLP receivers (gRPC 4317, HTTP 4318) and the query frontend (3200).
  Storage: local filesystem in single-binary mode.

- [x] **Step 1f** — Grafana: deployment, service (NodePort 30092), pvc, three configmaps

  **configmap-datasources.yaml** (mounted at `/etc/grafana/provisioning/datasources/`):
  ```yaml
  datasources:
    - name: Prometheus
      type: prometheus
      url: http://prometheus.jeeb-obs.svc.cluster.local:9090
      isDefault: true
    - name: Loki
      type: loki
      url: http://loki.jeeb-obs.svc.cluster.local:3100
    - name: Tempo
      type: tempo
      url: http://tempo.jeeb-obs.svc.cluster.local:3200
      jsonData:
        tracesToLogsV2:
          datasourceUid: loki       # click trace → jump to correlated logs
          tags: [{key: "app"}]
        logsToTraces:
          datasourceUid: tempo
  ```

  **configmap-dashboard-backend.yaml**: Pre-built dashboard with:
  - Request rate (req/s) by endpoint
  - P50/P95/P99 latency histogram
  - Error rate (4xx / 5xx)
  - Active connections
  - Go runtime metrics (goroutines, GC pause)

- [x] **Step 1g** — `templates/ingress.yaml`: `grafana.jeeb.local` → Grafana:3000

---

### Step 2 — Backend: add OpenTelemetry tracing

Add to `backend/go.mod`:
```
go.opentelemetry.io/otel                                          v1.27.0
go.opentelemetry.io/otel/sdk                                      v1.27.0
go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp  v1.27.0
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp     v0.52.0
go.opentelemetry.io/otel/trace                                    v1.27.0
```

- [x] **Step 2a** — Create `backend/internal/adapter/in/http/middleware/tracing.go`

  Exports `InitTracer(ctx, serviceName, otlpEndpoint) (func(), error)`:
  - Creates OTLP HTTP exporter pointing to `OTEL_EXPORTER_OTLP_ENDPOINT`
  - Registers a `TracerProvider` with resource attributes: `service.name=jeeb-backend`,
    `deployment.env=dev|uat` (from `GO_ENV`)
  - Sets global tracer provider + text map propagator (W3C TraceContext + Baggage)
  - Returns shutdown function (called on graceful shutdown)

- [x] **Step 2b** — Update `backend/cmd/api/main.go`:
  - Call `InitTracer` after config load, before server start
  - Defer the returned shutdown function

- [x] **Step 2c** — Update `backend/internal/adapter/in/http/router.go`:
  - Wrap the entire Chi router with `otelhttp.NewHandler(r, "jeeb-backend")`
  - This automatically creates a span per request with method, route, status code

- [x] **Step 2d** — Update `backend/internal/adapter/in/http/middleware/logging.go`:
  - Extract `trace.SpanFromContext(r.Context()).SpanContext().TraceID()` 
  - Add `trace_id` field to the slog record so every log line is correlated to its trace

---

### Step 3 — Backend: add Prometheus metrics

Add to `backend/go.mod`:
```
github.com/prometheus/client_golang  v1.19.1
```

- [x] **Step 3a** — Create `backend/internal/adapter/in/http/middleware/metrics.go`

  Defines:
  - `httpRequestsTotal` — Counter vec (method, route, status_code)
  - `httpRequestDuration` — Histogram vec (method, route, status_code), buckets: 5ms–10s
  - `httpRequestsInFlight` — Gauge

  Exports `PrometheusMiddleware(next http.Handler) http.Handler` that records all three.

- [x] **Step 3b** — Update `backend/internal/adapter/in/http/router.go`:
  - Add `PrometheusMiddleware` to the middleware chain (after RequestID, before Logging)
  - Register `GET /metrics` route (public, no auth) using `promhttp.Handler()`

- [x] **Step 3c** — Update `backend/cmd/api/main.go`:
  - Register custom collectors if any (Go runtime metrics are auto-registered)

---

### Step 4 — Wire Tempo endpoint into backend config

- [x] **Step 4a** — Update `k8s/charts/jeeb-app/templates/backend/configmap.yaml`:
  ```yaml
  OTEL_EXPORTER_OTLP_ENDPOINT: "http://tempo.jeeb-obs.svc.cluster.local:4318"
  OTEL_SERVICE_NAME: "jeeb-backend"
  ```

- [x] **Step 4b** — Update `backend/internal/config/config.go`:
  - Add `OtelEndpoint string \`envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT"\``
  - Add `OtelServiceName string \`envconfig:"OTEL_SERVICE_NAME" default:"jeeb-backend"\``
  - Pass to `InitTracer()`

- [x] **Step 4c** — Update `k8s/charts/jeeb-app/templates/backend/deployment.yaml`:
  - Add pod annotation: `prometheus.io/scrape: "true"`, `prometheus.io/port: "8080"`,
    `prometheus.io/path: "/metrics"`

---

### Step 5 — Hosts file + CoreDNS + apply.sh

- [x] **Step 5a** — Add to Windows hosts file (manual, admin shell):
  ```
  127.0.0.1 grafana.jeeb.local
  ```

- [x] **Step 5b** — Patch CoreDNS configmap to include `grafana.jeeb.local` in the hosts block

- [x] **Step 5c** — Update `k8s/apply.sh`:
  ```bash
  helm upgrade --install jeeb-obs k8s/charts/jeeb-obs --namespace jeeb-obs --create-namespace
  ```

- [x] **Step 5d** — Create `k8s/apply-obs.sh` as a convenience shortcut

---

### Step 6 — Rebuild and redeploy backend

- [x] **Step 6a** — Build and push updated backend image via Jenkins (or manually with `docker build/push`)
- [x] **Step 6b** — `helm upgrade jeeb-dev k8s/charts/jeeb-app -f k8s/charts/jeeb-app/values-dev.yaml`
- [x] **Step 6c** — Verify `/metrics` returns Prometheus text format
- [x] **Step 6d** — Verify traces appear in Tempo within 30s of making a request

---

## Verification

```bash
# 1. All obs pods running
kubectl get pods -n jeeb-obs

# 2. Prometheus scraping backend
curl http://localhost:30093/targets   # should show jeeb-backend-dev UP

# 3. Backend exposes metrics
curl http://localhost:30080/metrics   # Prometheus text format

# 4. Loki receiving logs
# In Grafana → Explore → Loki → {namespace="jeeb-dev"} — should see backend logs

# 5. Tempo receiving traces
# Make a request: curl http://api.jeeb-dev.local/health
# In Grafana → Explore → Tempo → search by service "jeeb-backend" — trace should appear

# 6. Log ↔ Trace correlation
# Click a trace span in Tempo → "Logs for this span" → should open Loki filtered by trace_id

# 7. Grafana dashboard
# http://grafana.jeeb.local → "Jeeb Backend" dashboard → request rate, latency, errors
```

---

## Resource Budget (Docker Desktop)

| Service | Memory Request | Memory Limit |
|---------|---------------|--------------|
| Prometheus | 256Mi | 512Mi |
| Loki | 128Mi | 256Mi |
| Tempo | 128Mi | 256Mi |
| Promtail | 64Mi | 128Mi |
| Grafana | 128Mi | 256Mi |
| **Total** | **704Mi** | **1.4Gi** |

All run in single-binary / single-replica mode to minimise local resource use.

---

## NodePort Additions

| Service | NodePort |
|---------|---------|
| Grafana | 30092 |
| Prometheus | 30093 |

---

## What This Unlocks

| Problem | Solution |
|---------|----------|
| "Which request caused the error?" | Trace ID in every log line → jump from log to Tempo trace |
| "Is the API getting slower over time?" | Prometheus latency histogram + Grafana dashboard |
| "What did pod X log before it crashed?" | Promtail ships logs to Loki before pod dies |
| "Where is the bottleneck in a request?" | Tempo waterfall view shows each span (MongoDB, Keycloak call, etc.) |
| "Is UAT behaving differently from dev?" | Both environments scraped; Grafana allows env filter |
