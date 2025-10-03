package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-study-lab/go-mall/api/request"
	"github.com/go-study-lab/go-mall/common/app"
	"github.com/go-study-lab/go-mall/common/errcode"
	"github.com/go-study-lab/go-mall/common/logger"
	"github.com/go-study-lab/go-mall/config"
	"github.com/go-study-lab/go-mall/logic/appservice"
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
		"type":     database.Master.Type,
		"max_life": database.Master.MaxLifeTime,
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

func TestGormLogger(c *gin.Context) {
	svc := appservice.NewDemoAppSvc(c)
	list, err := svc.GetDemoIdentities()
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}
	app.NewResponse(c).Success(list)
	return
}

func TestCreateDemoOrder(c *gin.Context) {
	request := new(request.DemoOrderCreate)
	err := c.ShouldBind(request)
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	// 验证用户信息 Token 然后把UserID赋值上去 这里测试就直接赋值了
	request.UserId = 123453453
	svc := appservice.NewDemoAppSvc(c)
	reply, err := svc.CreateDemoOrder(request)
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}
	app.NewResponse(c).Success(reply)
}
