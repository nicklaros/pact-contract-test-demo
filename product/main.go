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
		c.JSON(http.StatusOK, gin.H{
			"name": "The Best Product in The World",
		})
	})

	r.Run(fmt.Sprintf("0.0.0.0:%d", port))
}
