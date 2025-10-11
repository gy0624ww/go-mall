package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-study-lab/go-mall/api/request"
	"github.com/go-study-lab/go-mall/common/app"
	"github.com/go-study-lab/go-mall/common/errcode"
	"github.com/go-study-lab/go-mall/common/logger"
	"github.com/go-study-lab/go-mall/config"
	"github.com/go-study-lab/go-mall/library"
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

func TestForHttpToolGet(c *gin.Context) {
	ipDetail, err := library.NewWhoisLib(c).GetHostIpDetail()
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	app.NewResponse(c).Success(ipDetail)
}

func TestForHttpToolPost(c *gin.Context) {

	orderReply, err := library.NewDemoLib(c).TestPostCreateOrder()
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	app.NewResponse(c).Success(orderReply)
}

func TestMakeToken(c *gin.Context) {
	userSvc := appservice.NewUserAppSvc(c)
	token, err := userSvc.GenToken()
	if err != nil {
		if errors.Is(err, errcode.ErrUserInvalid) {
			logger.Error(c, "invalid user is unable to generate token", err)
			app.NewResponse(c).Error(errcode.ErrUserInvalid)
		} else {
			appErr := err.(*errcode.AppError)
			app.NewResponse(c).Error(appErr)
		}
		return
	}
	app.NewResponse(c).Success(token)
}

func TestAuthToken(c *gin.Context) {
	app.NewResponse(c).Success(gin.H{
		"user_id":    c.GetInt64("userId"),
		"session_id": c.GetString("sessionId"),
	})
	return
}

func TestRefreshToken(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	if refreshToken == "" {
		app.NewResponse(c).Error(errcode.ErrParams)
		return
	}
	userSvc := appservice.NewUserAppSvc(c)
	token, err := userSvc.TokenRefresh(refreshToken)
	if err != nil {
		if errors.Is(err, errcode.ErrTooManyRequests) {
			// 客户端有并发刷新token
			app.NewResponse(c).Error(errcode.ErrTooManyRequests)
			return
		} else {
			appErr := err.(*errcode.AppError)
			app.NewResponse(c).Error(appErr)
		}
		return
	}
	app.NewResponse(c).Success(token)
}
