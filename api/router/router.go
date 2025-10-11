package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-study-lab/go-mall/common/middleware"
)

func RegisterRoutes(engine *gin.Engine) {
	// use global middleware
	engine.Use(middleware.StartTrace(), middleware.LogAccess(), middleware.GinPanicRecovery())
	routeGroup := engine.Group("")
	registerBuildingRoutes(routeGroup)
	registerUserRoutes(routeGroup)
}
