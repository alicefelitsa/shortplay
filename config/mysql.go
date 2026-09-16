package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
	"time"
)

var Mysql *gorm.DB

func init() {
	InitMysql()
}

// InitMysql 初始化Mysql连接
func InitMysql() {
	dbAddress := Conf.GetString("mysql.address")
	dbName := Conf.GetString("mysql.database")
	dbUser := Conf.GetString("mysql.user")
	dbPasswd := Conf.GetString("mysql.password")
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=true&loc=Asia%%2FShanghai&timeout=5s&readTimeout=30s&writeTimeout=30s&interpolateParams=true", dbUser, dbPasswd, dbAddress, dbName)
	var err error
	Mysql, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("初始化Mysql时出错：", err)
	}
	// 设置连接池参数
	sqlDB, err := Mysql.DB()
	if err != nil {
		log.Fatal("获取Mysql连接池时出错：", err)
	}
	sqlDB.SetMaxOpenConns(100)                 // 最大连接数
	sqlDB.SetMaxIdleConns(10)                  // 空闲连接数
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // 连接最大存活时间
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 空闲连接最大存活时间
	if err = sqlDB.Ping(); err != nil {
		log.Fatal("连接到Mysql时出错：", err)
	}
	fmt.Println("Mysql连接成功！")
}

// PageLimit 处理Mysql数据分页
func PageLimit(c *gin.Context) string {
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

// PrintMysqlStats 连接池监控打印
func PrintMysqlStats() {
	sqlDB, err := Mysql.DB()
	if err != nil {
		fmt.Println("获取Mysql连接池失败：", err)
		return
	}
	stats := sqlDB.Stats()
	fmt.Printf("连接池状态:\n")
	fmt.Printf("最大打开连接数: %d\n", stats.MaxOpenConnections)
	fmt.Printf("打开连接数: %d\n", stats.OpenConnections)
	fmt.Printf("使用中连接数: %d\n", stats.InUse)
	fmt.Printf("空闲连接数: %d\n", stats.Idle)
	fmt.Printf("等待连接数: %d\n", stats.WaitCount)
	fmt.Printf("等待时间总计: %v\n", stats.WaitDuration)
	fmt.Printf("最大空闲时间关闭数: %d\n", stats.MaxIdleTimeClosed)
	fmt.Printf("最大生命周期关闭数: %d\n", stats.MaxLifetimeClosed)
}

// warmupConnections 预热指定数量的连接
func warmupConnections(db *gorm.DB, count int) {
	start := time.Now()
	var wg sync.WaitGroup
	errChan := make(chan error, count)
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := db.Exec("SELECT 1").Error; err != nil {
				errChan <- err
			}
		}()
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			log.Fatalf("连接池预热失败：%v", err)
		}
	}
	fmt.Printf("Mysql已预热%v个连接，耗时：%v\n", count, time.Since(start))
}

// keepAlive 定时保活连接池
func keepAlive(db *gorm.DB) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for t := range ticker.C {
		err := db.Exec("SELECT 1").Error
		sprintf := fmt.Sprintf("Mysql保活连接池：%s %v", t.Format("2006-01-02 15:04:05"), err)
		fmt.Println(sprintf)
		LogInfo("%s", sprintf)
	}
}
