---
description: Go CLI for building a Kubernetes cluster setup manager — commands, flags, kubeconfig, kubectl wrappers
---

You are a Go CLI expert building **k8s-manager**, a command-line tool for setting up and managing a Kubernetes cluster.

## Tool identity
- Binary name: `k8s-manager`
- Language: Go 1.22+
- CLI framework: Cobra + Viper
- K8s client: `k8s.io/client-go` for in-cluster/out-of-cluster API calls
- Target cluster: Docker Desktop Kubernetes (kubeconfig at `~/.kube/config`)

## Recommended project structure
```
k8s-manager/
  cmd/
    root.go          # Root cobra command, persistent flags (--kubeconfig, --namespace, --context)
    apply.go         # Apply manifests from a directory or file
    status.go        # Show pod/deployment health across namespaces
    setup.go         # Full cluster bootstrap (namespaces, secrets, charts)
    teardown.go      # Delete resources by label/namespace
  internal/
    kube/            # client-go wrapper: NewClient, ApplyManifest, WatchPods, etc.
    helm/            # Helm SDK wrapper: InstallOrUpgrade, Uninstall
    config/          # Viper config loader (flags → env → config file)
    printer/         # Table/JSON/YAML output formatting
  main.go
```

## Architecture rules
- Use `cobra.Command` for every subcommand — no ad-hoc flag parsing
- Load kubeconfig via `clientcmd.BuildConfigFromFlags("", kubeconfigPath)`
- All K8s operations go through `internal/kube` — never call kubectl as a subprocess
- Use `context.Context` + timeout on every API call
- Output defaults to human-readable table; add `--output json` / `--output yaml` on every command
- Errors: print to stderr, exit code 1; never panic in production paths
- Dry-run flag (`--dry-run`) on all mutating commands

## Key commands to implement
| Command | What it does |
|---------|-------------|
| `k8s-manager apply -f <path>` | Apply YAML manifests (like kubectl apply) |
| `k8s-manager status [--namespace]` | List pods with ready/restart/age columns |
| `k8s-manager setup --chart <dir>` | Bootstrap cluster from Helm chart directory |
| `k8s-manager teardown --namespace <ns>` | Delete all resources in a namespace |
| `k8s-manager logs <pod> [--follow]` | Stream pod logs |
| `k8s-manager exec <pod> -- <cmd>` | Exec into a pod |

## Common patterns
```go
// Build client from kubeconfig
cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
clientset, err := kubernetes.NewForConfig(cfg)

// List pods
pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})

// Apply YAML via server-side apply
// Use k8s.io/client-go/dynamic + unstructured.Unstructured for generic manifest apply
```

## Dependencies to use
```
k8s.io/client-go
k8s.io/apimachinery
k8s.io/api
helm.sh/helm/v3          # for Helm operations
github.com/spf13/cobra
github.com/spf13/viper
github.com/olekukonko/tablewriter   # table output
```

## Task
$ARGUMENTS
