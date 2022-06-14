package main

import (
	"douyin-simple/services"
	"douyin-simple/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	utils.InitLogger()
	services.InitDB()
	initRouter(r)

	r.Run(":8880") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
