package jenkins

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s-manager/internal/credentials"
	"k8s-manager/internal/logger"
	"k8s-manager/internal/util"
)

// CredentialsPatcher updates the jenkins-secret Kubernetes secret from
// secrets.yaml and rolls out a Jenkins restart so the new values take effect.
type CredentialsPatcher struct {
	namespace string
	creds     *credentials.Credentials
	dryRun    bool
}

func NewCredentialsPatcher(creds *credentials.Credentials, namespace string, dryRun bool) *CredentialsPatcher {
	ns := namespace
	if ns == "" {
		ns = "jeeb-infra"
	}
	return &CredentialsPatcher{namespace: ns, creds: creds, dryRun: dryRun}
}

func (p *CredentialsPatcher) Run(ctx context.Context) error {
	logger.Step("patching jenkins-secret in namespace %s ...", p.namespace)
	if err := p.patchSecret(ctx); err != nil {
		return fmt.Errorf("patch jenkins-secret: %w", err)
	}

	logger.Step("restarting jenkins deployment ...")
	if err := p.restartDeployment(ctx); err != nil {
		return fmt.Errorf("restart jenkins: %w", err)
	}

	logger.Step("waiting for rollout to complete ...")
	if err := p.waitRollout(ctx); err != nil {
		return fmt.Errorf("rollout status: %w", err)
	}

	logger.Step("done — jenkins credentials updated")
	return nil
}

func (p *CredentialsPatcher) patchSecret(ctx context.Context) error {
	patch := map[string]interface{}{
		"stringData": map[string]string{
			"admin-password": p.creds.JenkinsAdminPassword,
			"github-user":    p.creds.JenkinsGithubUser,
			"github-pat":     p.creds.JenkinsGithubPAT,
			"nexus-user":     p.creds.JenkinsNexusUser,
			"nexus-password": p.creds.JenkinsNexusPassword,
			"sonar-token":    p.creds.JenkinsSonarToken,
		},
	}

	patchJSON, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	if p.dryRun {
		logger.Step("[dry-run] kubectl patch secret jenkins-secret -n %s --type=merge -p '%s'",
			p.namespace, string(patchJSON))
		return nil
	}

	return util.RunCmd(ctx, "kubectl", "patch", "secret", "jenkins-secret",
		"-n", p.namespace,
		"--type=merge",
		"-p", string(patchJSON),
	)
}

func (p *CredentialsPatcher) restartDeployment(ctx context.Context) error {
	if p.dryRun {
		logger.Step("[dry-run] kubectl rollout restart deployment/jenkins -n %s", p.namespace)
		return nil
	}
	return util.RunCmd(ctx, "kubectl", "rollout", "restart", "deployment/jenkins", "-n", p.namespace)
}

func (p *CredentialsPatcher) waitRollout(ctx context.Context) error {
	if p.dryRun {
		return nil
	}
	return util.RunCmd(ctx, "kubectl", "rollout", "status", "deployment/jenkins",
		"-n", p.namespace, "--timeout=3m")
}
