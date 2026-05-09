package trustcert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

func TestExtractCertificate(t *testing.T) {
	sampleCertPEM := createTestCertificatePEM(t)
	certPEM, cert, err := extractCertificate(map[string][]byte{
		"tls.crt": sampleCertPEM,
	}, "tls.crt")
	if err != nil {
		t.Fatalf("extractCertificate() error = %v", err)
	}
	if got := cert.Subject.CommonName; got != "jeeb-dev.local" {
		t.Fatalf("CommonName = %q, want jeeb-dev.local", got)
	}
	if len(certPEM) == 0 {
		t.Fatal("certPEM should not be empty")
	}
}

func createTestCertificatePEM(t *testing.T) []byte {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "jeeb-dev.local",
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"jeeb-dev.local", "*.jeeb-dev.local"},
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("CreateCertificate() error = %v", err)
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})
}

func TestExtractCertificateMissingKey(t *testing.T) {
	_, _, err := extractCertificate(map[string][]byte{}, "tls.crt")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestStoreForScope(t *testing.T) {
	tests := map[string]storeSpec{
		"": {
			Name:     "Root",
			Location: "CurrentUser",
			Label:    `CurrentUser\Root`,
		},
		ScopeCurrentUser: {
			Name:     "Root",
			Location: "CurrentUser",
			Label:    `CurrentUser\Root`,
		},
		ScopeLocalMachine: {
			Name:     "Root",
			Location: "LocalMachine",
			Label:    `LocalMachine\Root`,
		},
	}

	for scope, want := range tests {
		got, err := storeForScope(scope)
		if err != nil {
			t.Fatalf("storeForScope(%q) error = %v", scope, err)
		}
		if got != want {
			t.Fatalf("storeForScope(%q) = %#v, want %#v", scope, got, want)
		}
	}
}
