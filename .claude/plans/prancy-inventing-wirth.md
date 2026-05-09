# Plan: k8s-manager Rewrite

## Context

Three improvements to the k8s-manager CLI:
1. All `fmt.Printf`/`fmt.Println` logging replaced with `github.com/rs/zerolog` (ConsoleWriter for CLI output), with `--log-level` flag controlling verbosity
2. Repeated patterns (HTTP polling, command exec, HTTP client setup) extracted into a shared `internal/util` package
3. `env/credentials.yaml` split into `env/secrets.yaml` (passwords/tokens only) + `env/config.yaml` (cluster topology overrides, all optional)

---

## New Files

| Path | Purpose |
|------|---------|
| `internal/logger/logger.go` | slog init, `Step`/`Info`/`Debug`/`Warn`/`Error` helpers |
| `internal/util/poll.go` | `PollUntil`, `PollHTTP` — eliminate 4 duplicate deadline loops |
| `internal/util/exec.go` | `RunCmd`, `RunCmdOutput`, `RunCmdStdin`, `DryRunOrExec`, `DryRunOrExecStdin` |
| `internal/util/http.go` | `NewBasicAuthClient`, `DoJSON`, `FetchRS256PEM`, `JWKToPEM` |
| `env/secrets.yaml.example` | Operator template (passwords, PATs, tokens) |
| `env/config.yaml.example` | Cluster topology template (all optional, shows defaults) |

---

## Step 1 — util package

Create `internal/util/` with three files. No dependencies on other internal packages.

**`poll.go`**
```go
type PollConfig struct {
    Timeout  time.Duration
    Interval time.Duration // default 5s
}

func PollUntil(ctx context.Context, cfg PollConfig, check func(ctx context.Context) error) error
func PollHTTP(ctx context.Context, cfg PollConfig, client *http.Client, method, url string, wantStatus int, beforeReq func(*http.Request)) error
```

**`exec.go`**
```go
func RunCmd(ctx context.Context, name string, args ...string) error
func RunCmdOutput(ctx context.Context, name string, args ...string) ([]byte, error)
func RunCmdStdin(ctx context.Context, stdin io.Reader, name string, args ...string) error
func DryRunOrExec(ctx context.Context, dryRun bool, name string, args ...string) error
func DryRunOrExecStdin(ctx context.Context, dryRun bool, stdin io.Reader, name string, args ...string) error
```

**`http.go`**
```go
func NewBasicAuthClient() *http.Client            // cookie jar + 30s timeout
func DoJSON(ctx context.Context, client *http.Client, url string, out any) error
func FetchRS256PEM(ctx context.Context, jwksURL string) (string, error)
func JWKToPEM(nB64, eB64 string) (string, error) // moved from setup.go + kongkey/update.go
```

---

## Step 2 — logger package

Add dependency: `github.com/rs/zerolog` — run `go get github.com/rs/zerolog`.

Create `internal/logger/logger.go`:

```go
// Init configures the global zerolog logger. Call from PersistentPreRunE.
// level is one of: "debug"|"info"|"warn"|"error" (case-insensitive).
// Uses zerolog.ConsoleWriter for human-readable CLI output.
func Init(level string)

// Step writes user-facing progress. Always visible (not filtered by log level).
func Step(format string, args ...any)
func StepMsg(msg string)

// These forward to zerolog; controlled by --log-level.
func Info(msg string, args ...any)
func Debug(msg string, args ...any)
func Warn(msg string, args ...any)
func Error(msg string, args ...any)
```

**Implementation sketch:**
```go
var log zerolog.Logger

func Init(level string) {
    lvl, err := zerolog.ParseLevel(strings.ToLower(level))
    if err != nil {
        lvl = zerolog.InfoLevel
    }
    zerolog.SetGlobalLevel(lvl)
    w := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly}
    log = zerolog.New(w).With().Timestamp().Logger()
}

func Debug(msg string, args ...any) { log.Debug().Fields(args).Msg(msg) }
func Info(msg string, args ...any)  { log.Info().Fields(args).Msg(msg) }
func Warn(msg string, args ...any)  { log.Warn().Fields(args).Msg(msg) }
func Error(msg string, args ...any) { log.Error().Fields(args).Msg(msg) }

func Step(format string, args ...any) { fmt.Fprintf(os.Stdout, format+"\n", args...) }
func StepMsg(msg string)              { fmt.Fprintln(os.Stdout, msg) }
```

**Output distinction:**
- `logger.Step` / `logger.StepMsg` → `fmt.Fprintf(os.Stdout)` — always printed (step banners, "→ Deploying...")
- `logger.Info/Debug/Warn/Error` → zerolog ConsoleWriter to stderr — filtered by `--log-level`
- `validate.go`, `printer/pods.go`, `healthcheck/check.go` `Print()` → keep `fmt` (user-facing tables, not logging)

---

## Step 3 — env file split

### `env/secrets.yaml` (operator fills in)
```yaml
jenkins:
  adminPassword: ""
  githubUser: ""
  githubPat: ""
  nexusUser: "admin"
  nexusPat: ""
  sonarToken: ""
  githubCredsId: "github-creds"
  jenkinsRepo: ""      # optional; defaults to https://github.com/<githubUser>/jenkins.git
  k8sRepo: ""
  backendRepo: ""
  frontendRepo: ""
  learningBackendRepo: ""
  learningFrontendRepo: ""
keycloak:
  adminUser: "admin"
  adminPassword: ""
mongodb:
  username: "jeeb"
  password: ""
nexus:
  adminPassword: ""
sonarqube:
  adminPassword: ""
grafana:
  adminPassword: ""
kong:
  keycloakPublicKey: ""  # filled by 'k8s-manager kong-key'
```

### `env/config.yaml` (all optional — cluster topology overrides)
```yaml
cluster:
  namespaces:
    dev: jeeb-dev
    infra: jeeb-infra
    obs: jeeb-obs
  releases:
    data: jeeb-data
    dev: jeeb-dev
    infra: jeeb-infra
    learning: jeeb-learning
    obs: jeeb-obs
  vault:
    pod: vault-0
    addr: http://127.0.0.1:8200
  keycloak:
    realm: jeeb
    clientId: jeeb-app
    host: keycloak.jeeb-dev.svc.cluster.local:8080
    nodePort: 30081
  kong:
    issuer: http://auth.jeeb-dev.local/realms/jeeb
    nodePort: 30088
  mongo:
    host: mongodb.jeeb-dev.svc.cluster.local:27017
    nodePort: 30017
  nodePorts:
    frontend: 30000
    backend: 30080
    jenkins: 30082
    nexusUI: 30083
    learningBackend: 30086
    learningFront: 30087
    sonarQube: 30090
    vault: 30091
    grafana: 30092
    prometheus: 30093
  nexus:
    registry: localhost:30050
  ingress:
    label: app.kubernetes.io/name=ingress-nginx
  rancher:
    nodePort: 30443
    hostname: rancher.jeeb-infra.local
    namespace: cattle-system
  vaultPaths:
    backend: secret/data/jeeb/backend/develop
    frontend: secret/data/jeeb/frontend/develop
    learningBackend: secret/data/jeeb/learning/backend/develop
    learningFrontend: secret/data/jeeb/learning/frontend/develop
```

`env/.gitignore` — add `secrets.yaml` and `config.yaml`.

---

## Step 4 — main.go flags

| Old flag | New flag | Notes |
|----------|----------|-------|
| `--credentials` | `--secrets` | Path to `env/secrets.yaml` |
| _(none)_ | `--config` | Path to `env/config.yaml` (optional; omit to use all defaults) |
| _(none)_ | `--log-level` | `debug/info/warn/error`, default `info` |

Backward compat: if `--secrets` not set and `env/secrets.yaml` not found, fall back to `env/credentials.yaml` with `logger.Warn("credentials.yaml is deprecated; rename to secrets.yaml")`. Use cobra's `MarkDeprecated("credentials", "use --secrets instead")` if keeping old flag as alias.

```go
root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
    logger.Init(logLevel)
    return nil
}
```

Thread `configFile *string` into all command constructors that call `config.LoadFromFile`.

---

## Step 5 — per-file modifications

### `internal/credentials/loader.go`
- Signature unchanged: `Load(path string)`
- Remove any awareness of `cluster:` key (that section moved to config.yaml)
- Update comments to reference secrets.yaml

### `internal/config/config.go`
- `LoadFromFile(path string)` signature unchanged
- Now always receives the config file path from main.go (separate from secrets path)
- No structural change to YAML parsing — `cluster:` key is identical

### `internal/setup/setup.go` (largest change)
- Replace all progress `fmt.Printf` → `logger.Step` / `logger.StepMsg`
- Replace warning `fmt.Printf` → `logger.Warn`
- Replace `r.kubectl(ctx, args...)` body → `util.DryRunOrExec(ctx, r.dryRun, "kubectl", args...)`
- Replace `waitForKeycloak` polling loop → `util.PollHTTP`
- Replace Nexus wait loop → `util.PollHTTP`
- Remove `fetchKeycloakPublicKey` + private `jwkToPEM` → use `util.FetchRS256PEM`
- Rename `r.credsPath` → `r.secretsPath`

### `internal/jenkins/seed.go`
- `newHTTPClient()` → `util.NewBasicAuthClient()`
- `waitForJenkins` loop → `util.PollHTTP` with `beforeReq` setting basic auth
- `waitForBuild` loops → `util.PollUntil`
- Progress `fmt.Printf` → `logger.Step` / `logger.Debug`

### `internal/kongkey/update.go`
- `fetchFromKeycloak()` → `util.FetchRS256PEM(ctx, jwksURL)`
- Remove private `jwkToPEM` (canonical copy now in `util/http.go`)
- `fmt.Printf` → `logger.Step` / `logger.Debug`
- `u.credsPath` → `u.secretsPath`

### `internal/rancher/deploy.go`
- Private `d.helm()` / `d.kubectl()` → `util.DryRunOrExec`
- `fmt.Printf` → `logger.Step` / `logger.Debug`

### `internal/helm/runner.go`
- `exec.CommandContext` → `util.RunCmd`
- `fmt.Printf` (dry-run print) → `logger.Step`

### `internal/kube/client.go`
- Namespace created/exists messages → `logger.Info` / `logger.Debug`
- Restart/log progress → `logger.Info`
- Table output in `PrintStatus` stays as `fmt` (via `printer.PrintPods`)

### `internal/healthcheck/check.go`
- `Print()` keeps all `fmt` — it is a user-facing results table
- HTTP check helpers: no fmt calls currently; add `logger.Debug` for request URLs if useful

### `internal/maintain/report.go`
- Section headers (`fmt.Println("=== ...")`) → `logger.StepMsg`
- Fix command output (the `fmt.Printf` blocks showing kubectl/vault commands to run) → keep `fmt` — those are user-facing instructions

---

## Implementation order

- [x] 1. `internal/util/` — build compiles, no callers yet
- [x] 2. `internal/logger/` — build compiles, no callers yet
- [x] 3. `cmd/k8s-manager/main.go` — add flags + `PersistentPreRunE`, keep old `--credentials` as deprecated alias
- [x] 4. `env/` — create example files, update `.gitignore`
- [x] 5. `internal/setup/setup.go` — bulk of the work
- [x] 6. `internal/kongkey/update.go`
- [x] 7. `internal/jenkins/seed.go`
- [x] 8. `internal/rancher/deploy.go`
- [x] 9. `internal/helm/runner.go`, `kube/client.go`

---

## Verification

```powershell
# From D:\personal\jeeb\k8s-manager

# Build
go build ./...

# Tests
go test ./...

# Vet
go vet ./...

# New flags visible
go run ./cmd/k8s-manager --help

# --log-level works (debug output visible)
go run ./cmd/k8s-manager --log-level debug --secrets env/secrets.yaml.example validate

# --config accepted (no error even if cluster.yaml has no overrides)
go run ./cmd/k8s-manager --secrets env/secrets.yaml.example --config env/config.yaml.example validate

# Dry-run setup: confirm logger.Step output, no panics
go run ./cmd/k8s-manager --dry-run --secrets env/secrets.yaml.example setup

# Deprecated --credentials flag still works with warning
go run ./cmd/k8s-manager --credentials env/credentials.yaml validate
```

---

## Files NOT changed

- `internal/validate/validate.go` — user-facing table, keep fmt
- `internal/printer/pods.go` — user-facing table, keep fmt
- `internal/setup/secrets_values.go` — pure data transformation, no output
- `internal/setup/secrets_values_test.go` — test file
