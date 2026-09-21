package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"shortplay/config"
	"shortplay/tools"
	"strings"
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

// formatTimeFields 将结果集中的 time.Time 字段格式化为字符串
func formatTimeFields(rows []map[string]interface{}) {
	for _, row := range rows {
		for col, val := range row {
			if t, ok := val.(time.Time); ok {
				row[col] = t.Format("2006-01-02 15:04:05")
			}
		}
	}
}

// AdminLogin 管理员登录
func (bc *BossController) AdminLogin(c *gin.Context) {
	data := make(map[string]interface{})
	_ = c.BindJSON(&data)
	// Pluck 进 []int：GORM 内部完成列类型→int 转换（避开 map 断言与 .Row() 的 nil panic）
	var ids []int
	bc.db.Table("admin").Where("account = ? and password = ?", data["account"], data["password"]).Pluck("id", &ids)
	if len(ids) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "账户或密码不正确"})
		return
	}
	token, err := tools.GenerateToken(ids[0], tools.RoleAdmin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "创建授权令牌失败：" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"token":   token,
		"account": data["account"],
	})
}

// AdminLogout 管理员退出（JWT 无状态，前端清除本地 token 即可）
func (bc *BossController) AdminLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "退出登录"})
}

// AuthUser 获取当前登录管理员信息
func (bc *BossController) AuthUser(c *gin.Context) {
	var code int
	var message string
	data := make(map[string]interface{})
	uid := c.GetInt("userID")
	resData := make([]map[string]interface{}, 0)
	err := bc.db.Raw("select * from admin where id = ?", uid).Scan(&resData).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
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
	c.JSON(http.StatusOK, gin.H{"code": code, "message": message, "data": data})
}

// ==================== 短剧管理 drama_book ====================

// GetDramaList 获取短剧列表（仅展示 flag=1）
func (bc *BossController) GetDramaList(c *gin.Context) {
	var code, count int
	conds := "flag = 1"
	bookName := c.Query("book_name")
	status := c.Query("status")
	typeId := c.Query("type_id")
	bookId := c.Query("book_id")
	if bookName != "" {
		conds += fmt.Sprintf(" and book_name like '%%%v%%'", bookName)
	}
	if status != "" {
		conds += fmt.Sprintf(" and status = '%v'", status)
	}
	if typeId != "" {
		conds += fmt.Sprintf(" and main_type_id = '%v'", typeId)
	}
	if bookId != "" {
		conds += fmt.Sprintf(" and book_id like '%%%v%%'", bookId)
	}
	where := " where " + conds
	data := make([]map[string]interface{}, 0)
	err := bc.db.Raw("select id,book_id,book_name,book_name_en,slug,cover,cover2,ratings,status,language,is_free,author,introduction,main_type_id,chapter_count,view_count,follow_count,created_at from drama_book" + where + " order by id desc" + pageLimit(c)).Scan(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	formatTimeFields(data)
	withCoverShow(data, playDomain(bc.db))
	err = bc.db.Raw("select count(id) from drama_book" + where).Scan(&count).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 501, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": code, "message": "操作成功", "count": count, "data": data})
}

// GetDramaDetail 获取短剧详情（按自增id）
func (bc *BossController) GetDramaDetail(c *gin.Context) {
	data := make([]map[string]interface{}, 0)
	err := bc.db.Raw("select * from drama_book where id = ?", c.Query("id")).Scan(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	formatTimeFields(data)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功", "data": data})
}

// AddDrama 添加短剧
func (bc *BossController) AddDrama(c *gin.Context) {
	data := make(map[string]interface{})
	_ = c.BindJSON(&data)
	delete(data, "cover_show")
	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()
	result := bc.db.Table("drama_book").Create(data)
	if result.Error == nil && result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功"})
	} else {
		msg := "操作失败"
		if result.Error != nil {
			msg = result.Error.Error()
		}
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": msg})
	}
}

// SaveDrama 修改短剧
func (bc *BossController) SaveDrama(c *gin.Context) {
	data := make(map[string]interface{})
	_ = c.BindJSON(&data)
	id := data["id"]
	delete(data, "id")
	delete(data, "created_at")
	delete(data, "cover_show")
	data["updated_at"] = time.Now()
	result := bc.db.Table("drama_book").Where("id = ?", id).Updates(data)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "操作失败"})
	}
}

// DelDrama 删除短剧（同时删除关联分集）
func (bc *BossController) DelDrama(c *gin.Context) {
	idList := tools.SplitIds(c.Query("ids"))
	// 先查出 book_id 以清理关联分集
	bookIds := make([]string, 0)
	_ = bc.db.Raw("select book_id from drama_book where id in (?)", idList).Scan(&bookIds).Error
	result := bc.db.Exec("delete from drama_book where id in (?)", idList)
	if len(bookIds) > 0 {
		bc.db.Exec("delete from drama_chapter where book_id in (?)", bookIds)
		bc.db.Exec("delete from drama_book_type where book_id in (?)", bookIds)
	}
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "操作失败"})
	}
}

// ==================== 分集管理 drama_chapter ====================

// GetChapterList 获取分集列表（按 book_id 过滤）
func (bc *BossController) GetChapterList(c *gin.Context) {
	var code, count int
	var where string
	bookId := c.Query("book_id")
	if bookId != "" {
		where = fmt.Sprintf(" where book_id = '%v'", bookId)
	} else {
		//默认列表只显示 flag=1 短剧关联的集（与 H5 展示同口径）；按剧集ID查询时不限
		where = " where exists (select 1 from drama_book b where b.book_id = drama_chapter.book_id and b.flag = 1)"
	}
	data := make([]map[string]interface{}, 0)
	err := bc.db.Raw("select id,book_id,chapter_id,chapter_name,chapter_index,chapter_index_str,is_unlock,chapter_price,duration,m3u8_flag,mp4_url,video_url,subtitle,created_at from drama_chapter" + where + " order by chapter_index asc" + pageLimit(c)).Scan(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	formatTimeFields(data)
	// 视频播放地址：/file + video_url 并做 HMAC 签名（CF Worker 校验），有效期读配置
	domain := playDomain(bc.db)
	secret := videoSecret(bc.db)
	expireMinutes := videoExpireMinutes(bc.db)
	for _, row := range data {
		vu, _ := row["video_url"].(string)
		if vu != "" {
			if !strings.HasPrefix(vu, "/") {
				vu = "/" + vu
			}
			row["play_url"] = tools.GenerateSignedVideoURL(domain, "/file"+vu, secret, expireMinutes)
		}
		// 字幕地址：JSON 数组，逐条签名供 <track> 加载
		subs := make([]string, 0)
		if subRaw, _ := row["subtitle"].(string); subRaw != "" && subRaw != "[]" {
			paths := make([]string, 0)
			if json.Unmarshal([]byte(subRaw), &paths) == nil {
				for _, p := range paths {
					if p == "" {
						continue
					}
					if !strings.HasPrefix(p, "/") {
						p = "/" + p
					}
					subs = append(subs, tools.GenerateSignedVideoURL(domain, "/file"+p, secret, expireMinutes))
				}
			}
		}
		row["subtitle_urls"] = subs
	}
	//关联剧名：按当页 book_id 批量查 drama_book 注入，供列表列展示（book_id 可能为数值类型，统一 Sprintf 归一化）
	bookIds := make([]string, 0)
	seen := map[string]bool{}
	for _, row := range data {
		if row["book_id"] == nil {
			continue
		}
		bid := fmt.Sprintf("%v", row["book_id"])
		if bid != "" && !seen[bid] {
			seen[bid] = true
			bookIds = append(bookIds, bid)
		}
	}
	if len(bookIds) > 0 {
		books := make([]map[string]interface{}, 0)
		if bc.db.Raw("select book_id,book_name from drama_book where book_id in (?)", bookIds).Scan(&books).Error == nil {
			nameMap := map[string]string{}
			for _, b := range books {
				nameMap[fmt.Sprintf("%v", b["book_id"])] = fmt.Sprintf("%v", b["book_name"])
			}
			for _, row := range data {
				if row["book_id"] != nil {
					row["book_name"] = nameMap[fmt.Sprintf("%v", row["book_id"])]
				}
			}
		}
	}
	err = bc.db.Raw("select count(id) from drama_chapter" + where).Scan(&count).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 501, "message": err.Error()})
		return
	}
	//剧名：按剧集ID查 drama_book 随列表返回，前端工具栏/播放标题展示（直接输入ID查询时无跳转参数）
	bookName := ""
	queryBookId := bookId
	if queryBookId == "" && len(data) > 0 && data[0]["book_id"] != nil {
		queryBookId = fmt.Sprintf("%v", data[0]["book_id"])
	}
	if queryBookId != "" {
		_ = bc.db.Raw("select book_name from drama_book where book_id = ?", queryBookId).Row().Scan(&bookName)
	}
	c.JSON(http.StatusOK, gin.H{"code": code, "message": "操作成功", "count": count, "data": data, "book_name": bookName})
}

// AddChapter 添加分集
func (bc *BossController) AddChapter(c *gin.Context) {
	data := make(map[string]interface{})
	_ = c.BindJSON(&data)
	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()
	result := bc.db.Table("drama_chapter").Create(data)
	if result.Error == nil && result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功"})
	} else {
		msg := "操作失败"
		if result.Error != nil {
			msg = result.Error.Error()
		}
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": msg})
	}
}

// SaveChapter 修改分集
func (bc *BossController) SaveChapter(c *gin.Context) {
	data := make(map[string]interface{})
	_ = c.BindJSON(&data)
	id := data["id"]
	delete(data, "id")
	delete(data, "created_at")
	data["updated_at"] = time.Now()
	result := bc.db.Table("drama_chapter").Where("id = ?", id).Updates(data)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "操作失败"})
	}
}

// DelChapter 删除分集
func (bc *BossController) DelChapter(c *gin.Context) {
	idList := tools.SplitIds(c.Query("ids"))
	result := bc.db.Exec("delete from drama_chapter where id in (?)", idList)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "操作失败"})
	}
}

// SetChapterUnlock 批量设置解锁状态（ids 逗号分隔；is_unlock 仅接受 0/1 校验后拼接）
func (bc *BossController) SetChapterUnlock(c *gin.Context) {
	ids := c.Query("ids")
	unlock := c.Query("is_unlock")
	if ids == "" || (unlock != "0" && unlock != "1") {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	unlockVal := 0
	if unlock == "1" {
		unlockVal = 1
	}
	result := bc.db.Exec("update drama_chapter set is_unlock = ?, updated_at = ? where id in (?)", unlockVal, time.Now(), tools.SplitIds(ids))
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "操作失败"})
	}
}

// ==================== 分类管理 drama_type ====================

// GetTypeList 获取分类列表
func (bc *BossController) GetTypeList(c *gin.Context) {
	var code, count int
	var where string
	typeName := c.Query("type_name")
	if typeName != "" {
		where = fmt.Sprintf(" where type_name like '%%%v%%'", typeName)
	}
	data := make([]map[string]interface{}, 0)
	err := bc.db.Raw("select * from drama_type" + where + " order by type_id asc" + pageLimit(c)).Scan(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	formatTimeFields(data)
	err = bc.db.Raw("select count(id) from drama_type" + where).Scan(&count).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 501, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": code, "message": "操作成功", "count": count, "data": data})
}

// AddType 添加分类
func (bc *BossController) AddType(c *gin.Context) {
	data := make(map[string]interface{})
	_ = c.BindJSON(&data)
	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()
	result := bc.db.Table("drama_type").Create(data)
	if result.Error == nil && result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "操作失败"})
	}
}

// SaveType 修改分类
func (bc *BossController) SaveType(c *gin.Context) {
	data := make(map[string]interface{})
	_ = c.BindJSON(&data)
	id := data["id"]
	delete(data, "id")
	delete(data, "created_at")
	data["updated_at"] = time.Now()
	result := bc.db.Table("drama_type").Where("id = ?", id).Updates(data)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "操作失败"})
	}
}

// DelType 删除分类
func (bc *BossController) DelType(c *gin.Context) {
	idList := tools.SplitIds(c.Query("ids"))
	result := bc.db.Exec("delete from drama_type where id in (?)", idList)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "操作失败"})
	}
}

// ==================== 站点配置 config ====================

// GetConfigSetting 获取站点配置
func (bc *BossController) GetConfigSetting(c *gin.Context) {
	data := make([]map[string]interface{}, 0)
	err := bc.db.Raw("select * from config order by id asc limit 1").Scan(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功", "data": data})
}

// SaveConfigSetting 保存站点配置
func (bc *BossController) SaveConfigSetting(c *gin.Context) {
	data := make(map[string]interface{})
	_ = c.BindJSON(&data)
	domain, _ := data["domain"].(string)
	accessMode, _ := data["access_mode"].(string)
	if accessMode != "pc" && accessMode != "h5" && accessMode != "all" {
		accessMode = "all"
	}
	displayType, _ := data["display_type"].(string)
	if displayType != "blank" && displayType != "jump" {
		displayType = "jump"
	}
	videoSecretKey, _ := data["video_secret_key"].(string)
	videoExpireMinutes := 1440
	if v, ok := data["video_expire_minutes"].(float64); ok && v > 0 {
		videoExpireMinutes = int(v)
	}
	var count int64
	_ = bc.db.Raw("select count(id) from config").Scan(&count).Error
	var result *gorm.DB
	if count == 0 {
		result = bc.db.Exec("insert into config (domain, access_mode, display_type, video_secret_key, video_expire_minutes) values (?, ?, ?, ?, ?)", domain, accessMode, displayType, videoSecretKey, videoExpireMinutes)
	} else {
		result = bc.db.Exec("update config set domain=?, access_mode=?, display_type=?, video_secret_key=?, video_expire_minutes=? order by id asc limit 1", domain, accessMode, displayType, videoSecretKey, videoExpireMinutes)
	}
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功"})
}

// ==================== 文件上传 ====================

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

// ==================== 字幕代理 ====================

// GetSubtitle 字幕代理（后台专用）：走 boss 组，与前端 h5 接口分离，后台不调用 /api/web。
// 复用 serveSubtitle 公共逻辑；播放器 fetch 字幕无法携带登录 token，路由在 BossAuth 白名单按前缀放行。
func (bc *BossController) GetSubtitle(c *gin.Context) {
	serveSubtitle(bc.db, c)
}
