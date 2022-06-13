package controller

import (
	models "douyin-simple/models"
	"douyin-simple/services"
	"douyin-simple/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type CommentActionResponse struct {
	models.Response
	Comment models.CommentRes `json:"comment,omitempty"`
}
type CommentListResponse struct {
	models.Response
	CommentList []models.CommentRes `json:"comment_list,omitempty"`
}

func OutPutGeneralResponse(c *gin.Context, statusCode int32, statusMsg string) {
	c.JSON(http.StatusOK, models.Response{StatusCode: statusCode, StatusMsg: statusMsg})
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	// 检查video_id
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		OutPutGeneralResponse(c, 1, "bad video_id")
		log.Println("CommentAction函数在解析video_id时出错：", err)
		utils.PrintLog(err, "[Error]")
		return
	}

	// 检查token
	token := c.Query("token")
	tokenData, err := services.ResolveToken(token)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad token"})
		return
	}
	user := tokenData.User

	// 获取动作类型
	actionType := c.Query("action_type")
	//if actionType != "1" && actionType != "2" {
	//	log.Println("无效的actionType")
	//	utils.PrintLog(errors.New("ActionType error"), "[Warn]")
	//	c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad action_type"})
	//	return
	//}
	if actionType == "1" {
		// 发布评论

		// 获取评论内容
		content := c.Query("comment_text")
		// 如果评论为空
		if content == "" {
			OutPutGeneralResponse(c, 1, "无效的评论内容")
			return
		}
		commentRes, err := services.CreateComment(user, videoId, content)
		// 如果评论创建失败
		if err != nil {
			OutPutGeneralResponse(c, 1, "评论创建失败")
			return
		}

		c.JSON(http.StatusOK, CommentActionResponse{
			Response: models.Response{StatusCode: 0},
			Comment:  commentRes})
	} else if actionType == "2" {
		// 删除评论

		// 获取评论id
		commentIdStr := c.Query("comment_id")
		commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
		if err != nil {
			log.Println("CommentAction中在解析comment_id时出错：", err)
			utils.PrintLog(err, "[Error]")
			OutPutGeneralResponse(c, 1, "无效的comment_id")
			return
		}
		err = services.DeleteComment(commentId)
		if err != nil {
			OutPutGeneralResponse(c, 1, "评论删除失败或评论不存在")
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: models.Response{StatusCode: 0},
			Comment:  models.CommentRes{}})
	}

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	// 检查videoId
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		log.Println("解析video_id时出错：", err)
		utils.PrintLog(err, "[Error]")
		OutPutGeneralResponse(c, 1, "无效的video_id")
		return
	}

	// 检查token
	token := c.Query("token")
	tokenData, err := services.ResolveToken(token)
	if err != nil {
		log.Println("CommentList中在通过token获取用户数据时出错：", err)
		utils.PrintLog(err, "[Error]")
		OutPutGeneralResponse(c, 1, "无效的token")
		return
	}
	user := tokenData.User

	// 获取评论
	commentList, err := services.GetCommentList(user.Id, videoId)
	if err != nil {
		OutPutGeneralResponse(c, 1, "评论获取失败")
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    models.Response{StatusCode: 0},
		CommentList: commentList,
	})

}
