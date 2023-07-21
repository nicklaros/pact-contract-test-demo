package main

import (
	"fmt"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestProductConsumer(t *testing.T) {
	pact := dsl.Pact{
		Consumer: "gateway_service",
		Provider: "product_service",
	}

	pact.Setup(true)

	url := fmt.Sprintf("http://localhost:%d", pact.Server.Port)

	pact.
		AddInteraction().
		Given("Product BEST exists").
		UponReceiving("A request to get product").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.String("/"),
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body: dsl.StructMatcher{
				"name": "The Best Product in The World",
			},
		})

	err := pact.Verify(func() error {
		resp := callProductService(url)

		assert.Equal(t, Product{
			Name: "The Best Product in The World",
		}, resp)

		return nil
	})

	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}

	pact.Teardown()
}

func TestInventoryConsumer(t *testing.T) {
	pact := dsl.Pact{
		Consumer: "gateway_service",
		Provider: "inventory_service",
	}

	pact.Setup(true)

	url := fmt.Sprintf("http://localhost:%d", pact.Server.Port)

	pact.
		AddInteraction().
		Given("Inventory BEST exists").
		UponReceiving("A request to get inventory").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.String("/"),
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body: dsl.StructMatcher{
				"stock": 87,
			},
		})

	err := pact.Verify(func() error {
		resp := callInventoryService(url)

		assert.Equal(t, Inventory{
			Stock: 87,
		}, resp)

		return nil
	})

	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}

	pact.Teardown()
}
