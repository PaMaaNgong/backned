package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/auth", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "from auth service"})
	})

	return r
}

func main() {
	r := setupRouter()
	_ = r.Run()
}
