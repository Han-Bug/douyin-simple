package controller

import (
	"douyin-simple/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

type VideoListResponse struct {
	models.Response
	VideoList []models.VideoRes `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	// 解析token

	// 将视频文件存入data
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 获取文件名称与用户名
	filename := filepath.Base(data.Filename)
	user := usersLoginInfo[token]
	// 打印相关信息
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	// 保存文件
	saveFile := filepath.Join("./public/video/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

func DoPublish(video models.VideoRes) {

}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {

	// TODO 解析用户token,
	token := c.PostForm("token")
	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "User_Login doesn't exist"})
		return
	}

	// TODO 从数据库获取指定用户的视频列表

	c.JSON(http.StatusOK, VideoListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
