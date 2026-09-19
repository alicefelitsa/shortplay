package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shortplay/config"
	"strings"
)

// BossAuth 管理员授权验证
func BossAuth(c *gin.Context) {
	path := c.Request.URL.Path
	// 白名单路径直接放行
	if path == "/api/boss/login" || path == "/api/boss/logout" || path == "/api/boss/captcha" {
		c.Next()
		return
	}
	// 字幕代理供播放器（ArtPlayer）直接 fetch，无法携带登录 token，按前缀放行
	if strings.HasPrefix(path, "/api/boss/GetSubtitle/") {
		c.Next()
		return
	}
	Authorization := c.GetHeader("Authorization")
	if Authorization == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    401,
			"message": "授权不可为空",
		})
		c.Abort()
		return
	}
	uid, _ := config.Redis.Get(config.Ctx, Authorization).Result()
	if uid == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    401,
			"message": "登录异常请登录",
		})
		c.Abort()
		return
	}
	// 认证通过
	c.Next()
}
