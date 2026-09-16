package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"shortplay/config"
	"shortplay/tools"
	"time"
)

type BossController struct {
	db *gorm.DB
}

// NewBossController 创建管理员控制器
func NewBossController() *BossController {
	return &BossController{
		db: config.Mysql,
	}
}

// AdminLogin 管理员登录
func (bc *BossController) AdminLogin(c *gin.Context) {
	var code int
	data := make(map[string]interface{})
	_ = c.BindJSON(&data)
	resData := make([]map[string]interface{}, 0)
	err := bc.db.Raw("select id from admin where account = ? and password = ?", data["account"], data["password"]).Scan(&resData).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	if len(resData) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "账户或密码不正确"})
		return
	}
	token := tools.CreateARandomString(30)
	err = config.Redis.Set(config.Ctx, token, resData[0]["id"], time.Minute*43200).Err()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "创建授权令牌失败：" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": "登录成功",
		"token":   token,
		"account": data["account"],
	})
}

// AdminLogout 管理员退出
func (bc *BossController) AdminLogout(c *gin.Context) {
	Authorization := c.GetHeader("Authorization")
	if Authorization != "" {
		config.Redis.Del(config.Ctx, Authorization)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "退出登录",
	})
}

// AuthUser 获取当前登录管理员信息
func (bc *BossController) AuthUser(c *gin.Context) {
	var code int
	var message string
	data := make(map[string]interface{})
	uid, _ := config.Redis.Get(config.Ctx, c.GetHeader("Authorization")).Result()
	resData := make([]map[string]interface{}, 0)
	err := bc.db.Raw("select * from admin where id = ?", uid).Scan(&resData).Error
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	if len(resData) > 0 {
		code = 0
		message = "操作成功"
		data["userId"] = resData[0]["id"]
		data["account"] = resData[0]["account"]
		data["nickname"] = resData[0]["account"]
	} else {
		code = 400
		message = "管理员不存在"
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

// GetDramaList 获取短剧列表
func (bc *BossController) GetDramaList(c *gin.Context) {
	var code, count int
	var where string
	title := c.Query("title")
	if title != "" {
		where += fmt.Sprintf("title like '%%%v%%' and ", title)
	}
	if where != "" {
		where = " where " + where[:len(where)-5]
	}
	data := make([]map[string]interface{}, 0)
	err := bc.db.Raw("select * from drama" + where + " order by id desc" + config.PageLimit(c)).Scan(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	for _, row := range data {
		for col, val := range row {
			if t, ok := val.(time.Time); ok {
				row[col] = t.Format("2006-01-02 15:04:05")
			}
		}
	}
	err = bc.db.Raw("select count(id) from drama" + where).Scan(&count).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 501, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": "操作成功",
		"count":   count,
		"data":    data,
	})
}

// AddDrama 添加短剧
func (bc *BossController) AddDrama(c *gin.Context) {
	data := make(map[string]interface{})
	_ = c.BindJSON(&data)
	data["ctime"] = time.Now()
	result := bc.db.Table("drama").Create(data)
	if result.Error == nil && result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "操作失败"})
	}
}

// SaveDrama 修改短剧
func (bc *BossController) SaveDrama(c *gin.Context) {
	data := make(map[string]interface{})
	_ = c.BindJSON(&data)
	result := bc.db.Table("drama").Where("id = ?", data["id"]).Updates(data)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "操作失败"})
	}
}

// DelDrama 删除短剧
func (bc *BossController) DelDrama(c *gin.Context) {
	ids := c.Query("ids")
	result := bc.db.Exec("delete from drama where id in(" + ids + ")")
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "操作失败"})
	}
}

// UploadImage 上传图片
func (bc *BossController) UploadImage(c *gin.Context) {
	filePath := "/uploads"
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "获取上传文件失败"})
		return
	}
	fileName, fileUrl, err := tools.SaveImageFile(c, "."+filePath, file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "上传失败：" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     0,
		"message":  "操作成功",
		"url":      filePath + fileUrl,
		"fileName": fileName,
	})
}
