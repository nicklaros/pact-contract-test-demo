package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/pact-foundation/pact-go/utils"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/pact-foundation/pact-go/v2/version"
	"github.com/stretchr/testify/assert"
)

var port, _ = utils.GetFreePort()

// The Provider verification
func TestPactProvider(t *testing.T) {
	version.CheckVersion()

	go runService(port)

	verifier := provider.NewVerifier()

	err := verifier.VerifyProvider(t, provider.VerifyRequest{
		ProviderBaseURL: fmt.Sprintf("http://localhost:%d", port),
		Provider:        "gateway_service",
		ProviderVersion: "1.0.0",
		PactFiles: []string{
			// Use local file for simplicity, in real scenario it is usually fetched from pact-broker.
			filepath.FromSlash("../fe/pacts/fe_service-gateway_service.json"),
		},
		DisableColoredOutput: true,
	})

	assert.NoError(t, err)
}
