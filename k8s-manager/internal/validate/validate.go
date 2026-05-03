package validate

import (
	"fmt"
	"strings"

	"k8s-manager/internal/credentials"
)

const (
	checkMark = "[✓]"
	crossMark = "[✗]"
)

// Run prints a pass/fail table for all credential fields and returns false if
// any required field is missing.
func Run(creds *credentials.Credentials) bool {
	fmt.Println("=== Credential Validation ===")
	fmt.Println()

	allOK := true

	fmt.Println("Required fields:")
	for name, val := range creds.RequiredFields() {
		if strings.TrimSpace(val) == "" {
			fmt.Printf("  %s %-30s not set\n", crossMark, name)
			allOK = false
		} else {
			fmt.Printf("  %s %s\n", checkMark, name)
		}
	}

	fmt.Println()
	fmt.Println("Optional fields (can be filled later):")
	for name, val := range creds.OptionalFields() {
		if strings.TrimSpace(val) == "" {
			fmt.Printf("  %s %-30s not set\n", crossMark, name)
		} else {
			fmt.Printf("  %s %s\n", checkMark, name)
		}
	}

	fmt.Println()
	if allOK {
		fmt.Println("All required credentials are set. Ready to run setup.")
	} else {
		fmt.Println("Some required credentials are missing. Fill them in credentials.yaml before running setup.")
	}

	if creds.KongKeycloakPublicKey == "" {
		fmt.Println("Tip: kong.keycloakPublicKey can only be filled after Keycloak is running.")
		fmt.Println("     Run 'k8s-manager kong-key' once Keycloak is up at http://localhost:30081")
	}

	return allOK
}
