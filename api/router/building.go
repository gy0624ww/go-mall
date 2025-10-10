package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-study-lab/go-mall/api/controller"
	"github.com/go-study-lab/go-mall/common/middleware"
)

func registerBuildingRoutes(rg *gin.RouterGroup) {
	// 这个路由组中的路由都以/building 开头
	g := rg.Group("/building/")
	// 测试ping
	g.GET("/ping", controller.TestPing)
	// 测试日志文件的读取
	g.GET("config-read", controller.TestConfigRead)
	// 测试日志门面Logger的使用
	g.GET("logger-test", controller.TestLogger)
	// 测试服务的访问日志
	g.POST("access-log-test", controller.TestAccessLog)
	// 测试统一响应--返回列表和分页
	g.GET("response-list", controller.TestResponseList)
	// 测试gorm的日志
	g.GET("gorm-logger-test", controller.TestGormLogger)
	g.POST("create-demo-order", controller.TestCreateDemoOrder)
	// 测试封装的httptool
	g.GET("httptool-get-test", controller.TestForHttpToolGet)
	g.GET("token-make-test", controller.TestMakeToken)
	g.GET("token-auth-test", middleware.AuthUser(), controller.TestAuthToken)
}
