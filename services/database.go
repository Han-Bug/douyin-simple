package services

import (
	"douyin-simple/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DSN = "douyin:665577733_douYIN@tcp(119.23.68.131:3306)/douyin?parseTime=true"
)

func ConnectDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		//fmt.Println("数据库连接失败：", err)
		utils.PrintLog(err, "[Fatal]")
		//log.Println(err)
	}
	//fmt.Println("数据库连接成功")
	return db, err
}
