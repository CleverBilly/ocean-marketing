package handler

import (
	"ocean-marketing/internal/config"
	"ocean-marketing/internal/middleware"

	"github.com/gin-gonic/gin"
)

var (
	exampleHandler *ExampleHandler
)

// Init 初始化handlers
func Init(cfg *config.Config) {
	exampleHandler = NewExampleHandler()
}

// 示例相关handlers
func GetExamples(c *gin.Context) {
	exampleHandler.GetExamples(c)
}

func GetExample(c *gin.Context) {
	exampleHandler.GetExample(c)
}

func CreateExample(c *gin.Context) {
	exampleHandler.CreateExample(c)
}

func UpdateExample(c *gin.Context) {
	exampleHandler.UpdateExample(c)
}

func DeleteExample(c *gin.Context) {
	exampleHandler.DeleteExample(c)
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return middleware.AuthMiddleware()
}

// 系统相关handlers
func Health(c *gin.Context) {
	HealthCheck(c)
}

func Ready(c *gin.Context) {
	ReadinessCheck(c)
}

func Live(c *gin.Context) {
	LivenessCheck(c)
}
