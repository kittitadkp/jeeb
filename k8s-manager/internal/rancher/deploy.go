package rancher

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"k8s-manager/internal/config"
)

const (
	certManagerVersion = "v1.15.3"
	rancherVersion     = "2.9.3"
)

type Deployer struct {
	cfg    *config.ClusterConfig
	dryRun bool
}

func NewDeployer(cfg *config.ClusterConfig, dryRun bool) *Deployer {
	return &Deployer{cfg: cfg, dryRun: dryRun}
}

func (d *Deployer) Run(ctx context.Context) error {
	steps := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"Add / update Helm repos", d.addRepos},
		{"Install cert-manager " + certManagerVersion, d.installCertManager},
		{"Wait for cert-manager pods", d.waitCertManager},
		{"Install Rancher " + rancherVersion, d.installRancher},
		{"Wait for Rancher rollout", d.waitRancher},
		{"Patch Rancher service to NodePort " + fmt.Sprint(d.cfg.RancherNodePort), d.patchNodePort},
	}

	fmt.Println("=== Rancher Setup ===")
	if d.dryRun {
		fmt.Println("DRY RUN — commands will be printed, not executed")
	}
	fmt.Println()

	for i, s := range steps {
		fmt.Printf("[%d/%d] %s\n", i+1, len(steps), s.name)
		if err := s.fn(ctx); err != nil {
			return fmt.Errorf("step %d (%s): %w", i+1, s.name, err)
		}
		fmt.Printf("      done\n\n")
	}

	fmt.Printf(`=== Rancher deployed ===

  Via NodePort  https://localhost:%d  (accept self-signed cert warning)
  Via ingress   https://%s

  Bootstrap password: admin  (you will be prompted to change it on first login)

  Add to C:\Windows\System32\drivers\etc\hosts (as Administrator):
    127.0.0.1  %s
`, d.cfg.RancherNodePort, d.cfg.RancherHostname, d.cfg.RancherHostname)
	return nil
}

func (d *Deployer) addRepos(ctx context.Context) error {
	cmds := [][]string{
		{"repo", "add", "jetstack", "https://charts.jetstack.io", "--force-update"},
		{"repo", "add", "rancher-stable", "https://releases.rancher.com/server-charts/stable", "--force-update"},
		{"repo", "update"},
	}
	for _, args := range cmds {
		if err := d.helm(ctx, args...); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deployer) installCertManager(ctx context.Context) error {
	return d.helm(ctx,
		"upgrade", "--install", "cert-manager", "jetstack/cert-manager",
		"--namespace", d.cfg.CertManagerNamespace,
		"--create-namespace",
		"--version", certManagerVersion,
		"--set", "crds.enabled=true",
	)
}

func (d *Deployer) waitCertManager(ctx context.Context) error {
	deployments := []string{"cert-manager", "cert-manager-webhook", "cert-manager-cainjector"}
	for _, dep := range deployments {
		fmt.Printf("      waiting for %s...\n", dep)
		if err := d.kubectl(ctx,
			"rollout", "status", "deployment/"+dep,
			"-n", d.cfg.CertManagerNamespace,
			"--timeout=120s",
		); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deployer) installRancher(ctx context.Context) error {
	return d.helm(ctx,
		"upgrade", "--install", "rancher", "rancher-stable/rancher",
		"--namespace", d.cfg.RancherNamespace,
		"--create-namespace",
		"--version", rancherVersion,
		"--set", "hostname="+d.cfg.RancherHostname,
		"--set", "ingress.tls.source=rancher",
		"--set", "ingress.ingressClassName=nginx",
		"--set", "replicas=1",
		"--set", "bootstrapPassword=admin",
	)
}

func (d *Deployer) waitRancher(ctx context.Context) error {
	return d.kubectl(ctx,
		"rollout", "status", "deployment/rancher",
		"-n", d.cfg.RancherNamespace,
		"--timeout=300s",
	)
}

func (d *Deployer) patchNodePort(ctx context.Context) error {
	patch := fmt.Sprintf(`[
    {"op":"replace","path":"/spec/type","value":"NodePort"},
    {"op":"add","path":"/spec/ports/0/nodePort","value":%d}
  ]`, d.cfg.RancherNodePort)
	return d.kubectl(ctx,
		"patch", "svc", "rancher",
		"-n", d.cfg.RancherNamespace,
		"--type=json",
		"-p="+patch,
	)
}

func (d *Deployer) helm(ctx context.Context, args ...string) error {
	if d.dryRun {
		fmt.Printf("      helm %s\n", strings.Join(args, " "))
		return nil
	}
	cmd := exec.CommandContext(ctx, "helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d *Deployer) kubectl(ctx context.Context, args ...string) error {
	if d.dryRun {
		fmt.Printf("      kubectl %s\n", strings.Join(args, " "))
		return nil
	}
	cmd := exec.CommandContext(ctx, "kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
