package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	runService(8082)
}

func runService(port int) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		id := c.DefaultQuery("id", "BEST")

		product := productById[id]

		c.JSON(http.StatusOK, product)
	})

	r.Run(fmt.Sprintf("0.0.0.0:%d", port))
}

type Product struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var productById = map[string]*Product{
	"BEST": {
		Id:   "BEST",
		Name: "The Best Product in The World",
	},
	"TEST_EXISTING_PRODUCT": {
		Id:   "TEST_EXISTING_PRODUCT",
		Name: "The existing product with a wonderful name",
	},
}

func addProduct(id string, name string) {
	productById[id] = &Product{
		Id:   id,
		Name: name,
	}
}

func deleteProduct(id string) {
	productById[id] = nil
}
