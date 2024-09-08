package main

import (
	"fmt"
	"testing"

	"pact-contract-test-demo/common/errors"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

type CaseFn func(t *testing.T, pact *consumer.V2HTTPMockProvider)

var testCases = []CaseFn{
	getExistingProduct,
	getNonExistingProduct,
}

func TestGatewayConsumer(t *testing.T) {
	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "fe_service",
		Provider: "gateway_service",
	})
	assert.NoError(t, err)

	for _, testCase := range testCases {
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
				"error": nil,
				"product": matchers.StructMatcher{
					"name":  "The existing product with a wonderful name",
					"stock": 101,
				},
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			baseURL := fmt.Sprintf("http://%s:%d", config.Host, config.Port)

			resp := callGatewayService(baseURL, "TEST_EXISTING_PRODUCT")

			assert.Equal(t, Resp{
				Error: nil,
				Product: &Product{
					Name:  "The existing product with a wonderful name",
					Stock: 101,
				},
			}, resp)

			return nil
		})

	assert.NoError(t, err)
}

func getNonExistingProduct(t *testing.T, pact *consumer.V2HTTPMockProvider) {
	err := pact.
		AddInteraction().
		Given("Product with id `TEST_NONEXISTING_PRODUCT` not exists").
		UponReceiving("A request to get product with id `TEST_NONEXISTING_PRODUCT`").
		WithRequest("GET", "/", func(request *consumer.V2RequestBuilder) {
			request.Query("id", matchers.String("TEST_NONEXISTING_PRODUCT"))
		}).
		WillRespondWith(404, func(resp *consumer.V2ResponseBuilder) {
			resp.JSONBody(matchers.StructMatcher{
				"error": matchers.StructMatcher{
					"code":    "C001",
					"message": "product not found",
				},
				"product": nil,
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			baseURL := fmt.Sprintf("http://%s:%d", config.Host, config.Port)

			resp := callGatewayService(baseURL, "TEST_NONEXISTING_PRODUCT")

			assert.Equal(t, Resp{
				Error: &errors.Error{
					Code:    "C001",
					Message: "product not found",
				},
				Product: nil,
			}, resp)

			return nil
		})

	assert.NoError(t, err)
}
