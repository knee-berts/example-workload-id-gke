package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/cloud-storage-bucket", uploadFileHandler)

	r.GET("/cloud-secrets-manager/:name", viewSecretsHandler)

	r.GET("/healthz", func(c *gin.Context) { 
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}