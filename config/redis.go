package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var Redis *redis.Client
var Ctx = context.Background() // 定义一个全局的上下文

func init() {
	InitRedis()
}

// InitRedis 初始化Redis连接
func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     Conf.GetString("redis.address"),  // Redis 服务器地址
		Password: Conf.GetString("redis.password"), // 密码，没有则为空
		DB:       Conf.GetInt("redis.db"),          // 默认数据库

		// 连接池配置（生产环境推荐）
		PoolSize:        1000,             // 连接池最大连接数
		MinIdleConns:    20,               // 最小空闲连接数
		DialTimeout:     5 * time.Second,  // 连接建立超时时间
		ReadTimeout:     3 * time.Second,  // 读超时
		WriteTimeout:    3 * time.Second,  // 写超时
		PoolTimeout:     4 * time.Second,  // 获取连接超时时间
		ConnMaxIdleTime: 10 * time.Minute, // 连接最大空闲时间
	})

	// 测试连接是否成功
	if _, err := Redis.Ping(Ctx).Result(); err != nil {
		log.Fatal("无法连接到Redis：", err)
	}
	fmt.Println("Redis连接成功!")
}
