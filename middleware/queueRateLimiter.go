package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"strconv"
	"time"
)

// QueueRateLimiter 并发请求控制器中间件（非阻塞版本）
// rps：每秒并发请求数
// burst：突发并发请求数
func QueueRateLimiter(rps int, burst int) gin.HandlerFunc {
	if rps <= 0 || burst <= 0 {
		log.Fatal("并发请求控制器参数错误：rps 和 burst 必须大于 0")
	}

	fmt.Printf("并发请求控制器启动：%d 请求/秒，突发 %d 请求\n", rps, burst)

	limiter := rate.NewLimiter(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		// 非阻塞检查：立即返回是否允许
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"error":   "请求过多，请稍后重试",
				"message": fmt.Sprintf("当前服务繁忙，请稍后再试（限制：%d 请求/秒）", rps),
			})
			return
		}

		// 设置限流头部信息
		c.Header("X-RateLimit-Limit", strconv.Itoa(rps))
		c.Header("X-RateLimit-Burst", strconv.Itoa(burst))
		c.Header("X-RateLimit-Remaining", strconv.FormatFloat(limiter.Tokens(), 'f', 0, 64))

		// 计算重置时间
		resetSeconds := float64(burst) / float64(rps)
		resetTime := time.Now().Add(time.Duration(resetSeconds * float64(time.Second)))
		c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))

		c.Next()
	}
}
