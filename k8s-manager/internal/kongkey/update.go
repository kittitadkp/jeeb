package kongkey

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"

	"k8s-manager/internal/config"
	"k8s-manager/internal/credentials"
	"k8s-manager/internal/setup"

	"gopkg.in/yaml.v3"
)

// Updater handles fetching, storing, and deploying the Keycloak public key.
type Updater struct {
	cfg       *config.ClusterConfig
	credsPath string
	chartsDir string
	outputDir string
	dryRun    bool
}

func NewUpdater(cfg *config.ClusterConfig, credsPath, chartsDir, outputDir string, dryRun bool) *Updater {
	return &Updater{
		cfg:       cfg,
		credsPath: credsPath,
		chartsDir: chartsDir,
		outputDir: outputDir,
		dryRun:    dryRun,
	}
}

// Run fetches the Keycloak RS256 public key (from --key or JWKS endpoint),
// saves it to credentials.yaml, regenerates values-secrets.yaml, and
// redeploys jeeb-infra so Kong picks up the new key.
func (u *Updater) Run(ctx context.Context, pemKey string) error {
	if pemKey == "" {
		var err error
		pemKey, err = u.fetchFromKeycloak()
		if err != nil {
			return fmt.Errorf("fetch key from Keycloak: %w\n\nIs Keycloak running at http://localhost:%d?",
				err, u.cfg.KeycloakNodePort)
		}
		fmt.Printf("Fetched RS256 public key from Keycloak.\n")
	}

	if !strings.Contains(pemKey, "BEGIN PUBLIC KEY") {
		return fmt.Errorf("provided key does not look like a PEM public key (missing 'BEGIN PUBLIC KEY')")
	}

	fmt.Println("Updating credentials.yaml...")
	creds, err := u.updateCredsFile(pemKey)
	if err != nil {
		return err
	}

	fmt.Println("Regenerating values-secrets.yaml...")
	secretsPath, err := setup.WriteSecretsValuesFile(u.outputDir, creds, u.cfg.NexusRegistry)
	if err != nil {
		return err
	}

	fmt.Printf("Redeploying jeeb-infra with updated Kong key...\n")
	runner := setup.NewRunner(u.cfg, creds, u.chartsDir, u.outputDir, u.dryRun)
	runner.SetSecretsFile(secretsPath)
	return runner.DeployInfra(ctx)
}

func (u *Updater) updateCredsFile(pemKey string) (*credentials.Credentials, error) {
	data, err := os.ReadFile(u.credsPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", u.credsPath, err)
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse %s: %w", u.credsPath, err)
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
		fmt.Printf("      [dry-run] would write updated kong.keycloakPublicKey to %s\n", u.credsPath)
	} else {
		if err := os.WriteFile(u.credsPath, updated, 0600); err != nil {
			return nil, fmt.Errorf("write %s: %w", u.credsPath, err)
		}
	}

	return credentials.Load(u.credsPath)
}

func (u *Updater) fetchFromKeycloak() (string, error) {
	url := fmt.Sprintf("http://localhost:%d/realms/%s/protocol/openid-connect/certs",
		u.cfg.KeycloakNodePort, u.cfg.KeycloakRealm)
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []struct {
			Kty string `json:"kty"`
			Alg string `json:"alg"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return "", fmt.Errorf("decode JWKS: %w", err)
	}

	for _, key := range jwks.Keys {
		if key.Kty == "RSA" && (key.Alg == "RS256" || key.Alg == "") {
			return jwkToPEM(key.N, key.E)
		}
	}
	return "", fmt.Errorf("no RSA key found in JWKS response")
}

func jwkToPEM(nB64, eB64 string) (string, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nB64)
	if err != nil {
		return "", fmt.Errorf("decode modulus: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eB64)
	if err != nil {
		return "", fmt.Errorf("decode exponent: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := int(new(big.Int).SetBytes(eBytes).Int64())

	pub := &rsa.PublicKey{N: n, E: e}
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", fmt.Errorf("marshal public key: %w", err)
	}

	var buf strings.Builder
	if err := pem.Encode(&buf, &pem.Block{Type: "PUBLIC KEY", Bytes: der}); err != nil {
		return "", fmt.Errorf("encode PEM: %w", err)
	}
	return buf.String(), nil
}
