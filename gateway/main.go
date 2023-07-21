package main

import (
	"net/http"
	"pact-contract-test-demo/common"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		productServiceURL := "http://localhost:8082"
		product := callProductService(productServiceURL)

		inventoryServiceURL := "http://localhost:8083"
		inventory := callInventoryService(inventoryServiceURL)

		c.JSON(http.StatusOK, gin.H{
			"product_service_url":   productServiceURL,
			"inventory_service_url": inventoryServiceURL,
			"product": gin.H{
				"name":  product.Name,
				"stock": inventory.Stock,
			},
		})
	})

	r.Run("0.0.0.0:8081")
}

type Product struct {
	Name string `json:"name"`
}

type Inventory struct {
	Stock int32 `json:"stock"`
}

func callProductService(url string) Product {
	var jsonResp Product
	common.CallService(url, &jsonResp)

	return jsonResp
}

func callInventoryService(url string) Inventory {
	var jsonResp Inventory
	common.CallService(url, &jsonResp)

	return jsonResp
}
