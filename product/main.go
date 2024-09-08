package main

import (
	"fmt"
	"net/http"

	"pact-contract-test-demo/common"

	"github.com/gin-gonic/gin"
)

func main() {
	runService(common.GetPortFromEnvVar(8082))
}

func runService(port int) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		id := c.DefaultQuery("id", "BEST")

		product := productById[id]

		if product != nil {
			c.JSON(http.StatusOK, gin.H{
				"product": product,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"product": nil,
			})
		}
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
