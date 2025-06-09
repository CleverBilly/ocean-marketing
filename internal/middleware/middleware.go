package middleware

import (
	"ocean-marketing/internal/config"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Register 注册中间件
func Register(r *gin.Engine, cfg *config.Config) {
	// CORS 跨域中间件
	r.Use(CORS())

	// 日志中间件
	r.Use(Logger())

	// Recovery 中间件
	r.Use(Recovery(cfg.Feishu))

	// 限流中间件
	r.Use(RateLimit())

	// 链路追踪中间件
	r.Use(Tracer())

	// Prometheus 指标中间件
	r.Use(Prometheus())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Prometheus 指标端点
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// pprof 性能分析
	if cfg.App.Mode == "debug" {
		pprof.Register(r)
	}
}
