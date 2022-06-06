package controller

import (
	"douyin-simple/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ConnectDatabase(dbdsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dbdsn), &gorm.Config{})
	return db, err
}

/*
	如果是点赞功能的话：
	1.判断是否登录，未登录直接返回错误结果（通过修改StatusCode值，可以参考其他模块的代码）
	2.判断是点赞还是取消点赞
	若是点赞：查询是否已经点赞过（在数据库中根据视频编号和用户编号查询指定点赞数据，
	数据存在则表示点过赞了，不存在则未点赞），若未点赞则进行点赞（数据库中插入数据），
	否则返回直接返回结果（StatusCode值为0，即代表正常执行）；取消点赞类似

*/
// FavoriteAction no practical effect, just check if token is valid
func ThumbUp(user models.User, db *gorm.DB, c *gin.Context) {
	//查询用户编号
	userId := user.Id
	//获取视频编号
	videoIdstr := c.Query("video_id")
	videoId, err := strconv.Atoi(videoIdstr)
	if err != nil {
		fmt.Println("videoId转换成int类型出错")
	}
	//通过用户编号和视频编号查询视频是否出现在点赞列表
	res := db.Model(&models.Favorite{}).Where("UserId = ? AND VideoId = ?", userId, videoId).Find(&user)
	if res == nil {
		favorite := models.Favorite{
			UserId:    userId,
			VideoId:   int64(videoId),
			CreatedAt: time.Now(),
		}
		db.Select("UserId", "VideoId", "CreatedAt").Create(&favorite)
	}
}

//取消点赞操作
func CancelThumbUp(db *gorm.DB, c *gin.Context) {
	//根据videoID将记录从数据库删除
	//获取videoId
	videoIdstr := c.Query("video_id")
	videoId, err := strconv.Atoi(videoIdstr)
	if err != nil {
		fmt.Println("videoId转换成int类型出错")
	}
	//执行删除操作
	favorite := models.Favorite{
		//UserId:    userId,
		VideoId:   int64(videoId),
		CreatedAt: time.Now(),
	}
	db.Where("video_id = ?", favorite.VideoId).Delete(&favorite)

}

func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, models.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "User_Login doesn't exist"})
	}
	//链接数据库
	dbdsn := "douyin:665577733_douYIN@tcp(119.23.68.131:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := ConnectDatabase(dbdsn)
	if err != nil {
		//fmt.Println("数据库失败")
		log.Fatalln(err)
	}
	fmt.Println("数据库连接成功")
	user, err := GetUserModelByToken(token)
	if err != nil {
		fmt.Printf("通过token获取user失败: %v ", err)
		log.Fatalln(err)
	}
	//获取操作请求，根据请求实现操作，点赞和取消点赞
	ActionType := c.Query("action_type")
	//获取actionType 1点赞 2取消点赞
	if ActionType == "1" {
		ThumbUp(user, db, c)
	} else if ActionType == "2" {
		CancelThumbUp(db, c)
	}

}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
