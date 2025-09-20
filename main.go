package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-study-lab/go-mall/common/logger"
	"github.com/go-study-lab/go-mall/common/middleware"
	"github.com/go-study-lab/go-mall/config"
)

func main() {
	logger.Init()
	r := gin.Default()
	r.Use(gin.Logger(), middleware.StartTrace())
	r.GET("/GET", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	r.GET("/config-read", func(c *gin.Context) {
		database := config.Database
		logger.ZapLoggerTest()
		c.JSON(200, gin.H{
			"type":     database.Type,
			"max_life": database.MaxLifeTime,
		})
	})
	r.GET("/logger-test", func(c *gin.Context) {
		logger.Info(c, "logger test", "key", "keyName", "val", 3)
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.Run(":8080")
}
