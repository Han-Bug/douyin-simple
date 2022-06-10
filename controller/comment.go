package controller

import (
	models "douyin-simple/models"
	"douyin-simple/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CommentActionResponse struct {
	models.Response
	Comment models.CommentRes `json:"comment,omitempty"`
}
type CommentListResponse struct {
	models.Response
	CommentList []models.CommentRes `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad video_id"})
		log.Println("CommentAction函数在解析video_id时出错：", err)
		utils.PrintLog(err, "[Error]")
		return
	}
	user, err := GetUserModelByToken(token)
	if err != nil {
		fmt.Println("CommentAction中获取用户信息时出错:", err)
		utils.PrintLog(err, "[Error]")
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad token"})
		return
	}
	actionType := c.Query("action_type")
	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad action_type"})
		return
	}
	// &数据库连接

	/*db, err := gorm.Open(
		mysql.Open(dbdsn),
	)
	// TODO 将错误信息打印至日志中
	if err != nil {
		fmt.Println("CommentAction中连接数据库出错：", err)
		utils.PrintLog(err, "[Error]")
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "database error"})
		return
	}*/

	db := ConnectDatabase(dbdsn, c)

	if actionType == "1" {
		// 发布评论
		content := c.PostForm("comment_text")
		// 如果评论为空
		if content == "" {
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "content cannot be empty"})
			return
		}
		// 创建评论
		comment := models.Comment{
			VideoId:   videoId,
			UserId:    user.Id,
			Content:   content,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		}
		db.Create(&comment)
		if db.Error != nil {
			fmt.Println("创建评论时出错：", err)
			utils.PrintLog(err, "[Error]")
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "error occur when creating comment"})
			return
		}
		userRes, err := GetUserByUserModel(user, user.Id)
		if err != nil {
			fmt.Println("获取用户信息时出错：", err)
			utils.PrintLog(err, "[Error]")
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "error occur when getting user"})
			return
		}
		commentRes := models.CommentRes{
			Id:         comment.Id,
			User:       userRes,
			Content:    content,
			CreateDate: comment.CreatedAt.String(),
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: models.Response{StatusCode: 0},
			Comment:  commentRes})
	} else if actionType == "2" {
		commentId, err := strconv.ParseInt(c.PostForm("comment_id"), 10, 64)
		if err != nil {
			fmt.Println("CommentAction中在解析comment_id时出错：", err)
			utils.PrintLog(err, "[Error]")
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad comment_id"})
			return
		}
		// 删除评论
		db.Delete(&models.Comment{}, commentId)
		if db.Error != nil {
			fmt.Println("删除评论时出错：", err)
			utils.PrintLog(err, "[Error]")
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "error occur when deleting comment"})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: models.Response{StatusCode: 0},
			Comment:  models.CommentRes{}})
	}

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		fmt.Println("解析video_id时出错：", err)
		utils.PrintLog(err, "[Error]")
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad video_id"})
		return
	}
	user, err := GetUserModelByToken(token)
	if err != nil {
		fmt.Println("CommentList中在通过token获取用户数据时出错：", err)
		utils.PrintLog(err, "[Error]")
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad token"})
		return
	}
	// &数据库连接
	//db, err := gorm.Open(
	//	mysql.Open(dbdsn),
	//)
	//// TODO 将错误信息打印至日志中
	//if err != nil {
	//	fmt.Println("连接数据库时出错：", err)
	//	utils.PrintLog(err, "[Fatal]")
	//	c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "database error"})
	//	return
	//}
	db := ConnectDatabase(dbdsn, c)
	// 获取评论列表
	var comments []models.Comment
	db.Where("video_id = ?", videoId).Order("created_at desc").Find(&comments)
	var commentReses []models.CommentRes
	for i := 0; i < len(comments); i++ {
		// 获取评论作者信息
		var author models.User
		db.Where("user_id = ?", comments[i].UserId).First(author)
		authorRes, err := GetUserByUserModel(author, user.Id)
		if err != nil {
			fmt.Println("获取用户信息时出错：", err)
			utils.PrintLog(err, "[Error]")
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "error occur when getting user"})
			return
		}
		commentReses = append(commentReses, models.CommentRes{
			Id:         comments[i].Id,
			User:       authorRes,
			Content:    comments[i].Content,
			CreateDate: comments[i].CreatedAt.String(),
		})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    models.Response{StatusCode: 0},
		CommentList: commentReses,
	})

}
