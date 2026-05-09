package redeploy

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"k8s-manager/internal/config"
	"k8s-manager/internal/logger"
	"k8s-manager/internal/util"
)

type JenkinsRedeployer struct {
	cfg       *config.ClusterConfig
	namespace string
	timeout   time.Duration
	dryRun    bool
}

func NewJenkinsRedeployer(cfg *config.ClusterConfig, namespace string, timeout time.Duration, dryRun bool) *JenkinsRedeployer {
	return &JenkinsRedeployer{cfg: cfg, namespace: namespace, timeout: timeout, dryRun: dryRun}
}

func (r *JenkinsRedeployer) Run(ctx context.Context) error {
	logger.Step("restarting jenkins deployment ...")
	if err := r.restart(ctx); err != nil {
		return fmt.Errorf("restart jenkins: %w", err)
	}

	logger.Step("waiting for rollout to complete ...")
	if err := r.waitRollout(ctx); err != nil {
		return fmt.Errorf("rollout status: %w", err)
	}

	logger.Step("waiting for jenkins to be ready ...")
	if err := r.healthCheck(ctx); err != nil {
		return fmt.Errorf("health check: %w", err)
	}

	logger.Step("done — jenkins is up")
	return nil
}

func (r *JenkinsRedeployer) restart(ctx context.Context) error {
	if r.dryRun {
		logger.Step("[dry-run] kubectl rollout restart deployment/jenkins -n %s", r.namespace)
		return nil
	}
	return util.RunCmd(ctx, "kubectl", "rollout", "restart", "deployment/jenkins", "-n", r.namespace)
}

func (r *JenkinsRedeployer) waitRollout(ctx context.Context) error {
	if r.dryRun {
		logger.Step("[dry-run] kubectl rollout status deployment/jenkins -n %s --timeout=%s", r.namespace, r.timeout)
		return nil
	}
	return util.RunCmd(ctx, "kubectl", "rollout", "status", "deployment/jenkins",
		"-n", r.namespace, "--timeout="+r.timeout.String())
}

func (r *JenkinsRedeployer) healthCheck(ctx context.Context) error {
	url := fmt.Sprintf("http://localhost:%d/login", r.cfg.JenkinsNodePort)
	if r.dryRun {
		logger.Step("[dry-run] poll %s", url)
		return nil
	}
	return util.PollHTTP(ctx,
		util.PollConfig{Timeout: r.timeout, Interval: 5 * time.Second},
		&http.Client{Timeout: 10 * time.Second},
		http.MethodGet, url, http.StatusOK, nil,
	)
}
