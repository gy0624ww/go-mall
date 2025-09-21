package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-study-lab/go-mall/common/errcode"
	"github.com/go-study-lab/go-mall/common/logger"
	"github.com/go-study-lab/go-mall/common/middleware"
	"github.com/go-study-lab/go-mall/config"
)

func main() {
	logger.Init()
	r := gin.Default()
	r.Use(middleware.StartTrace(), middleware.LogAccess(), middleware.GinPanicRecovery())
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
		var a map[string]string
		a["k"] = "v"
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data":   a,
		})
	})
	r.GET("/customized-error-test", func(c *gin.Context) {
		// 使用wrap包装原因error 生成项目error
		err := errors.New("a dao error")
		appErr := errcode.Wrap("包装错误", err)
		bAppErr := errcode.Wrap("再包装错误", appErr)
		logger.Error(c, "记录错误", "err", bAppErr)
		// 预定义的ErrServer, 给其追加错误原因的error
		err = errors.New("a domain error")
		apiErr := errcode.ErrServer.WithCause(err)
		logger.Error(c, "API 执行中出现错误", "err", apiErr)
		c.JSON(apiErr.HttpStatusCode(), gin.H{
			"code": apiErr.Code(),
			"msg":  apiErr.Msg(),
		})
	})

	r.Run(":8080")
}
