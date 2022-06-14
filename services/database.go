package services

import (
	"douyin-simple/utils"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	DSN       = "douyin:665577733_douYIN@tcp(119.23.68.131:3306)/douyin?parseTime=true"
	GLOBAL_DB *gorm.DB
)

func ConnectDatabase() (*gorm.DB, error) {
	if GLOBAL_DB == nil {
		// 数据库未初始化
		err := errors.New("数据库未初始化")
		utils.PrintLogError(err)
		return nil, err
	}
	return GLOBAL_DB, nil
}
func InitDB() {
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		utils.PrintLogFatal("数据库连接失败:", err)
		return
	}
	sqlDB, _ := db.DB()
	if err != nil {
		return
	}
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	GLOBAL_DB = db
}
