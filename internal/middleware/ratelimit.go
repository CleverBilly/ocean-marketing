package middleware

import (
	"net/http"
	"strconv"
	"time"

	"ocean-marketing/pkg/errno"
	"ocean-marketing/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimit 限流中间件
func RateLimit() gin.HandlerFunc {
	// 创建限流器，每分钟100个请求
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		// 使用客户端IP作为限流key
		key := c.ClientIP()

		context, err := instance.Get(c, key)
		if err != nil {
			response.InternalServerError(c, errno.InternalServerError)
			c.Abort()
			return
		}

		// 设置响应头
		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		if context.Reached {
			response.ErrorWithCode(c, http.StatusTooManyRequests, errno.ErrLimitExceed)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitWithConfig 自定义配置的限流中间件
func RateLimitWithConfig(period time.Duration, limit int64) gin.HandlerFunc {
	rate := limiter.Rate{
		Period: period,
		Limit:  limit,
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		key := c.ClientIP()

		context, err := instance.Get(c, key)
		if err != nil {
			response.InternalServerError(c, errno.InternalServerError)
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		if context.Reached {
			response.ErrorWithCode(c, http.StatusTooManyRequests, errno.ErrLimitExceed)
			c.Abort()
			return
		}

		c.Next()
	}
}
