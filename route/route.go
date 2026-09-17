package route

import (
	"net/http"
	"shortplay/controller"
	"shortplay/middleware"

	"github.com/gin-gonic/gin"
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
		// 登录鉴权
		boss.POST("/login", bossController.AdminLogin)
		boss.GET("/logout", bossController.AdminLogout)
		boss.GET("/auth/user", bossController.AuthUser)
		// 短剧管理
		boss.GET("/GetDramaList", bossController.GetDramaList)
		boss.GET("/GetDramaDetail", bossController.GetDramaDetail)
		boss.POST("/AddDrama", bossController.AddDrama)
		boss.POST("/SaveDrama", bossController.SaveDrama)
		boss.GET("/DelDrama", bossController.DelDrama)
		// 分集管理
		boss.GET("/GetChapterList", bossController.GetChapterList)
		boss.POST("/AddChapter", bossController.AddChapter)
		boss.POST("/SaveChapter", bossController.SaveChapter)
		boss.GET("/DelChapter", bossController.DelChapter)
		boss.GET("/SetChapterUnlock", bossController.SetChapterUnlock)
		// 分类管理
		boss.GET("/GetTypeList", bossController.GetTypeList)
		boss.POST("/AddType", bossController.AddType)
		boss.POST("/SaveType", bossController.SaveType)
		boss.GET("/DelType", bossController.DelType)
		// 站点配置
		boss.GET("/GetConfigSetting", bossController.GetConfigSetting)
		boss.POST("/SaveConfigSetting", bossController.SaveConfigSetting)
		// 文件上传
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
		web.GET("/GetSiteConfig", webController.GetSiteConfig)
		// 字幕代理（SRT 原样返回，供播放器同源加载）
		web.GET("/GetSubtitle/:file", webController.GetSubtitle)
	}
	// 上传的图片静态服务
	router.Static("/uploads", "./uploads")
	return router
}
