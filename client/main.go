package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	ProductServiceURL   string  `json:"product_service_url"`
	InventoryServiceURL string  `json:"inventory_service_url"`
	Product             Product `json:"product"`
}

type Product struct {
	Name  string `json:"name"`
	Stock int32  `json:"stock"`
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		gatewayServiceURL := "http://localhost:8081"

		var jsonResp Resp
		resp, _ := http.Get(gatewayServiceURL)
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &jsonResp)
		resp.Body.Close()

		c.JSON(http.StatusOK, gin.H{
			"gateway_service_url": gatewayServiceURL,
			"product": gin.H{
				"name":  jsonResp.Product.Name,
				"stock": jsonResp.Product.Stock,
			},
		})
	})

	r.Run("0.0.0.0:8080")
}
