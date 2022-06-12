package utils

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func ConnectDatabase(dbdsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dbdsn), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败：", err)
		PrintLog(err, "[Fatal]")
		//c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "an error occur when connecting database"})
		log.Fatalln(err)
		return nil, err
	}
	fmt.Println("数据库连接成功")
	return db, err
}
