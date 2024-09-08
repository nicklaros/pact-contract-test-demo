package main

import (
	"fmt"
	"net/http"

	"pact-contract-test-demo/common"

	"github.com/gin-gonic/gin"
)

var productServiceBaseURL = "http://localhost:8082"   // adjust port if needed
var inventoryServiceBaseURL = "http://localhost:8083" // adjust port if needed

func main() {
	runService(common.GetPortFromEnvVar(8081))
}

func runService(port int) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		id := c.DefaultQuery("id", "BEST")

		product := callProductService(productServiceBaseURL, id)

		if product == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"product_service_url":   common.GetServiceFullUrl(productServiceBaseURL, id),
				"inventory_service_url": common.GetServiceFullUrl(inventoryServiceBaseURL, id),
				"product":               nil,
				"error": gin.H{
					"code":    "C001",
					"message": "product not found",
				},
			})

			return
		}

		inventory := callInventoryService(inventoryServiceBaseURL, id)

		stock := int32(0)
		if inventory != nil {
			stock = inventory.Stock
		}

		c.JSON(http.StatusOK, gin.H{
			"product_service_url":   common.GetServiceFullUrl(productServiceBaseURL, id),
			"inventory_service_url": common.GetServiceFullUrl(inventoryServiceBaseURL, id),
			"product": gin.H{
				"name":  product.Name,
				"stock": stock,
			},
			"error": nil,
		})
	})

	r.Run(fmt.Sprintf("0.0.0.0:%d", port))
}

type ProductResp struct {
	Product *Product `json:"product"`
}

type Product struct {
	Name string `json:"name"`
}

type Inventory struct {
	Stock int32 `json:"stock"`
}

func callProductService(url string, productId string) *Product {
	var jsonResp ProductResp
	common.CallService(common.GetServiceFullUrl(url, productId), &jsonResp)

	return jsonResp.Product
}

func callInventoryService(url string, productId string) *Inventory {
	var jsonResp *Inventory
	common.CallService(common.GetServiceFullUrl(url, productId), &jsonResp)

	return jsonResp
}
