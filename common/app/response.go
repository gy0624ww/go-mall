package app

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-study-lab/go-mall/common/errcode"
	"github.com/go-study-lab/go-mall/common/logger"
)

type response struct {
	ctx        *gin.Context
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	RequestId  string      `json:"request_id"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *pagination `json:"pagination,omitempty"`
}

func NewResponse(c *gin.Context) *response {
	return &response{ctx: c}
}

// SetPagination 设置Response的分页信息
func (r *response) SetPagination(pagination *pagination) *response {
	r.Pagination = pagination
	return r
}

func (r *response) Success(data interface{}) {
	r.Code = errcode.Success.Code()
	r.Msg = errcode.Success.Msg()
	requestId := ""
	if _, exists := r.ctx.Get("traceid"); exists {
		val, _ := r.ctx.Get("traceid")
		requestId = val.(string)
	}
	r.RequestId = requestId
	r.Data = data
	r.ctx.JSON(errcode.Success.HttpStatusCode(), r)
}
func (r *response) SuccessOk() {
	r.Success("")
}

func (r *response) Error(err *errcode.AppError) {
	appErr := errcode.ErrServer.Clone() // 生成一个appErr 用作目标错误类型的判定
	if !errors.As(err, &appErr) {
		// 如果err不是appErr的类型, 把它变成appErr
		appErr = appErr.WithCause(err)
	}
	r.Code = appErr.Code()
	r.Msg = appErr.Msg()
	requestId := ""
	if _, exists := r.ctx.Get("traceid"); exists {
		val, _ := r.ctx.Get("traceid")
		requestId = val.(string)
	}
	r.RequestId = requestId
	// 兜底记一条响应错误，项目自定义的AppError中有错误链条,方便出错后排查问题
	logger.Error(r.ctx, "api_resonse_error", "err", err)
	r.ctx.JSON(appErr.HttpStatusCode(), r)
}
