package setup

import (
	"strings"
	"testing"

	"k8s-manager/internal/credentials"
)

func TestGenerateSecretsValues(t *testing.T) {
	creds := &credentials.Credentials{
		JenkinsAdminPassword:   "pass1",
		JenkinsGithubUser:      "kittitadkp",
		JenkinsGithubPAT:       "ghp_test",
		JenkinsNexusUser:       "admin",
		JenkinsNexusPAT:        "nexus_pass",
		JenkinsSonarToken:      "squ_abc",
		KeycloakAdminUser:      "admin",
		KeycloakAdminPassword:  "kc_pass",
		MongoDBUsername:        "jeeb",
		MongoDBPassword:        "mongo_pass",
		NexusAdminPassword:     "nx_pass",
		SonarQubeAdminPassword: "sonar_pass",
		GrafanaAdminPassword:   "grafana_pass",
		KongKeycloakPublicKey:  "-----BEGIN PUBLIC KEY-----\nABC123\n-----END PUBLIC KEY-----",
	}

	out, err := GenerateSecretsValues(creds, "localhost:30050")
	if err != nil {
		t.Fatal(err)
	}

	yaml := string(out)
	t.Log(yaml)

	checks := []string{
		"adminPassword: pass1",
		"githubUser: kittitadkp",
		"adminPassword: kc_pass",
		"password: mongo_pass",
		"adminPassword: grafana_pass",
		"dockerConfigJson:",
		"keycloakPublicKey:",
	}
	for _, check := range checks {
		if !strings.Contains(yaml, check) {
			t.Errorf("expected output to contain %q", check)
		}
	}

	// nexus dockerConfigJson must be non-empty
	if strings.Contains(yaml, `dockerConfigJson: ""`) {
		t.Error("dockerConfigJson should not be empty")
	}
}
