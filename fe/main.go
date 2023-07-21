package main

import (
	"fmt"
	"net/http"
	"pact-contract-test-demo/common"

	"github.com/gin-gonic/gin"
)

func main() {
	runService(8080)
}

func runService(port int) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		gatewayServiceURL := "http://localhost:8081"
		resp := callGatewayService(gatewayServiceURL)

		c.JSON(http.StatusOK, gin.H{
			"gateway_service_url": gatewayServiceURL,
			"product": gin.H{
				"name":  resp.Product.Name,
				"stock": resp.Product.Stock,
			},
		})
	})

	r.Run(fmt.Sprintf("0.0.0.0:%d", port))
}

type Resp struct {
	ProductServiceURL   string  `json:"product_service_url"`
	InventoryServiceURL string  `json:"inventory_service_url"`
	Product             Product `json:"product"`
}

type Product struct {
	Name  string `json:"name"`
	Stock int32  `json:"stock"`
}

func callGatewayService(url string) Resp {
	var jsonResp Resp
	common.CallService(url, &jsonResp)

	return jsonResp
}
