package controller

import (
	"douyin-simple/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type FeedResponse struct {
	models.Response
	VideoList []models.VideoRes `json:"video_list,omitempty"`
	NextTime  int64             `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	// 获取传入数据
	token := c.PostForm("token")
	latestTime := c.PostForm("latest_time")
	fmt.Println("latestTime: ", latestTime)
	fmt.Println("NowTime: ", time.Now())
	var user models.User
	if token != "" {
		// 解析token
		var err error
		user, err = GetUserModelByToken(token)
		if err != nil {
			// TODO 将错误打印至日志中
			fmt.Println(err)
			user = models.User{
				Id: -1,
			}
		}
	}
	// 连接数据库
	db, err := gorm.Open(
		mysql.Open(dbdsn),
	)

	if err != nil {
		// TODO 将错误打印至日志中
		fmt.Println(err)
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  "error occur when reading video data",
		})
		return
	}
	var videos []models.Video
	//db.Where("created_at > ?",).Order("created_at desc").Limit(20).Find(&videos)
	db.Order("created_at desc").Limit(20).Find(&videos)
	//if res.Error != nil {
	//	fmt.Println(res)
	//	fmt.Println(res.Error)
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: models.Response{StatusCode: 1, StatusMsg: "bad token"},
	//	})
	//	return
	//}
	var videoReses []models.VideoRes
	for i := 0; i < len(videos); i++ {
		// 查询视频的作者信息
		var author models.User
		db.Where("user_id = ?", videos[i].UserId).First(&author)
		authorRes, err2 := GetUserByUserModel(author, user.Id)
		if err2 != nil {
			// TODO 将错误打印至日志中
			fmt.Println(err)
			c.JSON(http.StatusOK, models.Response{
				StatusCode: 1,
				StatusMsg:  "error occur when getting author data",
			})
			return
		}
		// 查询视频的点赞数、评论数
		var favoriteCount int64
		db.Model(&models.Favorite{}).Where("video_id = ?", videos[i].Id).Count(&favoriteCount)
		var commentCount int64
		db.Model(&models.Comment{}).Where("video_id = ?", videos[i].Id).Count(&commentCount)
		// 查询该用户是否赞了该视频
		var isFavorite = false
		if user.Id > 0 {
			favorObj := models.Favorite{}
			db.Model(&models.Favorite{}).Where("video_id = ? And user_id = ?", videos[i].Id, user.Id).Find(&favorObj)
			if favorObj != (models.Favorite{}) {
				isFavorite = true
			}
		}
		videoReses = append(videoReses, models.VideoRes{
			Id:            videos[i].Id,
			Author:        authorRes,
			PlayUrl:       videos[i].PlayUrl,
			CoverUrl:      videos[i].CoverUrl,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
		})
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  models.Response{StatusCode: 0},
		VideoList: videoReses,
		NextTime:  time.Now().Unix(),
	})
}
