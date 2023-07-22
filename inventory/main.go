package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	runService(8083)
}

func runService(port int) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		id := c.DefaultQuery("id", "BEST")

		inventory := inventoryByProductId[id]

		c.JSON(http.StatusOK, inventory)
	})

	r.Run(fmt.Sprintf("0.0.0.0:%d", port))
}

type Inventory struct {
	Stock int32 `json:"stock"`
}

var inventoryByProductId = map[string]*Inventory{
	"BEST": {
		Stock: int32(87),
	},
	"TEST_EXISTING_PRODUCT": {
		Stock: int32(101),
	},
}

func addInventory(productId string, stock int32) {
	inventoryByProductId[productId] = &Inventory{
		Stock: stock,
	}
}
