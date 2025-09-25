package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-study-lab/go-mall/common/app"
	"github.com/go-study-lab/go-mall/common/logger"
	"github.com/go-study-lab/go-mall/config"
)

func TestPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
	return
}

func TestLogger(c *gin.Context) {
	logger.Info(c, "TestLogger", "message", "test logger")
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func TestAccessLog(c *gin.Context) {

}

func TestConfigRead(c *gin.Context) {
	database := config.Database
	c.JSON(http.StatusOK, gin.H{
		"type":     database.Type,
		"max_life": database.MaxLifeTime,
	})
	return
}

func TestResponseList(c *gin.Context) {
	pagination := app.NewPaginaton(c)
	// Mock fetch list data from db
	data := []struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		{
			Name: "Lily",
			Age:  26,
		},
		{
			Name: "Violet",
			Age:  25,
		},
	}
	pagination.SetTotalRows(2)
	app.NewResponse(c).SetPagination(pagination).Success(data)
	return
}
