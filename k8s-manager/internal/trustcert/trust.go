package trustcert

import (
	"context"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"k8s-manager/internal/kube"
	"k8s-manager/internal/logger"
)

const (
	ScopeCurrentUser  = "current-user"
	ScopeLocalMachine = "local-machine"
)

type storeSpec struct {
	Name     string
	Location string
	Label    string
}

type Importer struct {
	Kubeconfig string
	Namespace  string
	SecretName string
	CertKey    string
	Scope      string
	DryRun     bool
}

func NewImporter(kubeconfig, namespace, secretName, certKey, scope string, dryRun bool) *Importer {
	return &Importer{
		Kubeconfig: kubeconfig,
		Namespace:  namespace,
		SecretName: secretName,
		CertKey:    certKey,
		Scope:      scope,
		DryRun:     dryRun,
	}
}

func (i *Importer) Run(ctx context.Context) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("trust-cert is only supported on Windows")
	}

	store, err := storeForScope(i.Scope)
	if err != nil {
		return err
	}

	client, err := kube.NewClient(i.Kubeconfig)
	if err != nil {
		return err
	}

	secret, err := client.GetSecret(ctx, i.Namespace, i.SecretName)
	if err != nil {
		return err
	}

	certPEM, cert, err := extractCertificate(secret.Data, i.CertKey)
	if err != nil {
		return err
	}

	thumbprint := certificateThumbprint(cert)
	logger.Step("Certificate: %s", cert.Subject.String())
	logger.Step("Issuer: %s", cert.Issuer.String())
	logger.Step("Expires: %s", cert.NotAfter.Format(time.DateOnly))
	logger.Step("Thumbprint: %s", thumbprint)

	if i.DryRun {
		logger.Step("Would import %s/%s[%s] into %s", i.Namespace, i.SecretName, i.CertKey, store.Label)
		return nil
	}

	present, err := certificatePresent(ctx, store, thumbprint)
	if err != nil {
		return err
	}
	if present {
		logger.Step("Certificate already trusted in %s", store.Label)
		return nil
	}

	tempFile, err := os.CreateTemp("", "jeeb-dev-tls-*.crt")
	if err != nil {
		return fmt.Errorf("create temp cert file: %w", err)
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath)

	if _, err := tempFile.Write(certPEM); err != nil {
		tempFile.Close()
		return fmt.Errorf("write temp cert file: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temp cert file: %w", err)
	}

	if err := importCertificate(ctx, store, tempPath); err != nil {
		return err
	}

	logger.Step("Trusted certificate from %s/%s in %s", i.Namespace, i.SecretName, store.Label)
	logger.StepMsg("Restart Chrome or Edge if the warning persists.")
	return nil
}

func extractCertificate(secretData map[string][]byte, key string) ([]byte, *x509.Certificate, error) {
	certPEM, ok := secretData[key]
	if !ok {
		return nil, nil, fmt.Errorf("secret is missing %q", key)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, nil, fmt.Errorf("secret %q does not contain a PEM certificate", key)
	}
	if block.Type != "CERTIFICATE" {
		return nil, nil, fmt.Errorf("secret %q contains PEM block %q, expected CERTIFICATE", key, block.Type)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse certificate from %q: %w", key, err)
	}

	return certPEM, cert, nil
}

func certificateThumbprint(cert *x509.Certificate) string {
	sum := sha1.Sum(cert.Raw)
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

func storeForScope(scope string) (storeSpec, error) {
	switch scope {
	case "", ScopeCurrentUser:
		return storeSpec{
			Name:     "Root",
			Location: "CurrentUser",
			Label:    `CurrentUser\Root`,
		}, nil
	case ScopeLocalMachine:
		return storeSpec{
			Name:     "Root",
			Location: "LocalMachine",
			Label:    `LocalMachine\Root`,
		}, nil
	default:
		return storeSpec{}, fmt.Errorf("unsupported scope %q (valid: %s, %s)", scope, ScopeCurrentUser, ScopeLocalMachine)
	}
}

func certificatePresent(ctx context.Context, store storeSpec, thumbprint string) (bool, error) {
	script := fmt.Sprintf(
		"$store = New-Object System.Security.Cryptography.X509Certificates.X509Store(%s, %s); "+
			"$store.Open([System.Security.Cryptography.X509Certificates.OpenFlags]::ReadOnly); "+
			"try { $match = $store.Certificates | Where-Object Thumbprint -eq %s; if ($null -ne $match) { Write-Output 'present' } else { Write-Output 'missing' } } "+
			"finally { $store.Close() }",
		psQuote(store.Name),
		psQuote(store.Location),
		psQuote(thumbprint),
	)

	out, err := runPowerShell(ctx, script)
	if err != nil {
		return false, fmt.Errorf("check Windows certificate store: %w", err)
	}

	return strings.Contains(strings.ToLower(string(out)), "present"), nil
}

func importCertificate(ctx context.Context, store storeSpec, certPath string) error {
	script := fmt.Sprintf(
		"$store = New-Object System.Security.Cryptography.X509Certificates.X509Store(%s, %s); "+
			"$store.Open([System.Security.Cryptography.X509Certificates.OpenFlags]::ReadWrite); "+
			"try { $cert = New-Object System.Security.Cryptography.X509Certificates.X509Certificate2(%s); $store.Add($cert) } "+
			"finally { if ($null -ne $cert) { $cert.Dispose() }; $store.Close() }",
		psQuote(store.Name),
		psQuote(store.Location),
		psQuote(certPath),
	)

	if _, err := runPowerShell(ctx, script); err != nil {
		return fmt.Errorf("import certificate into %s: %w", store.Label, err)
	}
	return nil
}

func runPowerShell(ctx context.Context, script string) ([]byte, error) {
	cmd := exec.CommandContext(ctx,
		"powershell",
		"-NoProfile",
		"-NonInteractive",
		"-ExecutionPolicy", "Bypass",
		"-Command", script,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("%s", msg)
	}
	return out, nil
}

func psQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
