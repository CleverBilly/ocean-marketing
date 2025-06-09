package middleware

import (
	"ocean-marketing/internal/pkg/tracer"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// Tracer 链路追踪中间件
func Tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中提取span上下文
		var span opentracing.Span
		wireContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header),
		)

		if err != nil {
			// 创建新的span
			span = tracer.StartSpan(c.Request.URL.Path)
		} else {
			// 基于已有上下文创建span
			span = tracer.StartSpanFromContext(wireContext, c.Request.URL.Path)
		}

		defer span.Finish()

		// 设置标签
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.String())
		ext.Component.Set(span, "gin-server")

		// 将span上下文注入到请求中
		c.Set("span", span)

		// 处理请求
		c.Next()

		// 设置状态码
		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))

		// 如果有错误，标记span
		if len(c.Errors) > 0 {
			ext.Error.Set(span, true)
			span.SetTag("error.message", c.Errors.String())
		}
	}
}

// GetSpan 从gin上下文中获取span
func GetSpan(c *gin.Context) opentracing.Span {
	if span, exists := c.Get("span"); exists {
		return span.(opentracing.Span)
	}
	return nil
}
