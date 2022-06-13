package controller

import (
	"douyin-simple/models"
	"douyin-simple/services"
	"douyin-simple/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	models.Response
	VideoList []models.VideoRes `json:"video_list,omitempty"`
	NextTime  int64             `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	// 检查token
	token := c.PostForm("token")
	var user models.User
	if token == "" {
		user = models.User{Id: -1}
	} else {
		tokenData, err := services.ResolveToken(token)
		if err != nil {
			log.Println("无效的token")
			utils.PrintLog(err, "[Info]")
			user = models.User{Id: -1}
		} else {
			user = tokenData.User
		}
	}

	// 检查latest_time
	latestTimeStr := c.PostForm("latest_time")
	var latestTime time.Time
	if latestTimeStr == "" {
		latestTime = time.Now()
	} else {

		ltNum, err := strconv.ParseInt(latestTimeStr, 10, 64)
		if err != nil {
			log.Println("无效的latest_time")
			utils.PrintLog(err, "[Info]")
			latestTime = time.Now()
		} else {
			latestTime = time.UnixMilli(ltNum)
		}
	}

	fmt.Println("latestTimeStr: ", latestTimeStr)
	fmt.Println("NowTime: ", time.Now())

	videoResList, nextTime, err := services.GetVideoResList(latestTime, user.Id)
	if err != nil {
		OutPutGeneralResponse(c, 1, "视频获取失败")
		return
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  models.Response{StatusCode: 0},
		VideoList: videoResList,
		NextTime:  nextTime.UnixMilli(),
	})
}
