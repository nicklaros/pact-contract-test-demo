package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/pact-foundation/pact-go/utils"
)

var port, _ = utils.GetFreePort()

// The Provider verification
func TestPactProvider(t *testing.T) {
	go runService(port)

	pact := dsl.Pact{
		Provider: "product_service",
	}

	// Verify the Provider - Tag-based Published Pacts for any known consumers
	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: fmt.Sprintf("http://localhost:%d", port),
		PactURLs: []string{
			// Hard coded for simplicity, in real scenario it is usually fetched from pact-broker.
			filepath.FromSlash("../gateway/pacts/gateway_service-product_service.json"),
		},
		ProviderVersion: "1.0.0",
		PactLogDir:      "./logs",
	})

	if err != nil {
		t.Logf("verify provider: %s", err)
	}
}
