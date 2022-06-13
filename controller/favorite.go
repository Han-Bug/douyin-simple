package controller

import (
	"douyin-simple/models"
	"douyin-simple/services"
	"douyin-simple/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

//func ConnectDatabase(dbdsn string) (*gorm.DB, error) {
//	db, err := gorm.Open(mysql.Open(dbdsn), &gorm.Config{})
//	if err != nil {
//		fmt.Println("数据库连接失败：", err)
//		utils.PrintLog(err, "[Fatal]")
//		log.Fatalln(err)
//		return db, err
//	}
//	fmt.Println("数据库连接成功")
//	return db, nil
//}

/*
	如果是点赞功能的话：
	1.判断是否登录，未登录直接返回错误结果（通过修改StatusCode值，可以参考其他模块的代码）
	2.判断是点赞还是取消点赞
	若是点赞：查询是否已经点赞过（在数据库中根据视频编号和用户编号查询指定点赞数据，
	数据存在则表示点过赞了，不存在则未点赞），若未点赞则进行点赞（数据库中插入数据），
	否则返回直接返回结果（StatusCode值为0，即代表正常执行）；取消点赞类似

*/
//func ThumbUp(user models.User, db *gorm.DB, c *gin.Context) {
//	//查询用户编号
//	userId := user.Id
//	//获取视频编号
//	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
//	if err != nil {
//		fmt.Println("videoId转换成int类型出错")
//		utils.PrintLog(err, "[Error]")
//		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad token"})
//		return
//	}
//	//通过用户编号和视频编号查询视频是否出现在点赞列表
//	res := db.Model(&models.Favorite{}).Where("UserId = ? AND VideoId = ?", userId, videoId).Find(&user)
//	//对查询过程的错误进行处理
//	if res.Error != nil {
//		fmt.Printf("点赞操作在查询数据库时出错：%v", res.Error)
//		utils.PrintLog(err, "[Error]")
//		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad token"})
//		log.Fatalln(res.Error)
//		return
//	}
//	////定义一个结构体用于存储查询的数据
//	//favorite := models.Favorite{
//	//	UserId:    userId,
//	//	VideoId:   int64(videoId),
//	//	CreatedAt: time.Now(),
//	//}
//	////把点赞视频的数据插入到favorite表中
//	//res = db.Select("UserId", "VideoId", "CreatedAt").Create(&favorite)
//	////对插入结果进行错误处理
//	//if res.Error != nil {
//	//	fmt.Printf("点赞视频插入到数据库中出错：%v", res.Error)
//	//	utils.PrintLog(err, "[Error]")
//	//	c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad token"})
//	//	log.Fatalln(res.Error)
//	//}
//	services.CreateFavorite(userId, videoId)
//
//}

// CancelThumbUp 取消点赞操作
//func CancelThumbUp(db *gorm.DB, c *gin.Context) {
//	//根据videoID将记录从数据库删除
//	//获取videoId
//	videoIdstr := c.Query("video_id")
//	videoId, err := strconv.Atoi(videoIdstr)
//	if err != nil {
//		fmt.Println("videoId转换成int类型出错")
//		utils.PrintLog(err, "[Error]")
//		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad token"})
//
//	}
//	//执行删除操作
//	favorite := models.Favorite{
//		//UserId:    userId,
//		VideoId:   int64(videoId),
//		CreatedAt: time.Now(),
//	}
//	db.Where("video_id = ?", favorite.VideoId).Delete(&favorite)
//
//}

func FavoriteAction(c *gin.Context) {
	// 检查videoId
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		log.Println("videoId转换成int类型出错")
		utils.PrintLog(err, "[Error]")
		OutPutGeneralResponse(c, 1, "无效的videoId")
		return
	}

	// 检查token
	token := c.Query("token")
	tokenData, err := services.ResolveToken(token)
	if err != nil {
		OutPutGeneralResponse(c, 1, "无效的token")
		return
	}

	actionType := c.Query("action_type")

	if actionType == "1" {
		err := services.CreateFavorite(tokenData.User.Id, videoId)
		if err != nil {
			OutPutGeneralResponse(c, 1, "点赞创建失败或已点赞")
			return
		}
		c.JSON(http.StatusOK, models.Response{StatusCode: 0})

	} else if actionType == "2" {
		err := services.DeleteFavorite(tokenData.User.Id, videoId)
		if err != nil {
			OutPutGeneralResponse(c, 1, "点赞取消失败或未点赞")
			return
		}
		c.JSON(http.StatusOK, models.Response{StatusCode: 0})

	}
	OutPutGeneralResponse(c, 1, "无效的actionType")
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
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
	videoList, err := services.FindFavoriteVideoList(userId, curUser.Id)
	if err != nil {
		OutPutGeneralResponse(c, 1, "点赞视频列表获取失败")
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
