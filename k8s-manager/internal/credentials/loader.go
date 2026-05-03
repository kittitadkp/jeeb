package credentials

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Credentials holds all service passwords/tokens loaded from credentials.yaml.
type Credentials struct {
	JenkinsAdminPassword string
	JenkinsGithubUser    string
	JenkinsGithubPAT     string
	JenkinsNexusUser     string
	JenkinsNexusPAT      string
	JenkinsSonarToken    string

	KeycloakAdminUser     string
	KeycloakAdminPassword string

	MongoDBUsername string
	MongoDBPassword string

	NexusAdminPassword string

	SonarQubeAdminPassword string

	GrafanaAdminPassword string

	KongKeycloakPublicKey string
}

// credentialsFile mirrors the structure of credentials.yaml.
type credentialsFile struct {
	Jenkins struct {
		AdminPassword string `yaml:"adminPassword"`
		GithubUser    string `yaml:"githubUser"`
		GithubPat     string `yaml:"githubPat"`
		NexusUser     string `yaml:"nexusUser"`
		NexusPat      string `yaml:"nexusPat"`
		SonarToken    string `yaml:"sonarToken"`
	} `yaml:"jenkins"`
	Keycloak struct {
		AdminUser     string `yaml:"adminUser"`
		AdminPassword string `yaml:"adminPassword"`
	} `yaml:"keycloak"`
	MongoDB struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"mongodb"`
	Nexus struct {
		AdminPassword string `yaml:"adminPassword"`
	} `yaml:"nexus"`
	SonarQube struct {
		AdminPassword string `yaml:"adminPassword"`
	} `yaml:"sonarqube"`
	Grafana struct {
		AdminPassword string `yaml:"adminPassword"`
	} `yaml:"grafana"`
	Kong struct {
		KeycloakPublicKey string `yaml:"keycloakPublicKey"`
	} `yaml:"kong"`
}

func Load(path string) (*Credentials, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}

	var f credentialsFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	c := &Credentials{
		JenkinsAdminPassword:   f.Jenkins.AdminPassword,
		JenkinsGithubUser:      f.Jenkins.GithubUser,
		JenkinsGithubPAT:       f.Jenkins.GithubPat,
		JenkinsNexusUser:       f.Jenkins.NexusUser,
		JenkinsNexusPAT:        f.Jenkins.NexusPat,
		JenkinsSonarToken:      f.Jenkins.SonarToken,
		KeycloakAdminUser:      f.Keycloak.AdminUser,
		KeycloakAdminPassword:  f.Keycloak.AdminPassword,
		MongoDBUsername:        f.MongoDB.Username,
		MongoDBPassword:        f.MongoDB.Password,
		NexusAdminPassword:     f.Nexus.AdminPassword,
		SonarQubeAdminPassword: f.SonarQube.AdminPassword,
		GrafanaAdminPassword:   f.Grafana.AdminPassword,
		KongKeycloakPublicKey:  f.Kong.KeycloakPublicKey,
	}

	// Always return c so callers like validate can inspect partial data.
	return c, c.validate()
}

// RequiredFields returns the list of field names and their values used for
// required-field validation. Kong public key is intentionally excluded — it
// can only be filled after Keycloak is running (use 'k8s-manager kong-key').
func (c *Credentials) RequiredFields() map[string]string {
	return map[string]string{
		"jenkins.adminPassword":  c.JenkinsAdminPassword,
		"jenkins.githubUser":     c.JenkinsGithubUser,
		"jenkins.githubPat":      c.JenkinsGithubPAT,
		"jenkins.nexusPat":       c.JenkinsNexusPAT,
		"keycloak.adminPassword": c.KeycloakAdminPassword,
		"mongodb.password":       c.MongoDBPassword,
		"nexus.adminPassword":    c.NexusAdminPassword,
		"grafana.adminPassword":  c.GrafanaAdminPassword,
	}
}

// OptionalFields returns fields that are expected to be empty initially.
func (c *Credentials) OptionalFields() map[string]string {
	return map[string]string{
		"jenkins.sonarToken":     c.JenkinsSonarToken,
		"kong.keycloakPublicKey": c.KongKeycloakPublicKey,
	}
}

func (c *Credentials) validate() error {
	var missing []string
	for k, v := range c.RequiredFields() {
		if strings.TrimSpace(v) == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required credentials: %s", strings.Join(missing, ", "))
	}
	return nil
}
