package config

import (
	"fmt"
	"github.com/tagphi/czdb-search-golang/pkg/db"
	"log"
)

var Cz88Ip *db.DBSearcher

func init() {
	InitCz88Ip()
}

// InitCz88Ip 初始化IP纯真社区版
// IP 库为启动强依赖：加载失败直接 log.Fatal 终止启动（业务要求：缺库不允许运行）。
// 路径为相对路径，需从项目根目录运行（与 ./uploads、config.yaml 一致）。
func InitCz88Ip() {
	var err error
	Cz88Ip, err = db.InitDBSearcher("./cz88_public_v4.czdb", "s4s5fO8FegK89uxtvM8seg==", db.MEMORY)
	if err != nil {
		log.Fatal("初始化IP纯真社区版失败: ", err)
	}
	fmt.Println("初始化IP纯真社区版成功!")
}
