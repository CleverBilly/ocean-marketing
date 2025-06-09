package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"ocean-marketing/internal/config"
	"ocean-marketing/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// FeishuMessage 飞书消息结构
type FeishuMessage struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

// Recovery 恢复中间件
func Recovery(feishuCfg config.FeishuConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := make([]byte, 4096)
				length := runtime.Stack(stack, false)
				stackTrace := string(stack[:length])

				// 记录错误日志
				logger.Error("发生panic",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("ip", c.ClientIP()),
					zap.String("stack", stackTrace),
				)

				// 发送飞书通知
				if feishuCfg.WebhookURL != "" {
					go sendFeishuNotification(feishuCfg.WebhookURL, err, c, stackTrace)
				}

				// 返回500错误
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "内部服务器错误",
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}

// sendFeishuNotification 发送飞书通知
func sendFeishuNotification(webhookURL string, err interface{}, c *gin.Context, stackTrace string) {
	message := FeishuMessage{
		MsgType: "text",
	}

	text := fmt.Sprintf("🚨 服务异常告警\n"+
		"时间: %s\n"+
		"路径: %s %s\n"+
		"IP: %s\n"+
		"错误: %v\n"+
		"堆栈: %s",
		time.Now().Format("2006-01-02 15:04:05"),
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
		err,
		stackTrace,
	)

	message.Content.Text = text

	jsonData, jsonErr := json.Marshal(message)
	if jsonErr != nil {
		logger.Error("构造飞书消息失败", zap.Error(jsonErr))
		return
	}

	resp, httpErr := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if httpErr != nil {
		logger.Error("发送飞书通知失败", zap.Error(httpErr))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("飞书通知响应异常", zap.Int("status", resp.StatusCode))
	}
}
