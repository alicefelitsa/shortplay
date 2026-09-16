package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// IPRateLimiter IP请求限制器
type IPRateLimiter struct {
	mu        sync.RWMutex
	limiters  map[string]*limiterEntry
	rate      rate.Limit
	burst     int
	ttl       time.Duration // IP记录清理时间
	freezeDur time.Duration // 冻结时长
}

type limiterEntry struct {
	limiter     *rate.Limiter
	lastAccess  time.Time
	frozenUntil time.Time
}

// NewIPRateLimiter 创建IP请求限制器
// requestsPerMinute: 每分钟允许同一个IP的请求数
// burst: 突发请求数
// freezeDuration: 超限冻结时间
func NewIPRateLimiter(requestsPerMinute, burst int, freezeDuration time.Duration) gin.HandlerFunc {
	limiter := &IPRateLimiter{
		limiters:  make(map[string]*limiterEntry),
		rate:      rate.Limit(float64(requestsPerMinute) / 60.0),
		burst:     burst,
		ttl:       5 * time.Minute,
		freezeDur: freezeDuration,
	}

	go limiter.cleanup()
	return limiter.middleware()
}

func (l *IPRateLimiter) middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// 白名单（本地、内网）
		if isInternalIP(ip) {
			c.Next()
			return
		}

		entry := l.getOrCreateLimiter(ip)

		// 检查冻结
		if time.Now().Before(entry.frozenUntil) {
			retryAfter := entry.frozenUntil.Sub(time.Now())
			c.Header("Retry-After", strconv.FormatFloat(retryAfter.Seconds(), 'f', 0, 64))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
				"data": gin.H{
					"retry_after": retryAfter.Seconds(),
				},
			})
			return
		}

		// 令牌桶限流
		if !entry.limiter.Allow() {
			entry.frozenUntil = time.Now().Add(l.freezeDur)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": fmt.Sprintf("请求频率超限，IP已冻结%s", l.freezeDur),
			})
			return
		}

		// 更新最后访问时间
		entry.lastAccess = time.Now()

		// 添加限流头
		c.Header("X-RateLimit-Limit", strconv.Itoa(l.burst))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(int(entry.limiter.Tokens())))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))

		c.Next()
	}
}

func (l *IPRateLimiter) getOrCreateLimiter(ip string) *limiterEntry {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, exists := l.limiters[ip]
	if !exists || time.Since(entry.lastAccess) > l.ttl {
		entry = &limiterEntry{
			limiter:    rate.NewLimiter(l.rate, l.burst),
			lastAccess: time.Now(),
		}
		l.limiters[ip] = entry
	}

	entry.lastAccess = time.Now()
	return entry
}

func (l *IPRateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for ip, entry := range l.limiters {
			if now.Sub(entry.lastAccess) > l.ttl {
				delete(l.limiters, ip)
			}
		}
		l.mu.Unlock()
	}
}

func isInternalIP(ip string) bool {
	return ip == "127.0.0.1" || ip == "::1" ||
		ip == "localhost" || len(ip) == 0
}
