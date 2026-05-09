package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ClusterConfig struct {
	NamespaceDev   string
	NamespaceInfra string
	NamespaceObs   string

	ReleaseData     string
	ReleaseDev      string
	ReleaseInfra    string
	ReleaseLearning string
	ReleaseObs      string

	VaultPod  string
	VaultAddr string // in-pod address used via kubectl exec

	KeycloakRealm    string
	KeycloakClientID string
	KongIssuer       string

	MongoHost           string // in-cluster hostname:port
	KeycloakHost        string // in-cluster hostname:port
	KeycloakHostnameURL string // public URL Keycloak advertises as issuer base (e.g. http://auth.jeeb-dev.local)

	// NodePorts — match k8s/charts/*/values.yaml
	MongoNodePort        int // 30017
	FrontendNodePort     int // 30000
	BackendNodePort      int // 30080
	KeycloakNodePort     int // 30081
	JenkinsNodePort      int // 30082
	JenkinsAgentNodePort int // 30500
	NexusUINodePort      int // 30083
	LearningNodePort     int // 30086  (learning-backend)
	LearningFrontPort    int // 30087  (learning-frontend)
	KongNodePort         int // 30088
	SonarQubeNodePort    int // 30090
	VaultNodePort        int // 30091
	GrafanaNodePort      int // 30092
	PrometheusNodePort   int // 30093
	NexusRegistry        string // localhost:30050

	// In-cluster service addresses for Jenkins pipelines
	NexusRegistryHost string // nexus.<infra-ns>.svc.cluster.local:5000
	SonarQubeURL      string // http://sonarqube.<infra-ns>.svc.cluster.local:9000

	IngressLabel string

	// Rancher (optional, deployed via external Helm repos)
	RancherNodePort   int    // 30443
	RancherHostname   string // rancher.jeeb-infra.local
	RancherNamespace  string // cattle-system
	CertManagerNamespace string // cert-manager

	// VaultPaths use the KV v2 API format (secret/data/...) matching Helm chart values.
	// kvCLIPath() converts them to vault CLI format (secret/...) for kv put commands.
	VaultPathBackend          string
	VaultPathFrontend         string
	VaultPathLearningBackend  string
	VaultPathLearningFrontend string
}

// clusterFile mirrors the cluster: section of credentials.yaml.
type clusterFile struct {
	Cluster struct {
		Namespaces struct {
			Dev   string `yaml:"dev"`
			Infra string `yaml:"infra"`
			Obs   string `yaml:"obs"`
		} `yaml:"namespaces"`
		Releases struct {
			Data     string `yaml:"data"`
			Dev      string `yaml:"dev"`
			Infra    string `yaml:"infra"`
			Learning string `yaml:"learning"`
			Obs      string `yaml:"obs"`
		} `yaml:"releases"`
		Vault struct {
			Pod  string `yaml:"pod"`
			Addr string `yaml:"addr"`
		} `yaml:"vault"`
		Keycloak struct {
			Realm       string `yaml:"realm"`
			ClientID    string `yaml:"clientId"`
			Host        string `yaml:"host"`
			HostnameURL string `yaml:"hostnameUrl"`
			NodePort    int    `yaml:"nodePort"`
		} `yaml:"keycloak"`
		Kong struct {
			Issuer   string `yaml:"issuer"`
			NodePort int    `yaml:"nodePort"`
		} `yaml:"kong"`
		Mongo struct {
			Host     string `yaml:"host"`
			NodePort int    `yaml:"nodePort"`
		} `yaml:"mongo"`
		NodePorts struct {
			Frontend        int `yaml:"frontend"`
			Backend         int `yaml:"backend"`
			Jenkins         int `yaml:"jenkins"`
			JenkinsAgent    int `yaml:"jenkinsAgent"`
			NexusUI         int `yaml:"nexusUI"`
			LearningBackend int `yaml:"learningBackend"`
			LearningFront   int `yaml:"learningFront"`
			SonarQube       int `yaml:"sonarQube"`
			Vault           int `yaml:"vault"`
			Grafana         int `yaml:"grafana"`
			Prometheus      int `yaml:"prometheus"`
		} `yaml:"nodePorts"`
		Nexus struct {
			Registry string `yaml:"registry"`
		} `yaml:"nexus"`
		Ingress struct {
			Label string `yaml:"label"`
		} `yaml:"ingress"`
		Rancher struct {
			NodePort             int    `yaml:"nodePort"`
			Hostname             string `yaml:"hostname"`
			Namespace            string `yaml:"namespace"`
			CertManagerNamespace string `yaml:"certManagerNamespace"`
		} `yaml:"rancher"`
		VaultPaths struct {
			Backend          string `yaml:"backend"`
			Frontend         string `yaml:"frontend"`
			LearningBackend  string `yaml:"learningBackend"`
			LearningFrontend string `yaml:"learningFrontend"`
		} `yaml:"vaultPaths"`
	} `yaml:"cluster"`
}

// LoadFromFile reads the cluster: section from a credentials YAML file (if it
// exists) and applies hardcoded defaults for any unset fields. Priority: yaml > default.
func LoadFromFile(path string) *ClusterConfig {
	var f clusterFile
	if path != "" {
		if data, err := os.ReadFile(path); err == nil {
			_ = yaml.Unmarshal(data, &f)
		}
	}
	c := f.Cluster

	return &ClusterConfig{
		NamespaceDev:   coalesce(c.Namespaces.Dev, "jeeb-dev"),
		NamespaceInfra: coalesce(c.Namespaces.Infra, "jeeb-infra"),
		NamespaceObs:   coalesce(c.Namespaces.Obs, "jeeb-obs"),

		ReleaseData:     coalesce(c.Releases.Data, "jeeb-data"),
		ReleaseDev:      coalesce(c.Releases.Dev, "jeeb-dev"),
		ReleaseInfra:    coalesce(c.Releases.Infra, "jeeb-infra"),
		ReleaseLearning: coalesce(c.Releases.Learning, "jeeb-learning"),
		ReleaseObs:      coalesce(c.Releases.Obs, "jeeb-obs"),

		VaultPod:  coalesce(c.Vault.Pod, "vault-0"),
		VaultAddr: coalesce(c.Vault.Addr, "http://127.0.0.1:8200"),

		KeycloakRealm:       coalesce(c.Keycloak.Realm, "jeeb"),
		KeycloakClientID:    coalesce(c.Keycloak.ClientID, "jeeb-app"),
		KongIssuer:          coalesce(c.Kong.Issuer, "http://auth.jeeb-dev.local/realms/jeeb"),
		KeycloakHostnameURL: coalesce(c.Keycloak.HostnameURL, "https://auth.jeeb-dev.local"),

		MongoHost:    coalesce(c.Mongo.Host, "mongodb.jeeb-dev.svc.cluster.local:27017"),
		KeycloakHost: coalesce(c.Keycloak.Host, "keycloak.jeeb-dev.svc.cluster.local:8080"),

		MongoNodePort:        coalesceInt(c.Mongo.NodePort, 30017),
		FrontendNodePort:     coalesceInt(c.NodePorts.Frontend, 30000),
		BackendNodePort:      coalesceInt(c.NodePorts.Backend, 30080),
		KeycloakNodePort:     coalesceInt(c.Keycloak.NodePort, 30081),
		JenkinsNodePort:      coalesceInt(c.NodePorts.Jenkins, 30082),
		JenkinsAgentNodePort: coalesceInt(c.NodePorts.JenkinsAgent, 30500),
		NexusUINodePort:      coalesceInt(c.NodePorts.NexusUI, 30083),
		LearningNodePort:     coalesceInt(c.NodePorts.LearningBackend, 30086),
		LearningFrontPort:    coalesceInt(c.NodePorts.LearningFront, 30087),
		KongNodePort:         coalesceInt(c.Kong.NodePort, 30088),
		SonarQubeNodePort:    coalesceInt(c.NodePorts.SonarQube, 30090),
		VaultNodePort:        coalesceInt(c.NodePorts.Vault, 30091),
		GrafanaNodePort:      coalesceInt(c.NodePorts.Grafana, 30092),
		PrometheusNodePort:   coalesceInt(c.NodePorts.Prometheus, 30093),
		NexusRegistry:        coalesce(c.Nexus.Registry, "localhost:30050"),

		IngressLabel: coalesce(c.Ingress.Label, "app.kubernetes.io/name=ingress-nginx"),

		NexusRegistryHost: "nexus." + coalesce(c.Namespaces.Infra, "jeeb-infra") + ".svc.cluster.local:5000",
		SonarQubeURL:      "http://sonarqube." + coalesce(c.Namespaces.Infra, "jeeb-infra") + ".svc.cluster.local:9000",

		RancherNodePort:      coalesceInt(c.Rancher.NodePort, 30443),
		RancherHostname:      coalesce(c.Rancher.Hostname, "rancher.jeeb-infra.local"),
		RancherNamespace:     coalesce(c.Rancher.Namespace, "cattle-system"),
		CertManagerNamespace: coalesce(c.Rancher.CertManagerNamespace, "cert-manager"),

		// Paths use KV v2 API format (secret/data/...) — matches Helm chart vault.path values.
		VaultPathBackend:          coalesce(c.VaultPaths.Backend, "secret/data/jeeb/backend/develop"),
		VaultPathFrontend:         coalesce(c.VaultPaths.Frontend, "secret/data/jeeb/frontend/develop"),
		VaultPathLearningBackend:  coalesce(c.VaultPaths.LearningBackend, "secret/data/jeeb/learning/backend/develop"),
		VaultPathLearningFrontend: coalesce(c.VaultPaths.LearningFrontend, "secret/data/jeeb/learning/frontend/develop"),
	}
}

// KVCLIPath converts a KV v2 API path ("secret/data/X") to the vault CLI path
// ("secret/X") used by `vault kv put`.
func KVCLIPath(apiPath string) string {
	if after, ok := cutPrefix(apiPath, "secret/data/"); ok {
		return "secret/" + after
	}
	return apiPath
}

func cutPrefix(s, prefix string) (string, bool) {
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):], true
	}
	return s, false
}

func coalesce(v, def string) string {
	if v != "" {
		return v
	}
	return def
}

func coalesceInt(v, def int) int {
	if v != 0 {
		return v
	}
	return def
}
