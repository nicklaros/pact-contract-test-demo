package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Name string `json:"name"`
}

type Inventory struct {
	Stock int32 `json:"stock"`
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		productServiceURL := "http://localhost:8082"
		inventoryServiceURL := "http://localhost:8083"

		var productJsonResp Product
		resp, _ := http.Get(productServiceURL)
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &productJsonResp)
		resp.Body.Close()

		var inventoryJsonResp Inventory
		resp, _ = http.Get(inventoryServiceURL)
		body, _ = io.ReadAll(resp.Body)
		json.Unmarshal(body, &inventoryJsonResp)
		resp.Body.Close()

		c.JSON(http.StatusOK, gin.H{
			"product_service_url":   productServiceURL,
			"inventory_service_url": inventoryServiceURL,
			"product": gin.H{
				"name":  productJsonResp.Name,
				"stock": inventoryJsonResp.Stock,
			},
		})
	})

	r.Run("0.0.0.0:8081")
}
