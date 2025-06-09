package router

import (
	"ocean-marketing/internal/config"
	"ocean-marketing/internal/handler"

	"github.com/gin-gonic/gin"
)

// Register 注册路由
func Register(r *gin.Engine, cfg *config.Config) {
	// 健康检查路由（不需要认证）
	r.GET("/health", handler.Health)
	r.GET("/ready", handler.Ready)
	r.GET("/live", handler.Live)

	// API 路由组
	api := r.Group("/api")
	{
		// API v1 路由组
		v1Group := api.Group("/v1")

		// 注册各模块路由
		RegisterExampleRoutes(v1Group)
	}
}
