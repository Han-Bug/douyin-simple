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

type UserListResponse struct {
	models.Response
	UserList []models.UserRes `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	// 检查token
	token := c.Query("token")
	tokenData, err := services.ResolveToken(token)
	if err != nil {
		OutPutGeneralResponse(c, 1, "无效token")
		return
	}
	curUser := tokenData.User

	// 检查to_user_id
	targetUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		log.Println("RelationAction中解析to_user_id时出错：")
		utils.PrintLog(err, "[Error]")
		OutPutGeneralResponse(c, 1, "无效的to_user_id")
		return
	}

	actionType := c.Query("action_type")

	if actionType == "1" {
		if curUser.Id == targetUserId {
			OutPutGeneralResponse(c, 1, "不能关注自己")
			return
		}
		err := services.CreateRelation(curUser.Id, targetUserId)
		if err != nil {
			OutPutGeneralResponse(c, 1, "关注失败或已关注")
			return
		}
		c.JSON(http.StatusOK, models.Response{StatusCode: 0})

	} else if actionType == "2" {
		err := services.DeleteRelation(curUser.Id, targetUserId)
		if err != nil {
			OutPutGeneralResponse(c, 1, "取消关注失败或未关注")
			return
		}
		c.JSON(http.StatusOK, models.Response{StatusCode: 0})

	}
	OutPutGeneralResponse(c, 1, "无效的action_type")
}

// FollowList 关注列表
func FollowList(c *gin.Context) {
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
		log.Println("FollowList函数中解析user_id时出错：", err)
		utils.PrintLog(err, "[Error]")
		OutPutGeneralResponse(c, 1, "无效的user_id")
		return
	}

	userList, err := services.FindFollowUserList(userId, curUser.Id)
	if err != nil {
		OutPutGeneralResponse(c, 1, "获取关注列表失败")
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{StatusCode: 0},
		UserList: userList,
	})
}

// FollowerList 粉丝列表
func FollowerList(c *gin.Context) {
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
		log.Println("FollowList函数中解析user_id时出错：", err)
		utils.PrintLog(err, "[Error]")
		OutPutGeneralResponse(c, 1, "无效的user_id")
		return
	}

	userList, err := services.FindFollowerUserList(userId, curUser.Id)
	if err != nil {
		OutPutGeneralResponse(c, 1, "获取粉丝列表失败")
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{StatusCode: 0},
		UserList: userList,
	})
}
