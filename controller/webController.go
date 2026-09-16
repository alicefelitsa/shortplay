package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"shortplay/config"
	"time"
)

type WebController struct {
	db *gorm.DB
}

// NewWebController 创建前端控制器
func NewWebController() *WebController {
	return &WebController{
		db: config.Mysql,
	}
}

// GetDramaList 获取短剧列表（公开接口）
func (wc *WebController) GetDramaList(c *gin.Context) {
	var code, count int
	data := make([]map[string]interface{}, 0)
	_ = wc.db.Raw("select * from drama where status=1 order by id desc" + config.PageLimit(c)).Scan(&data).Error
	for _, row := range data {
		for col, val := range row {
			if t, ok := val.(time.Time); ok {
				row[col] = t.Format("2006-01-02 15:04:05")
			}
		}
	}
	_ = wc.db.Raw("select count(id) from drama where status=1").Scan(&count).Error
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": "操作成功",
		"count":   count,
		"data":    data,
	})
}

// GetDramaDetail 获取短剧详情
func (wc *WebController) GetDramaDetail(c *gin.Context) {
	var code int
	data := make([]map[string]interface{}, 0)
	episodes := make([]map[string]interface{}, 0)
	id := c.Query("id")
	_ = wc.db.Raw("select * from drama where id = ?", id).Scan(&data).Error
	if len(data) > 0 {
		for col, val := range data[0] {
			if t, ok := val.(time.Time); ok {
				data[0][col] = t.Format("2006-01-02 15:04:05")
			}
		}
		// 获取剧集列表
		_ = wc.db.Raw("select * from episode where drama_id = ? order by sort asc", id).Scan(&episodes).Error
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     code,
		"message":  "操作成功",
		"data":     data,
		"episodes": episodes,
	})
}

// GetCategory 获取分类列表
func (wc *WebController) GetCategory(c *gin.Context) {
	data := make([]map[string]interface{}, 0)
	_ = wc.db.Raw("select * from category where status=1 order by sort asc").Scan(&data).Error
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
		"data":    data,
	})
}

// Search 搜索短剧
func (wc *WebController) Search(c *gin.Context) {
	var count int
	keyword := c.Query("keyword")
	data := make([]map[string]interface{}, 0)
	query := "select * from drama where status=1 and title like ? order by id desc" + config.PageLimit(c)
	_ = wc.db.Raw(query, "%"+keyword+"%").Scan(&data).Error
	countQuery := "select count(id) from drama where status=1 and title like ?"
	_ = wc.db.Raw(countQuery, "%"+keyword+"%").Scan(&count).Error
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
		"count":   count,
		"data":    data,
	})
}
