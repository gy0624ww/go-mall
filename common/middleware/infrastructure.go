package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-study-lab/go-mall/common/util"
)

// infrastructure 中存放项目运行需要的基础中间件
func StartTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.Request.Header.Get("traceid")
		pSpanId := c.Request.Header.Get("spanid")
		spanId := util.GenerateSpanID(c.Request.RemoteAddr)
		if traceId == "" { // 如果traceId 为空，正面是链路的发端，把它设置成此次的spanId, 一般是初始节点，网关之类
			traceId = spanId
		}
		c.Set("traceid", traceId)
		c.Set("spanid", spanId)
		c.Set("pspanid", pSpanId)
		c.Next()
	}
}
