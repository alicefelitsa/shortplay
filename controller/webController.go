package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"shortplay/config"
	"shortplay/tools"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// playDomain 读取播放域名（config.domain），去掉尾部斜杠
func playDomain(db *gorm.DB) string {
	var domain string
	_ = db.Raw("select domain from config order by id asc limit 1").Scan(&domain).Error
	return strings.TrimRight(domain, "/")
}

// videoSecret 读取视频签名密钥（config.video_secret_key，与 CF Worker SECRET_KEY 一致）
func videoSecret(db *gorm.DB) string {
	var key string
	_ = db.Raw("select video_secret_key from config order by id asc limit 1").Scan(&key).Error
	return key
}

// videoExpireMinutes 读取视频签名有效期（config.video_expire_minutes，分钟），非法时回退 1440
func videoExpireMinutes(db *gorm.DB) int {
	var minutes int
	_ = db.Raw("select video_expire_minutes from config order by id asc limit 1").Scan(&minutes).Error
	if minutes <= 0 {
		minutes = 1440
	}
	return minutes
}

// withCoverShow 为结果集拼接展示封面：播放域名 + /file + cover2
// （对接 CF Worker 回源 B2，图片无签名可直接访问，路径需带 /file 前缀）
func withCoverShow(rows []map[string]interface{}, domain string) {
	for _, row := range rows {
		cover2, _ := row["cover2"].(string)
		if cover2 != "" {
			if !strings.HasPrefix(cover2, "/") {
				cover2 = "/" + cover2
			}
			row["cover_show"] = domain + "/file" + cover2
		} else {
			cover, _ := row["cover"].(string)
			row["cover_show"] = cover
		}
	}
}

// GetDramaList 获取短剧列表（仅已发布，支持按分类过滤）
func (wc *WebController) GetDramaList(c *gin.Context) {
	var code, count int
	var where = " where status = 'PUBLISHED' and flag = 1"
	typeId := c.Query("type_id")
	if typeId != "" {
		// 通过关联表筛选该分类下的短剧
		where += " and book_id in (select book_id from drama_book_type where type_id = " + typeId + ")"
	}
	data := make([]map[string]interface{}, 0)
	_ = wc.db.Raw("select book_id,book_name,book_name_en,slug,cover,cover2,ratings,introduction,chapter_count,view_count,follow_count,is_free,author,main_type_id,last_update_text,view_count_text from drama_book" + where + " order by shelf_time desc, id desc" + config.PageLimit(c)).Scan(&data).Error
	_ = wc.db.Raw("select count(id) from drama_book" + where).Scan(&count).Error
	withCoverShow(data, playDomain(wc.db))
	c.JSON(http.StatusOK, gin.H{"code": code, "message": "操作成功", "count": count, "data": data})
}

// GetDramaDetail 获取短剧详情 + 分集列表（按 book_id 或 slug）
func (wc *WebController) GetDramaDetail(c *gin.Context) {
	var code int
	bookId := c.Query("book_id")
	slug := c.Query("slug")
	data := make([]map[string]interface{}, 0)
	episodes := make([]map[string]interface{}, 0)
	if bookId != "" {
		_ = wc.db.Raw("select * from drama_book where book_id = ? and flag = 1", bookId).Scan(&data).Error
	} else if slug != "" {
		_ = wc.db.Raw("select * from drama_book where slug = ? and flag = 1", slug).Scan(&data).Error
	}
	if len(data) > 0 {
		for col, val := range data[0] {
			if t, ok := val.(time.Time); ok {
				data[0][col] = t.Format("2006-01-02 15:04:05")
			}
		}
		domain := playDomain(wc.db)
		withCoverShow(data, domain)
		bid, _ := data[0]["book_id"].(string)
		// 分集列表：只返回播放需要的字段
		_ = wc.db.Raw("select chapter_id,chapter_name,chapter_index,chapter_index_str,is_unlock,chapter_price,duration,cover,mp4_url,video_url,m3u8_flag from drama_chapter where book_id = ? order by chapter_index asc", bid).Scan(&episodes).Error
		// 视频播放地址：/file + video_url 并做 HMAC 签名（CF Worker 校验），有效期读配置
		secret := videoSecret(wc.db)
		expireMinutes := videoExpireMinutes(wc.db)
		for _, ep := range episodes {
			vu, _ := ep["video_url"].(string)
			if vu != "" {
				if !strings.HasPrefix(vu, "/") {
					vu = "/" + vu
				}
				ep["play_url"] = tools.GenerateSignedVideoURL(domain, "/file"+vu, secret, expireMinutes)
			}
		}
		// 推荐位
		recommends := make([]map[string]interface{}, 0)
		_ = wc.db.Raw("select b.book_id,b.book_name,b.slug,b.cover,b.cover2,b.ratings,b.chapter_count from drama_recommend r join drama_book b on b.book_id = r.recommend_book_id where r.book_id = ? and b.flag = 1 order by r.sort_no asc limit 12", bid).Scan(&recommends).Error
		withCoverShow(recommends, domain)
		c.JSON(http.StatusOK, gin.H{"code": code, "message": "操作成功", "data": data, "episodes": episodes, "recommends": recommends})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 404, "message": "短剧不存在", "data": data, "episodes": episodes})
}

// GetCategory 获取分类列表
func (wc *WebController) GetCategory(c *gin.Context) {
	data := make([]map[string]interface{}, 0)
	_ = wc.db.Raw("select type_id,type_name,replace_name from drama_type order by type_id asc").Scan(&data).Error
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功", "data": data})
}

// Search 搜索短剧（按剧名，中英文）
func (wc *WebController) Search(c *gin.Context) {
	var count int
	keyword := c.Query("keyword")
	data := make([]map[string]interface{}, 0)
	like := "%" + keyword + "%"
	where := " where status = 'PUBLISHED' and flag = 1 and (book_name like ? or book_name_en like ? or book_name_lower like ?)"
	_ = wc.db.Raw("select book_id,book_name,book_name_en,slug,cover,cover2,ratings,chapter_count,view_count,is_free,author from drama_book"+where+" order by view_count desc"+config.PageLimit(c), like, like, strings.ToLower(like)).Scan(&data).Error
	_ = wc.db.Raw("select count(id) from drama_book"+where, like, like, strings.ToLower(like)).Scan(&count).Error
	withCoverShow(data, playDomain(wc.db))
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功", "count": count, "data": data})
}

// GetSiteConfig 获取站点公开配置（访问方式、展示类型、域名）
func (wc *WebController) GetSiteConfig(c *gin.Context) {
	var cfg struct {
		Domain      string
		AccessMode  string
		DisplayType string
	}
	_ = wc.db.Raw("select domain, access_mode, display_type from config order by id asc limit 1").Scan(&cfg).Error
	if cfg.AccessMode == "" {
		cfg.AccessMode = "all"
	}
	if cfg.DisplayType == "" {
		cfg.DisplayType = "jump"
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
		"data": gin.H{
			"domain":       cfg.Domain,
			"access_mode":  cfg.AccessMode,
			"display_type": cfg.DisplayType,
		},
	})
}

// GetSubtitle 字幕代理：服务端拉取字幕原样返回（播放器自行解析 SRT），规避浏览器跨域 fetch 限制
// 路由格式 /GetSubtitle/{id}_{i}.srt，.srt 后缀便于播放器按扩展名识别字幕类型
func (wc *WebController) GetSubtitle(c *gin.Context) {
	parts := strings.Split(strings.TrimSuffix(c.Param("file"), ".srt"), "_")
	if len(parts) != 2 {
		c.String(http.StatusNotFound, "subtitle not found")
		return
	}
	idx := 0
	_, _ = fmt.Sscanf(parts[1], "%d", &idx)
	var subtitle string
	_ = wc.db.Raw("select subtitle from drama_chapter where id = ?", parts[0]).Row().Scan(&subtitle)
	paths := make([]string, 0)
	if subtitle != "" && subtitle != "[]" {
		_ = json.Unmarshal([]byte(subtitle), &paths)
	}
	if idx < 0 || idx >= len(paths) {
		c.String(http.StatusNotFound, "subtitle not found")
		return
	}
	p := paths[idx]
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	// 字幕无需签名（CF Worker 仅校验视频签名），直接拼路径拉取
	resp, err := http.Get(playDomain(wc.db) + "/file" + p)
	if err != nil {
		c.String(http.StatusBadGateway, "fetch subtitle failed: "+err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.String(http.StatusBadGateway, "fetch subtitle failed: status "+resp.Status)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "text/plain; charset=utf-8", body)
}
