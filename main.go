package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-study-lab/go-mall/config"
)

func main() {
	r := gin.Default()
	r.GET("/GET", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	r.GET("/config-read", func(c *gin.Context) {
		database := config.Database
		c.JSON(200, gin.H{
			"type":     database.Type,
			"max_life": database.MaxLifeTime,
		})
	})
	r.Run(":8080")
}
