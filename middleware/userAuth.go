package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shortplay/config"
)

// UserAuth 用户授权验证
func UserAuth(c *gin.Context) {
	path := c.Request.URL.Path
	// 白名单路径直接放行
	if path == "/api/user/login" || path == "/api/user/register" {
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
			"message": "账户异常请登录",
		})
		c.Abort()
		return
	}

	// 认证通过
	c.Next()
}
