package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ocean-marketing/internal/config"
	"ocean-marketing/internal/handler"
	"ocean-marketing/internal/middleware"
	"ocean-marketing/internal/pkg/database"
	"ocean-marketing/internal/pkg/logger"
	"ocean-marketing/internal/pkg/migration"
	"ocean-marketing/internal/pkg/redis"
	"ocean-marketing/internal/pkg/tracer"
	"ocean-marketing/internal/router"
	"ocean-marketing/pkg/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title Ocean Marketing API
// @version 1.0
// @description 这是一个基于Gin框架的完整项目骨架
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 初始化配置
	cfg := config.Init()

	// 初始化日志
	logger.Init(cfg.Log)

	// 初始化数据库
	database.Init(cfg.Database)

	// 数据库迁移
	migration.AutoMigrate()
	migration.CreateTables()
	migration.SeedData()

	// 初始化Redis
	redis.Init(cfg.Redis)

	// 初始化链路追踪
	tracer.Init(cfg.Tracer)

	// 初始化JWT
	jwt.Init(cfg.JWT)

	// 初始化handlers
	handler.Init(cfg)

	// 设置Gin模式
	gin.SetMode(cfg.App.Mode)

	// 创建Gin引擎
	r := gin.New()

	// 注册中间件
	middleware.Register(r, cfg)

	// 设置静态文件服务
	r.Static("/static", "./web/static")
	r.StaticFile("/", "./web/static/index.html")

	// 注册路由
	router.Register(r, cfg)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    cfg.App.Port,
		Handler: r,
	}

	// 启动服务器
	go func() {
		logger.Info("服务器启动", zap.String("addr", cfg.App.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("服务器关闭中...")

	// 5秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("服务器关闭失败", zap.Error(err))
	}

	logger.Info("服务器已关闭")
}
