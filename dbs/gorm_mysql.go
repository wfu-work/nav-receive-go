package dbs

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"log"
	"nav-rtlogging-go/configs"
	"nav-rtlogging-go/global"
	"sync"
	"time"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func GormMysql() *gorm.DB {
	s := global.NAV_CONFIG.Mysql
	return initMysqlDatabase(s)
}

func GormMysqlByConfig(s configs.Mysql) *gorm.DB {
	return initMysqlDatabase(s)
}

func initMysqlDatabase(s configs.Mysql) *gorm.DB {
	if s.DBName == "" {
		return nil
	}
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
			s.Username, s.Password, s.Host, s.Port, s.DBName, "utf8mb4", true, "Local")
		var err error
		DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("❌数据库连接失败: %v", err)
		}
		sqlDB, err := DB.DB()
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetMaxIdleConns(25)
		sqlDB.SetConnMaxLifetime(time.Hour)
		if err := sqlDB.Ping(); err != nil {
			log.Fatalf("❌无法连接 MySQL: %v", err)
		}
		log.Printf("✅成功连接到 MySQL: %s - %d\n", s.Host, s.Port)
		global.NAV_DB = DB
	})
	return DB
}
