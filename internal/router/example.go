package router

import (
	"ocean-marketing/internal/handler"
	"ocean-marketing/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterExampleRoutes 注册示例模块的路由
func RegisterExampleRoutes(v1 *gin.RouterGroup) {
	exampleHandler := handler.NewExampleHandler()

	// 示例业务路由
	examples := v1.Group("/examples")
	{
		examples.GET("", exampleHandler.GetExamples)                                       // 公开访问
		examples.GET("/:id", exampleHandler.GetExample)                                    // 公开访问
		examples.POST("", middleware.AuthMiddleware(), exampleHandler.CreateExample)       // 需要认证
		examples.PUT("/:id", middleware.AuthMiddleware(), exampleHandler.UpdateExample)    // 需要认证
		examples.DELETE("/:id", middleware.AuthMiddleware(), exampleHandler.DeleteExample) // 需要认证
	}
}
