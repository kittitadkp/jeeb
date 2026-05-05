package setup

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"k8s-manager/internal/credentials"

	"gopkg.in/yaml.v3"
)

// GenerateSecretsValues builds a Helm values YAML from credentials.
// nexusRegistry is the host:port used for the Nexus docker pull secret (e.g. "localhost:30050").
func GenerateSecretsValues(creds *credentials.Credentials, nexusRegistry string) ([]byte, error) {
	data := map[string]interface{}{
		"mongodb": map[string]interface{}{
			"credentials": map[string]interface{}{
				"username": creds.MongoDBUsername,
				"password": creds.MongoDBPassword,
			},
		},
		"keycloak": map[string]interface{}{
			"credentials": map[string]interface{}{
				"adminUser":     creds.KeycloakAdminUser,
				"adminPassword": creds.KeycloakAdminPassword,
			},
		},
		"nexus": map[string]interface{}{
			"dockerConfigJson": nexusDockerConfigJSON(nexusRegistry, creds.NexusAdminPassword),
			"credentials": map[string]interface{}{
				"adminPassword": creds.NexusAdminPassword,
			},
		},
		"jenkins": map[string]interface{}{
			"credentials": map[string]interface{}{
				"adminPassword": creds.JenkinsAdminPassword,
				"githubUser":    creds.JenkinsGithubUser,
				"githubPat":     creds.JenkinsGithubPAT,
				"nexusUser":     creds.JenkinsNexusUser,
				"nexusPat":      creds.JenkinsNexusPAT,
				"sonarToken":    creds.JenkinsSonarToken,
			},
		},
		"sonarqube": map[string]interface{}{
			"credentials": map[string]interface{}{
				"adminPassword": creds.SonarQubeAdminPassword,
			},
		},
		"kong": map[string]interface{}{
			"keycloakPublicKey": creds.KongKeycloakPublicKey,
		},
		"grafana": map[string]interface{}{
			"adminPassword": creds.GrafanaAdminPassword,
		},
	}
	return yaml.Marshal(data)
}

// WriteSecretsValuesFile writes the generated secrets values to outputDir/values-secrets.yaml.
// The file is created with mode 0600 so only the owner can read it.
func WriteSecretsValuesFile(outputDir string, creds *credentials.Credentials, nexusRegistry string) (string, error) {
	content, err := GenerateSecretsValues(creds, nexusRegistry)
	if err != nil {
		return "", fmt.Errorf("generate secrets values: %w", err)
	}
	path := filepath.Join(outputDir, "values-secrets.yaml")
	if err := os.WriteFile(path, content, 0600); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	return path, nil
}

// nexusDockerConfigJSON produces the base64-encoded dockerconfigjson for Nexus.
// Format: base64({"auths":{"<registry>":{"auth":"base64(admin:<password>)"}}})
func nexusDockerConfigJSON(registry, adminPassword string) string {
	auth := base64.StdEncoding.EncodeToString([]byte("admin:" + adminPassword))
	cfg := map[string]interface{}{
		"auths": map[string]interface{}{
			registry: map[string]interface{}{
				"auth": auth,
			},
		},
	}
	j, _ := json.Marshal(cfg)
	return base64.StdEncoding.EncodeToString(j)
}
