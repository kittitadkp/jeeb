package kongkey

import (
	"context"
	"fmt"
	"os"
	"strings"

	"k8s-manager/internal/config"
	"k8s-manager/internal/credentials"
	"k8s-manager/internal/logger"
	"k8s-manager/internal/setup"
	"k8s-manager/internal/util"

	"gopkg.in/yaml.v3"
)

// Updater handles fetching, storing, and deploying the Keycloak public key.
type Updater struct {
	cfg         *config.ClusterConfig
	secretsPath string
	chartsDir   string
	outputDir   string
	dryRun      bool
}

func NewUpdater(cfg *config.ClusterConfig, secretsPath, chartsDir, outputDir string, dryRun bool) *Updater {
	return &Updater{
		cfg:         cfg,
		secretsPath: secretsPath,
		chartsDir:   chartsDir,
		outputDir:   outputDir,
		dryRun:      dryRun,
	}
}

// FetchPublicKey fetches the RS256 public key from Keycloak's JWKS endpoint.
// Exported so the setup runner can call it inline without duplicating logic.
func FetchPublicKey(ctx context.Context, keycloakNodePort int, realm string) (string, error) {
	jwksURL := fmt.Sprintf("http://localhost:%d/realms/%s/protocol/openid-connect/certs",
		keycloakNodePort, realm)
	return util.FetchRS256PEM(ctx, jwksURL)
}

// Run fetches the Keycloak RS256 public key (from --key or JWKS endpoint),
// saves it to secrets.yaml, regenerates values-secrets.yaml, and
// redeploys jeeb-infra so Kong picks up the new key.
func (u *Updater) Run(ctx context.Context, pemKey string) error {
	if pemKey == "" {
		jwksURL := fmt.Sprintf("http://localhost:%d/realms/%s/protocol/openid-connect/certs",
			u.cfg.KeycloakNodePort, u.cfg.KeycloakRealm)
		var err error
		pemKey, err = util.FetchRS256PEM(ctx, jwksURL)
		if err != nil {
			return fmt.Errorf("fetch key from Keycloak: %w\n\nIs Keycloak running at http://localhost:%d?",
				err, u.cfg.KeycloakNodePort)
		}
		logger.Step("Fetched RS256 public key from Keycloak.")
	}

	if !strings.Contains(pemKey, "BEGIN PUBLIC KEY") {
		return fmt.Errorf("provided key does not look like a PEM public key (missing 'BEGIN PUBLIC KEY')")
	}

	logger.Step("Updating secrets.yaml...")
	creds, err := u.updateSecretsFile(pemKey)
	if err != nil {
		return err
	}

	logger.Step("Regenerating values-secrets.yaml...")
	secretsPath, err := setup.WriteSecretsValuesFile(u.outputDir, creds, u.cfg.NexusRegistry)
	if err != nil {
		return err
	}

	logger.Step("Redeploying jeeb-infra with updated Kong key...")
	runner := setup.NewRunner(u.cfg, creds, u.chartsDir, u.outputDir, u.secretsPath, u.dryRun)
	runner.SetSecretsFile(secretsPath)
	return runner.DeployInfra(ctx)
}

// UpdateCredsFile writes the PEM key into secrets.yaml and reloads credentials.
// Exported so the setup runner can update creds inline without re-deploying infra itself.
func (u *Updater) UpdateCredsFile(pemKey string) (*credentials.Credentials, error) {
	return u.updateSecretsFile(pemKey)
}

func (u *Updater) updateSecretsFile(pemKey string) (*credentials.Credentials, error) {
	data, err := os.ReadFile(u.secretsPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", u.secretsPath, err)
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse %s: %w", u.secretsPath, err)
	}

	kong, ok := raw["kong"].(map[string]interface{})
	if !ok {
		kong = map[string]interface{}{}
		raw["kong"] = kong
	}
	kong["keycloakPublicKey"] = pemKey

	updated, err := yaml.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal credentials: %w", err)
	}

	if u.dryRun {
		logger.Step("      [dry-run] would write updated kong.keycloakPublicKey to %s", u.secretsPath)
	} else {
		if err := os.WriteFile(u.secretsPath, updated, 0600); err != nil {
			return nil, fmt.Errorf("write %s: %w", u.secretsPath, err)
		}
	}

	return credentials.Load(u.secretsPath)
}
