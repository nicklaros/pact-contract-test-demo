package main

import (
	"fmt"
	"testing"

	"pact-contract-test-demo/common/errors"

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

	for _, testCase := range testCases {
		testCase(t, pact, url)
	}

	pact.Teardown()
}

func getExistingProduct(t *testing.T, pact dsl.Pact, url string) {
	pact.
		AddInteraction().
		UponReceiving("A request to get product with id `TEST_EXISTING_PRODUCT`").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.String("/"),
			Query: dsl.MapMatcher{
				"id": dsl.String("TEST_EXISTING_PRODUCT"),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body: dsl.StructMatcher{
				"product_service_url":   "http://localhost:8082",
				"inventory_service_url": "http://localhost:8083",
				"error":                 nil,
				"product": dsl.StructMatcher{
					"name":  "The existing product with a wonderful name",
					"stock": 87,
				},
			},
		})

	err := pact.Verify(func() error {
		resp := callGatewayService(url, "TEST_EXISTING_PRODUCT")

		assert.Equal(t, Resp{
			ProductServiceURL:   "http://localhost:8082",
			InventoryServiceURL: "http://localhost:8083",
			Error:               nil,
			Product: &Product{
				Name:  "The existing product with a wonderful name",
				Stock: 87,
			},
		}, resp)

		return nil
	})

	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

func getNonExistingProduct(t *testing.T, pact dsl.Pact, url string) {
	pact.
		AddInteraction().
		UponReceiving("A request to get product with id `TEST_NONEXISTING_PRODUCT`").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.String("/"),
			Query: dsl.MapMatcher{
				"id": dsl.String("TEST_NONEXISTING_PRODUCT"),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 404,
			Body: dsl.StructMatcher{
				"product_service_url":   "http://localhost:8082",
				"inventory_service_url": "http://localhost:8083",
				"error": dsl.StructMatcher{
					"code":    "C001",
					"message": "product not found",
				},
				"product": nil,
			},
		})

	err := pact.Verify(func() error {
		resp := callGatewayService(url, "TEST_NONEXISTING_PRODUCT")

		assert.Equal(t, Resp{
			ProductServiceURL:   "http://localhost:8082",
			InventoryServiceURL: "http://localhost:8083",
			Error: &errors.Error{
				Code:    "C001",
				Message: "product not found",
			},
			Product: nil,
		}, resp)

		return nil
	})

	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

type CaseFn func(t *testing.T, pact dsl.Pact, url string)

var testCases = []CaseFn{
	getExistingProduct,
	getNonExistingProduct,
}
