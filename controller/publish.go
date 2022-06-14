package controller

import (
	"douyin-simple/models"
	"douyin-simple/services"
	"douyin-simple/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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
	// 检查token
	token := c.Query("token")
	tokenData, err := services.ResolveToken(token)
	if err != nil {
		OutPutGeneralResponse(c, 1, "无效token")
		return
	}
	user := tokenData.User

	// 获取title
	title := c.PostForm("title")

	// 将视频文件存入data
	data, err := c.FormFile("data")
	if err != nil {
		log.Println("视频文件放入data时出错：")
		utils.PrintLog(err, "[Error]")
		OutPutGeneralResponse(c, 1, "视频文件读取失败")
		return
	}

	// 获取文件名称与用户名
	filename := filepath.Base(data.Filename)
	// 如果title为空
	if title == "" {
		title = filename
	}
	// 打印相关信息
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	// 保存文件
	filePath := filepath.Join("/public/video/", finalName)
	if err := c.SaveUploadedFile(data, "."+filePath); err != nil {
		log.Println("保存文件时出错")
		utils.PrintLog(err, "[Error]")
		OutPutGeneralResponse(c, 1, "文件保存出错")
		return
	}

	playUrl := strings.Join([]string{BASE_URL, filePath}, "")
	coverPath := filepath.Join("/public/video_pic/", finalName)
	err = utils.ReadFrameAsJpeg("."+filePath, "."+coverPath)
	if err != nil {
		OutPutGeneralResponse(c, 1, "视频截帧出错")
		return
	}
	coverUrl := strings.Join([]string{BASE_URL, coverPath}, "")

	// 此处使用的是临时封面
	// coverUrl := "https://bkimg.cdn.bcebos.com/pic/80cb39dbb6fd5266d0168896e952802bd40735fa9855?x-bce-process=image/resize,m_lfit,w_536,limit_1/format,f_jpg"
	err = services.CreateVideo(title, playUrl, coverUrl, user.Id)
	if err != nil {
		OutPutGeneralResponse(c, 1, "视频信息存储出错")
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList get a publishing list of one user
func PublishList(c *gin.Context) {
	// 检查token
	token := c.Query("token")
	tokenData, err := services.ResolveToken(token)
	if err != nil {
		OutPutGeneralResponse(c, 1, "无效token")
		return
	}
	curUser := tokenData.User

	// 检查user_id
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		log.Println("目标用户的id数据有误：")
		utils.PrintLog(err, "[Error]")
		OutPutGeneralResponse(c, 1, "无效的user_id")
		return
	}

	// 获取指定用户信息
	author, err2 := services.GetUserResByUserId(userId, curUser.Id)
	if err2 != nil {
		OutPutGeneralResponse(c, 1, "用户信息查询失败或不存在该用户")
		return
	}

	videoResList, err := services.GetVideoResByUserId(author.Id, curUser.Id)
	if err != nil {
		OutPutGeneralResponse(c, 1, "视频列表获取失败")
		return
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		VideoList: videoResList,
	})

}
