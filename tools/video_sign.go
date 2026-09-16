package tools

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

// GenerateSignedVideoURL 生成带 HMAC-SHA256 签名的视频播放地址，
// 与 CF Worker 的 SECRET_KEY 校验逻辑一致（message = path:expires）。
// domain 为播放域名（不含尾斜杠）；videoPath 为含 /file 前缀的路径，
// 如 /file/video/42000024682/701478952_episode_1/video.mp4；
// expireMinutes 为签名有效期（分钟）。
func GenerateSignedVideoURL(domain, videoPath, secretKey string, expireMinutes int) string {
	expires := time.Now().Add(time.Duration(expireMinutes) * time.Minute).Unix()
	message := fmt.Sprintf("%s:%d", videoPath, expires)
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(message))
	sig := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("%s%s?expires=%d&signature=%s", domain, videoPath, expires, sig)
}
