package main

import (
	"fmt"
	"net/http"
	"pact-contract-test-demo/common"

	"github.com/gin-gonic/gin"
)

func main() {
	runService(8081)
}

func runService(port int) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		id := c.DefaultQuery("id", "BEST")

		productServiceURL := "http://localhost:8082"
		inventoryServiceURL := "http://localhost:8083"

		product := callProductService(productServiceURL, id)

		if product == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"product_service_url":   productServiceURL,
				"inventory_service_url": inventoryServiceURL,
				"product":               nil,
				"error": gin.H{
					"code":    "C001",
					"message": "product not found",
				},
			})

			return
		}

		inventory := callInventoryService(inventoryServiceURL, id)

		stock := int32(0)
		if inventory != nil {
			stock = inventory.Stock
		}

		c.JSON(http.StatusOK, gin.H{
			"product_service_url":   productServiceURL,
			"inventory_service_url": inventoryServiceURL,
			"product": gin.H{
				"name":  product.Name,
				"stock": stock,
			},
			"error": nil,
		})
	})

	r.Run(fmt.Sprintf("0.0.0.0:%d", port))
}

type Product struct {
	Name string `json:"name"`
}

type Inventory struct {
	Stock int32 `json:"stock"`
}

func callProductService(url string, productId string) *Product {
	var jsonResp *Product
	common.CallService(fmt.Sprintf("%s/?id=%s", url, productId), &jsonResp)

	return jsonResp
}

func callInventoryService(url string, productId string) *Inventory {
	var jsonResp *Inventory
	common.CallService(url, &jsonResp)

	return jsonResp
}
