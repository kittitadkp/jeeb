package util

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

// NewBasicAuthClient returns an HTTP client with a cookie jar (required for Jenkins CSRF)
// and a 30s timeout.
func NewBasicAuthClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	return &http.Client{Timeout: 30 * time.Second, Jar: jar}
}

// DoJSON performs a GET to urlStr and JSON-decodes the response body into out.
func DoJSON(ctx context.Context, client *http.Client, urlStr string, out any) error {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d from %s", resp.StatusCode, urlStr)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// FetchRS256PEM fetches the JWKS from jwksURL and returns the first RSA key as a PEM public key.
func FetchRS256PEM(ctx context.Context, jwksURL string) (string, error) {
	var jwks struct {
		Keys []struct {
			Kty string `json:"kty"`
			Alg string `json:"alg"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := DoJSON(ctx, nil, jwksURL, &jwks); err != nil {
		return "", fmt.Errorf("fetch JWKS from %s: %w", jwksURL, err)
	}
	for _, key := range jwks.Keys {
		if key.Kty == "RSA" && (key.Alg == "RS256" || key.Alg == "") {
			return JWKToPEM(key.N, key.E)
		}
	}
	return "", fmt.Errorf("no RSA key found in JWKS at %s", jwksURL)
}

// JWKToPEM converts a JWK RSA modulus+exponent (base64url) to a PEM-encoded public key.
func JWKToPEM(nB64, eB64 string) (string, error) {
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
