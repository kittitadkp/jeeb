# Plan: Add `redeploy-jenkins` Command to k8s-manager CLI

## Context

The user needs a single command that completely cycles Jenkins: **Helm uninstall → Helm install → seed** (create seed job + generate all pipeline jobs). This is for when Jenkins is in a bad state and a plain rollout restart isn't enough — it wipes and recreates the release from scratch, then re-seeds.

---

## What the command does (3 steps)

| Step | Action | Reuses |
|------|--------|--------|
| 1. Helm uninstall | `helm uninstall jeeb-infra -n jeeb-infra` — ignores "not found" | `helm.Run()` + `util.DryRunOrExec()` |
| 2. Helm install | `setup.NewRunner(...).Deploy(ctx, []string{"infra"})` — generates values file + runs helm upgrade --install | existing `setup` package |
| 3. Seed | `jenkins.NewSeeder(...).Run(ctx)` — waits for Jenkins, configures env vars, creates/runs seed job | existing `jenkins.Seeder` |

---

## Files to create / modify

| File | Action |
|------|--------|
| `k8s-manager/internal/redeploy/jenkins.go` | **Created** — `JenkinsRedeployer` struct + `Run` method (own package to avoid import cycle with `setup`) |
| `k8s-manager/cmd/k8s-manager/main.go` | **Modified** — added `newRedeployJenkinsCmd()` and registered it in `root.AddCommand(...)` |

---

## `internal/jenkins/redeploy.go`

```go
package jenkins

import (
    "context"
    "fmt"
    "strings"

    "k8s-manager/internal/config"
    "k8s-manager/internal/credentials"
    "k8s-manager/internal/helm"
    "k8s-manager/internal/logger"
    "k8s-manager/internal/setup"
)

type Redeployer struct {
    cfg        *config.ClusterConfig
    creds      *credentials.Credentials
    chartsDir  string
    outputDir  string
    secretsFile string
    groovyPath string
    namespace  string
    dryRun     bool
}

func NewRedeployer(cfg *config.ClusterConfig, creds *credentials.Credentials,
    chartsDir, outputDir, secretsFile, groovyPath, namespace string, dryRun bool) *Redeployer {
    return &Redeployer{
        cfg: cfg, creds: creds,
        chartsDir: chartsDir, outputDir: outputDir,
        secretsFile: secretsFile, groovyPath: groovyPath,
        namespace: namespace, dryRun: dryRun,
    }
}

func (r *Redeployer) Run(ctx context.Context) error {
    steps := []struct {
        name string
        fn   func(context.Context) error
    }{
        {"helm uninstall jeeb-infra", r.uninstall},
        {"helm install jeeb-infra",   r.install},
        {"seed jenkins",              r.seed},
    }
    for _, s := range steps {
        logger.Step("==> %s", s.name)
        if err := s.fn(ctx); err != nil {
            return fmt.Errorf("%s: %w", s.name, err)
        }
    }
    return nil
}

func (r *Redeployer) uninstall(ctx context.Context) error {
    err := helm.Run(ctx, r.dryRun, "uninstall", cfg.ReleaseInfra, "-n", r.namespace)
    if err != nil && !strings.Contains(err.Error(), "not found") {
        return err
    }
    return nil
}

func (r *Redeployer) install(ctx context.Context) error {
    runner := setup.NewRunner(r.cfg, r.creds, r.chartsDir, r.outputDir, r.secretsFile, r.dryRun)
    return runner.Deploy(ctx, []string{"infra"})
}

func (r *Redeployer) seed(ctx context.Context) error {
    return NewSeeder(r.cfg, r.creds, r.groovyPath, "", r.dryRun).Run(ctx)
}
```

> **Key reuse:** `helm.Run()` already handles dry-run; `setup.Runner.Deploy(["infra"])` handles values-file generation + helm upgrade --install; `jenkins.Seeder` already handles waiting + seeding.

---

## `cmd/k8s-manager/main.go` changes

Add constructor function:

```go
func newRedeployJenkinsCmd(secretsFile, configFile, chartsDir, outputDir *string, dryRun *bool) *cobra.Command {
    var namespace string
    var groovyPath string

    cmd := &cobra.Command{
        Use:   "redeploy-jenkins",
        Short: "Helm uninstall + install jeeb-infra, then seed Jenkins",
        Long: `Completely reinstalls the jeeb-infra Helm release and re-seeds Jenkins.
Steps:
  1. helm uninstall jeeb-infra  (ignores "not found")
  2. helm upgrade --install jeeb-infra  (same as: k8s-manager deploy infra)
  3. seed Jenkins  (same as: k8s-manager seed)

Use this when Jenkins is in a broken state that a rollout restart cannot fix.`,
        RunE: func(cmd *cobra.Command, args []string) error {
            secretsPath := resolvedSecretsFile(*secretsFile)
            creds, err := loadCreds(secretsPath)
            if err != nil {
                return err
            }
            cfg := config.LoadFromFile(resolvedConfigFile(*configFile))
            r := jenkins.NewRedeployer(cfg, creds, *chartsDir, *outputDir, secretsPath, groovyPath, namespace, *dryRun)
            return r.Run(cmd.Context())
        },
    }

    cmd.Flags().StringVarP(&namespace, "namespace", "n", "jeeb-infra", "namespace where jeeb-infra is deployed")
    cmd.Flags().StringVar(&groovyPath, "groovy-path", defaultGroovyPath(), "path to seed.groovy")
    return cmd
}
```

Register it in `root.AddCommand(...)` alongside the others.

---

## Verification

```bash
cd k8s-manager

# 1. Compile
go build ./cmd/k8s-manager/...

# 2. Dry-run — no cluster needed, prints what would happen
./k8s-manager redeploy-jenkins --dry-run

# 3. Live run
./k8s-manager redeploy-jenkins

# 4. Verify Jenkins is healthy afterward
./k8s-manager check
```

Expected dry-run output:
```
==> helm uninstall jeeb-infra
[dry-run] helm uninstall jeeb-infra -n jeeb-infra
==> helm install jeeb-infra
[dry-run] helm upgrade --install jeeb-infra <charts-dir>/jeeb-infra --namespace jeeb-infra --create-namespace -f <values-file>
==> seed jenkins
[dry-run] would read seed.groovy, create seed job, trigger build
```
