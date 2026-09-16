package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shortplay/controller"
	"shortplay/middleware"
)

// SetupRouter 设置路由地址
func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.Cors())
	//router.Use(middleware.NewIPRateLimiter(120, 20, 1*time.Minute))
	//router.Use(middleware.QueueRateLimiter(config.Conf.GetInt("server.queueCapacity"), config.Conf.GetInt("server.suddenCapacity")))
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "页面不存在",
		})
	})
	// 管理后台接口
	boss := router.Group("/api/boss", middleware.BossAuth)
	{
		bossController := controller.NewBossController()
		boss.POST("/login", bossController.AdminLogin)
		boss.GET("/logout", bossController.AdminLogout)
		boss.GET("/auth/user", bossController.AuthUser)
		boss.GET("/GetDramaList", bossController.GetDramaList)
		boss.POST("/AddDrama", bossController.AddDrama)
		boss.POST("/SaveDrama", bossController.SaveDrama)
		boss.GET("/DelDrama", bossController.DelDrama)
		boss.POST("/UploadImage", bossController.UploadImage)
	}
	// 前端公开接口
	web := router.Group("/api/web")
	{
		webController := controller.NewWebController()
		web.GET("/GetDramaList", webController.GetDramaList)
		web.GET("/GetDramaDetail", webController.GetDramaDetail)
		web.GET("/GetCategory", webController.GetCategory)
		web.GET("/Search", webController.Search)
	}
	// 上传的图片静态服务
	router.Static("/uploads", "./uploads")
	return router
}
