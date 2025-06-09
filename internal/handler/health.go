package handler

import (
	"net/http"
	"time"

	"ocean-marketing/internal/pkg/database"
	"ocean-marketing/internal/pkg/redis"

	"github.com/gin-gonic/gin"
)

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}

// HealthCheck 健康检查
// @Summary 健康检查
// @Description 检查服务及其依赖的健康状态
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse "服务正常"
// @Failure 503 {object} HealthResponse "服务异常"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	services := make(map[string]string)
	overallStatus := "healthy"

	// 检查数据库连接
	if db := database.GetDB(); db != nil {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			services["database"] = "unhealthy"
			overallStatus = "unhealthy"
		} else {
			services["database"] = "healthy"
		}
	} else {
		services["database"] = "unhealthy"
		overallStatus = "unhealthy"
	}

	// 检查Redis连接
	if redisClient := redis.GetClient(); redisClient != nil {
		if _, err := redisClient.Ping(c).Result(); err != nil {
			services["redis"] = "unhealthy"
			overallStatus = "unhealthy"
		} else {
			services["redis"] = "healthy"
		}
	} else {
		services["redis"] = "unhealthy"
		overallStatus = "unhealthy"
	}

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Services:  services,
		Version:   "1.0.0",
	}

	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// ReadinessCheck 就绪检查
// @Summary 就绪检查
// @Description 检查服务是否准备好接收请求
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "服务就绪"
// @Failure 503 {object} map[string]interface{} "服务未就绪"
// @Router /ready [get]
func ReadinessCheck(c *gin.Context) {
	// 简单的就绪检查，可以根据需要扩展
	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now(),
	})
}

// LivenessCheck 存活检查
// @Summary 存活检查
// @Description 检查服务是否存活
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "服务存活"
// @Router /live [get]
func LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now(),
	})
}
