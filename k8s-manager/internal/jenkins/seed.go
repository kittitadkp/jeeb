package jenkins

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"k8s-manager/internal/config"
	"k8s-manager/internal/credentials"
)

// Seeder creates and runs the Jenkins seed job that generates all pipeline jobs.
type Seeder struct {
	jenkinsURL string
	groovyPath string
	user       string
	password   string
	dryRun     bool
	client     *http.Client

	// values substituted into seed.groovy placeholders and set as global env vars
	githubCredsId       string
	jenkinsRepo         string
	k8sRepo             string
	backendRepo         string
	frontendRepo        string
	learningBackendRepo  string
	learningFrontendRepo string
	nexusRegistryHost   string
	sonarQubeURL        string
}

func NewSeeder(cfg *config.ClusterConfig, creds *credentials.Credentials, groovyPath, jenkinsURL string, dryRun bool) *Seeder {
	url := jenkinsURL
	if url == "" {
		url = fmt.Sprintf("http://localhost:%d", cfg.JenkinsNodePort)
	}
	return &Seeder{
		jenkinsURL: strings.TrimRight(url, "/"),
		groovyPath: groovyPath,
		user:       "admin",
		password:   creds.JenkinsAdminPassword,
		dryRun:     dryRun,
		client:     &http.Client{Timeout: 30 * time.Second},

		githubCredsId:       creds.JenkinsGithubCredsId,
		jenkinsRepo:         creds.JenkinsJenkinsRepo,
		k8sRepo:             creds.JenkinsK8sRepo,
		backendRepo:         creds.JenkinsBackendRepo,
		frontendRepo:        creds.JenkinsFrontendRepo,
		learningBackendRepo:  creds.JenkinsLearningBackendRepo,
		learningFrontendRepo: creds.JenkinsLearningFrontendRepo,
		nexusRegistryHost:   cfg.NexusRegistryHost,
		sonarQubeURL:        cfg.SonarQubeURL,
	}
}

func (s *Seeder) Run(ctx context.Context) error {
	fmt.Printf("      waiting for Jenkins at %s ...\n", s.jenkinsURL)
	if err := s.waitForJenkins(ctx); err != nil {
		return fmt.Errorf("wait for Jenkins: %w", err)
	}
	fmt.Println("      Jenkins is ready")

	if s.dryRun {
		fmt.Printf("      [dry-run] would read %s, create seed job, trigger build\n", s.groovyPath)
		return nil
	}

	groovy, err := os.ReadFile(s.groovyPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", s.groovyPath, err)
	}

	crumbField, crumb, err := s.getCrumb(ctx)
	if err != nil {
		return fmt.Errorf("get crumb: %w", err)
	}

	fmt.Println("      configuring Jenkins global environment variables ...")
	if err := s.configureGlobalEnvVars(ctx, crumbField, crumb); err != nil {
		return fmt.Errorf("configure global env vars: %w", err)
	}

	fmt.Println("      creating/updating seed job ...")
	script := s.substituteGroovy(string(groovy))
	if err := s.createOrUpdateJob(ctx, crumbField, crumb, script); err != nil {
		return fmt.Errorf("create seed job: %w", err)
	}

	fmt.Println("      triggering seed build ...")
	queueURL, err := s.triggerBuild(ctx, crumbField, crumb)
	if err != nil {
		return fmt.Errorf("trigger seed build: %w", err)
	}

	fmt.Println("      waiting for seed build to complete ...")
	return s.waitForBuild(ctx, queueURL)
}

func (s *Seeder) waitForJenkins(ctx context.Context) error {
	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.jenkinsURL+"/api/json", nil)
		if err != nil {
			return err
		}
		req.SetBasicAuth(s.user, s.password)
		resp, err := s.client.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
	return fmt.Errorf("Jenkins not ready after 3 minutes at %s", s.jenkinsURL)
}

func (s *Seeder) getCrumb(ctx context.Context) (field, value string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.jenkinsURL+"/crumbIssuer/api/json", nil)
	if err != nil {
		return "", "", err
	}
	req.SetBasicAuth(s.user, s.password)
	resp, err := s.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", "", nil // CSRF disabled
	}
	var data struct {
		CrumbRequestField string `json:"crumbRequestField"`
		Crumb             string `json:"crumb"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", "", fmt.Errorf("decode crumb: %w", err)
	}
	return data.CrumbRequestField, data.Crumb, nil
}

func (s *Seeder) createOrUpdateJob(ctx context.Context, crumbField, crumb, groovyContent string) error {
	xmlBody := buildSeedJobXML(groovyContent)

	// Check whether the job already exists to decide create vs update endpoint.
	checkReq, err := http.NewRequestWithContext(ctx, http.MethodGet, s.jenkinsURL+"/job/seed/api/json", nil)
	if err != nil {
		return err
	}
	checkReq.SetBasicAuth(s.user, s.password)
	checkResp, err := s.client.Do(checkReq)
	if err != nil {
		return err
	}
	checkResp.Body.Close()

	endpoint := s.jenkinsURL + "/createItem?name=seed"
	if checkResp.StatusCode == http.StatusOK {
		endpoint = s.jenkinsURL + "/job/seed/config.xml"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(xmlBody))
	if err != nil {
		return err
	}
	req.SetBasicAuth(s.user, s.password)
	req.Header.Set("Content-Type", "application/xml")
	if crumbField != "" {
		req.Header.Set(crumbField, crumb)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Jenkins returned %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (s *Seeder) triggerBuild(ctx context.Context, crumbField, crumb string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.jenkinsURL+"/job/seed/build", nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(s.user, s.password)
	if crumbField != "" {
		req.Header.Set(crumbField, crumb)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("trigger build returned %d: %s", resp.StatusCode, string(body))
	}
	return resp.Header.Get("Location"), nil
}

func (s *Seeder) waitForBuild(ctx context.Context, queueItemURL string) error {
	buildURL := s.jenkinsURL + "/job/seed/lastBuild/api/json"

	// Resolve queue item → actual build number so we track the right build.
	if queueItemURL != "" {
		url := strings.TrimRight(queueItemURL, "/") + "/api/json"
		deadline := time.Now().Add(2 * time.Minute)
		for time.Now().Before(deadline) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return err
			}
			req.SetBasicAuth(s.user, s.password)
			resp, err := s.client.Do(req)
			if err == nil && resp.StatusCode == http.StatusOK {
				var item struct {
					Executable struct {
						Number int `json:"number"`
					} `json:"executable"`
				}
				_ = json.NewDecoder(resp.Body).Decode(&item)
				resp.Body.Close()
				if item.Executable.Number > 0 {
					buildURL = fmt.Sprintf("%s/job/seed/%d/api/json", s.jenkinsURL, item.Executable.Number)
					break
				}
			} else if resp != nil {
				resp.Body.Close()
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(3 * time.Second):
			}
		}
	}

	deadline := time.Now().Add(5 * time.Minute)
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, buildURL, nil)
		if err != nil {
			return err
		}
		req.SetBasicAuth(s.user, s.password)
		resp, err := s.client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			var build struct {
				Building bool   `json:"building"`
				Result   string `json:"result"`
			}
			_ = json.NewDecoder(resp.Body).Decode(&build)
			resp.Body.Close()
			if !build.Building && build.Result != "" {
				if build.Result == "SUCCESS" {
					fmt.Println("      seed complete — pipeline jobs created")
					return nil
				}
				return fmt.Errorf("seed build result: %s (check %s/job/seed/lastBuild/console)", build.Result, s.jenkinsURL)
			}
		} else if resp != nil {
			resp.Body.Close()
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
	return fmt.Errorf("seed build did not complete within 5 minutes")
}

func (s *Seeder) substituteGroovy(content string) string {
	r := strings.NewReplacer(
		"@@GITHUB_CREDS_ID@@", s.githubCredsId,
		"@@JENKINS_REPO@@", s.jenkinsRepo,
		"@@BACKEND_REPO@@", s.backendRepo,
		"@@FRONTEND_REPO@@", s.frontendRepo,
		"@@LEARNING_BACKEND_REPO@@", s.learningBackendRepo,
		"@@LEARNING_FRONTEND_REPO@@", s.learningFrontendRepo,
	)
	return r.Replace(content)
}

// configureGlobalEnvVars sets Jenkins global environment variables via the
// Script Console API so that jeebPipeline.groovy can read them at pipeline runtime.
func (s *Seeder) configureGlobalEnvVars(ctx context.Context, crumbField, crumb string) error {
	script := fmt.Sprintf(`
import jenkins.model.Jenkins
import hudson.slaves.EnvironmentVariablesNodeProperty

def instance = Jenkins.getInstance()
def globalProps = instance.getGlobalNodeProperties()
def envProps = globalProps.getAll(EnvironmentVariablesNodeProperty.class)
if (!envProps) {
    globalProps.add(new EnvironmentVariablesNodeProperty())
    envProps = globalProps.getAll(EnvironmentVariablesNodeProperty.class)
}
def envVars = envProps[0].getEnvVars()
envVars.put('NEXUS_REGISTRY_HOST', '%s')
envVars.put('SONAR_URL', '%s')
envVars.put('K8S_REPO_URL', '%s')
envVars.put('GITHUB_CREDS_ID', '%s')
instance.save()
println 'env vars configured'
`, s.nexusRegistryHost, s.sonarQubeURL, s.k8sRepo, s.githubCredsId)

	if s.dryRun {
		fmt.Printf("      [dry-run] would configure global env vars: NEXUS_REGISTRY_HOST=%s SONAR_URL=%s K8S_REPO_URL=%s GITHUB_CREDS_ID=%s\n",
			s.nexusRegistryHost, s.sonarQubeURL, s.k8sRepo, s.githubCredsId)
		return nil
	}

	form := "script=" + netURLQueryEscape(script)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.jenkinsURL+"/scriptText",
		strings.NewReader(form))
	if err != nil {
		return err
	}
	req.SetBasicAuth(s.user, s.password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if crumbField != "" {
		req.Header.Set(crumbField, crumb)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("scriptText returned %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// netURLQueryEscape percent-encodes a string for use as a form value.
func netURLQueryEscape(s string) string {
	// Encode every byte that isn't unreserved per RFC 3986.
	const unreserved = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.~"
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if strings.IndexByte(unreserved, c) >= 0 {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

func buildSeedJobXML(groovyContent string) string {
	// Escape any CDATA end sequence inside the groovy script.
	safe := strings.ReplaceAll(groovyContent, "]]>", "]]]]><![CDATA[>")
	return fmt.Sprintf(`<?xml version='1.1' encoding='UTF-8'?>
<project>
  <description>Seed job — generates all jeeb pipeline jobs from seed.groovy</description>
  <keepDependencies>false</keepDependencies>
  <properties/>
  <scm class="hudson.scm.NullSCM"/>
  <canRoam>true</canRoam>
  <disabled>false</disabled>
  <blockBuildWhenDownstreamBuilding>false</blockBuildWhenDownstreamBuilding>
  <blockBuildWhenUpstreamBuilding>false</blockBuildWhenUpstreamBuilding>
  <triggers/>
  <concurrentBuild>false</concurrentBuild>
  <builders>
    <javaposse.jobdsl.plugin.ExecuteDslScripts plugin="job-dsl">
      <scriptText><![CDATA[%s]]></scriptText>
      <usingScriptText>true</usingScriptText>
      <sandbox>false</sandbox>
      <ignoreMissingFiles>false</ignoreMissingFiles>
      <failOnMissingPlugin>false</failOnMissingPlugin>
      <unstableOnDeprecation>false</unstableOnDeprecation>
      <removedJobAction>IGNORE</removedJobAction>
      <removedViewAction>IGNORE</removedViewAction>
      <removedConfigFilesAction>IGNORE</removedConfigFilesAction>
      <lookupStrategy>JENKINS_ROOT</lookupStrategy>
    </javaposse.jobdsl.plugin.ExecuteDslScripts>
  </builders>
  <publishers/>
  <buildWrappers/>
</project>`, safe)
}
