package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/pact-foundation/pact-go/utils"
	"github.com/pact-foundation/pact-go/v2/models"
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
		Provider:        "inventory_service",
		ProviderVersion: "1.0.0",
		PactFiles: []string{
			// Use local file for simplicity, in real scenario it is usually fetched from pact-broker.
			filepath.FromSlash("../gateway/pacts/gateway_service-inventory_service.json"),
		},
		StateHandlers: models.StateHandlers{
			"Inventory with product id `TEST_EXISTING_PRODUCT` exists": func(setup bool, s models.ProviderState) (models.ProviderStateResponse, error) {
				// ... do something, such as insert "TEST_EXISTING_PRODUCT" in the database
				addInventory("TEST_EXISTING_PRODUCT", 101)

				// Optionally (if there are generators in the pact) return provider state values to be used in the verification
				return nil, nil
			},
		},
		DisableColoredOutput: true,
	})

	assert.NoError(t, err)
}
