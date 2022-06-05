package controller

import (
	"douyin-simple/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	BASE_URL = "http://172.27.157.250:8880"
)

type VideoListResponse struct {
	models.Response
	VideoList []models.VideoRes `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	title := c.PostForm("token")
	// 如果token为空
	if token == "" {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1, StatusMsg: "need to login",
		})
		return
	}
	// 解析token
	user, err := GetUserModelByToken(token)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1, StatusMsg: "bad token",
		})
		return
	}
	// 将视频文件存入data
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1, StatusMsg: "video data error",
		})
		return
	}

	// 获取文件名称与用户名
	filename := filepath.Base(data.Filename)
	//user := usersLoginInfo[token]
	// 打印相关信息
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	// 保存文件
	filePath := filepath.Join("/public/video/", finalName)
	if err := c.SaveUploadedFile(data, "."+filePath); err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 在数据库中存储相关视频信息
	// 连接数据库
	db, err := gorm.Open(
		mysql.Open(dbdsn),
	)

	if err != nil {
		// TODO 将错误打印至日志中
		fmt.Println(err)
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  "error occur when saving video data",
		})
		// TODO 将已下载的视频文件删除
		return
	}
	videoUrl := strings.Join([]string{BASE_URL, filePath}, "")
	// TODO 从视频中截取封面
	// 此处使用的是临时封面
	video := models.Video{
		UserId:    strconv.FormatInt(user.Id, 10),
		Title:     title,
		PlayUrl:   videoUrl,
		CoverUrl:  "https://bkimg.cdn.bcebos.com/pic/80cb39dbb6fd5266d0168896e952802bd40735fa9855?x-bce-process=image/resize,m_lfit,w_536,limit_1/format,f_jpg",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}
	res := db.Create(&video)
	if res.Error != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  "error occur when saving video data",
		})
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList get a publishing list of one user
func PublishList(c *gin.Context) {
	token := c.PostForm("token")
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	// 如果目标用户id数据错误
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  "bad id"})
		return
	}
	// 判断是否登录
	var curUser models.User
	if token != "" {
		curUser, err = GetUserModelByToken(token)
		if err != nil {
			curUser = models.User{
				Id: -1,
			}
		}
	}
	// 获取指定用户数据库信息
	author, err2 := GetUserModelById(userId)
	if err2 != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  "author not exist"})
		return
	}
	// 获取指定用户详细信息 包含是否关注了该用户
	authorRes, err := GetUserByUserModel(author, curUser.Id)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  "error occur when getting author info"})
		return
	}
	// 获取用户的发布列表
	db, err := gorm.Open(
		mysql.Open(dbdsn),
	)
	if err != nil {
		// TODO 将错误打印至日志中
		fmt.Println(err)
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  "error occur when reading author data",
		})
		return
	}
	var videos []models.Video
	db.Where("user_id = ?", author.Id).Order("created_at desc").Find(videos)
	var videoReses []models.VideoRes
	for i := 0; i < len(videos); i++ {
		// 查询视频的点赞数、评论数
		var favoriteCount int64
		db.Model(&models.Favorite{}).Where("video_id = ?", videos[i].Id).Count(&favoriteCount)
		var commentCount int64
		db.Model(&models.Comment{}).Where("video_id = ?", videos[i].Id).Count(&commentCount)
		// 查询该用户是否赞了该视频
		var isFavorite = false
		if curUser.Id >= 0 {
			favorObj := models.Favorite{}
			db.Model(&models.Favorite{}).Where("video_id = ? And user_id = ?", videos[i].Id, curUser.Id).Find(&favorObj)
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
	c.JSON(http.StatusOK, VideoListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		VideoList: videoReses,
	})

}
