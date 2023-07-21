package main

import (
	"fmt"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestGatewayConsumer(t *testing.T) {
	pact := dsl.Pact{
		Consumer: "fe_service",
		Provider: "gateway_service",
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
				"product_service_url":   "http://localhost:8082",
				"inventory_service_url": "http://localhost:8083",
				"product": dsl.StructMatcher{
					"name":  "The Best Product in The World",
					"stock": 87,
				},
			},
		})

	err := pact.Verify(func() error {
		resp := callGatewayService(url)

		assert.Equal(t, Resp{
			ProductServiceURL:   "http://localhost:8082",
			InventoryServiceURL: "http://localhost:8083",
			Product: Product{
				Name:  "The Best Product in The World",
				Stock: 87,
			},
		}, resp)

		return nil
	})

	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}

	pact.Teardown()
}
