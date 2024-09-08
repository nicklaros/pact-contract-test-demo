package main

import (
	"fmt"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

type CaseFn func(t *testing.T, pact *consumer.V2HTTPMockProvider)

var productTestCases = []CaseFn{
	getExistingProduct,
	getNonExistingProduct,
}

var inventoryTestCases = []CaseFn{
	getExistingInventory,
}

func TestProductConsumer(t *testing.T) {
	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "gateway_service",
		Provider: "product_service",
	})
	assert.NoError(t, err)

	for _, testCase := range productTestCases {
		testCase(t, pact)
	}
}

func TestInventoryConsumer(t *testing.T) {
	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "gateway_service",
		Provider: "inventory_service",
	})
	assert.NoError(t, err)

	for _, testCase := range inventoryTestCases {
		testCase(t, pact)
	}
}

func getExistingProduct(t *testing.T, pact *consumer.V2HTTPMockProvider) {
	err := pact.
		AddInteraction().
		Given("Product with id `TEST_EXISTING_PRODUCT` exists").
		UponReceiving("A request to get product with id `TEST_EXISTING_PRODUCT`").
		WithRequest("GET", "/", func(request *consumer.V2RequestBuilder) {
			request.Query("id", matchers.String("TEST_EXISTING_PRODUCT"))
		}).
		WillRespondWith(200, func(resp *consumer.V2ResponseBuilder) {
			resp.JSONBody(matchers.StructMatcher{
				"product": matchers.StructMatcher{
					"name": "The existing product with a wonderful name",
				},
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			baseURL := fmt.Sprintf("http://%s:%d", config.Host, config.Port)

			resp := callProductService(baseURL, "TEST_EXISTING_PRODUCT")

			assert.Equal(t, &Product{
				Name: "The existing product with a wonderful name",
			}, resp)

			return nil
		})

	assert.NoError(t, err)
}

func getNonExistingProduct(t *testing.T, pact *consumer.V2HTTPMockProvider) {
	err := pact.
		AddInteraction().
		Given("Product with id `TEST_NONEXISTING_PRODUCT` does not exists").
		UponReceiving("A request to get product with id `TEST_NONEXISTING_PRODUCT`").
		WithRequest("GET", "/", func(request *consumer.V2RequestBuilder) {
			request.Query("id", matchers.String("TEST_NONEXISTING_PRODUCT"))
		}).
		WillRespondWith(404, func(resp *consumer.V2ResponseBuilder) {
			resp.JSONBody(matchers.StructMatcher{
				"product": nil,
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			baseURL := fmt.Sprintf("http://%s:%d", config.Host, config.Port)

			resp := callProductService(baseURL, "TEST_NONEXISTING_PRODUCT")

			assert.Nil(t, resp)

			return nil
		})

	assert.NoError(t, err)
}

func getExistingInventory(t *testing.T, pact *consumer.V2HTTPMockProvider) {
	err := pact.
		AddInteraction().
		Given("Inventory with product id `TEST_EXISTING_PRODUCT` exists").
		UponReceiving("A request to get inventory").
		WithRequest("GET", "/", func(request *consumer.V2RequestBuilder) {
			request.Query("id", matchers.String("TEST_EXISTING_PRODUCT"))
		}).
		WillRespondWith(200, func(resp *consumer.V2ResponseBuilder) {
			resp.JSONBody(matchers.StructMatcher{
				"stock": 101,
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			baseURL := fmt.Sprintf("http://%s:%d", config.Host, config.Port)

			resp := callInventoryService(baseURL, "TEST_EXISTING_PRODUCT")

			assert.Equal(t, &Inventory{
				Stock: 101,
			}, resp)

			return nil
		})

	assert.NoError(t, err)
}
