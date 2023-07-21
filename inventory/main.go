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
		c.JSON(http.StatusOK, gin.H{
			"stock": 87,
		})
	})

	r.Run(fmt.Sprintf("0.0.0.0:%d", port))
}
