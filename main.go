package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-study-lab/go-mall/api/router"
	"github.com/go-study-lab/go-mall/common/enum"
	"github.com/go-study-lab/go-mall/config"
)

func main() {

	if config.App.Env == enum.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.New()
	router.RegisterRoutes(g)
	g.Run(":8080")
}
