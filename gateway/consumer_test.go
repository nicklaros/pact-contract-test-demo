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

	for _, testCase := range productTestCases {
		testCase(t, pact, url)
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

	for _, testCase := range inventoryTestCases {
		testCase(t, pact, url)
	}

	pact.Teardown()
}

func getExistingProduct(t *testing.T, pact dsl.Pact, url string) {
	pact.
		AddInteraction().
		Given("Product with id `TEST_EXISTING_PRODUCT` exists").
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
				"name": "The existing product with a wonderful name",
			},
		})

	err := pact.Verify(func() error {
		resp := callProductService(url, "TEST_EXISTING_PRODUCT")

		assert.Equal(t, &Product{
			Name: "The existing product with a wonderful name",
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
		Given("Product with id `TEST_NONEXISTING_PRODUCT` does not exists").
		UponReceiving("A request to get product with id `TEST_NONEXISTING_PRODUCT`").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.String("/"),
			Query: dsl.MapMatcher{
				"id": dsl.String("TEST_NONEXISTING_PRODUCT"),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body:   nil,
		})

	err := pact.Verify(func() error {
		resp := callProductService(url, "TEST_NONEXISTING_PRODUCT")

		assert.Nil(t, resp)

		return nil
	})

	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

func getExistingInventory(t *testing.T, pact dsl.Pact, url string) {
	pact.
		AddInteraction().
		Given("Inventory with product id `TEST_EXISTING_PRODUCT` exists").
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
		resp := callInventoryService(url, "TEST_EXISTING_PRODUCT")

		assert.Equal(t, &Inventory{
			Stock: 87,
		}, resp)

		return nil
	})

	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

type CaseFn func(t *testing.T, pact dsl.Pact, url string)

var productTestCases = []CaseFn{
	getExistingProduct,
	getNonExistingProduct,
}

var inventoryTestCases = []CaseFn{
	getExistingInventory,
}
