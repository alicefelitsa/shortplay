package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 本文件存放 web / boss 两个控制器共用的辅助函数，
// 使前端控制器（webController）与后台控制器（bossController）相互独立、互不依赖。

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

// serveSubtitle 字幕代理公共逻辑：查 drama_chapter.subtitle 取第 idx 条路径，
// 服务端 http.Get(domain + "/file" + path) 原样返回（text/plain，不做格式转换），
// 仅为规避浏览器跨域 fetch 限制。web / boss 两组各自暴露路由，后台不调用前端接口。
func serveSubtitle(db *gorm.DB, c *gin.Context) {
	parts := strings.Split(strings.TrimSuffix(c.Param("file"), ".srt"), "_")
	if len(parts) != 2 {
		c.String(http.StatusNotFound, "subtitle not found")
		return
	}
	idx := 0
	_, _ = fmt.Sscanf(parts[1], "%d", &idx)
	var subtitle string
	_ = db.Raw("select subtitle from drama_chapter where id = ?", parts[0]).Row().Scan(&subtitle)
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
	resp, err := http.Get(playDomain(db) + "/file" + p)
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

// pageLimit 处理Mysql数据分页：读取 page/limit 查询参数，拼成 limit 子句（同包共用，无需导出）。
func pageLimit(c *gin.Context) string {
	limit, err := strconv.Atoi(c.Query("limit"))
	page, err := strconv.Atoi(c.Query("page"))
	if page == 0 || limit == 0 || err != nil {
		return ""
	} else {
		page = (page - 1) * limit
		res := fmt.Sprintf(" limit %v,%v", page, limit)
		return res
	}
}
