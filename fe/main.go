package main

import (
	"fmt"
	"net/http"

	"pact-contract-test-demo/common"
	"pact-contract-test-demo/common/errors"

	"github.com/gin-gonic/gin"
)

var gatewayServiceBaseURL = "http://localhost:8081" // adjust port if needed

func main() {
	runService(common.GetPortFromEnvVar(8080))
}

func runService(port int) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		id := c.DefaultQuery("id", "BEST")

		resp := callGatewayService(gatewayServiceBaseURL, id)

		c.JSON(http.StatusOK, gin.H{
			"gateway_service_url": gatewayServiceBaseURL,
			"product": gin.H{
				"name":  resp.Product.Name,
				"stock": resp.Product.Stock,
			},
		})
	})

	r.Run(fmt.Sprintf("0.0.0.0:%d", port))
}

type Resp struct {
	ProductServiceURL   string        `json:"product_service_url"`
	InventoryServiceURL string        `json:"inventory_service_url"`
	Product             *Product      `json:"product"`
	Error               *errors.Error `json:"error"`
}

type Product struct {
	Name  string `json:"name"`
	Stock int32  `json:"stock"`
}

func callGatewayService(url string, productId string) Resp {
	var jsonResp Resp
	common.CallService(fmt.Sprintf("%s/?id=%s", url, productId), &jsonResp)

	return jsonResp
}
