package controller

import (
	"douyin-simple/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type UserListResponse struct {
	models.Response
	UserList []models.UserRes `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	targetUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad to_user_id"})
		return
	}
	curUser, err := GetUserModelByToken(token)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad token"})
		return
	}
	actionType := c.Query("action_type")
	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad action_type"})
		return
	}
	// &数据库连接
	db, err := gorm.Open(
		mysql.Open(dbdsn),
	)
	// TODO 将错误信息打印至日志中
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "database error"})
		return
	}
	// 查看是否已经关注
	isRelated := false
	var relations []models.Relation
	db.Where("user_id = ? AND follower_id = ?", curUser.Id, targetUserId).Find(relations)
	if len(relations) > 0 {
		isRelated = true
	}

	if actionType == "1" {
		// 关注
		if isRelated {
			fmt.Println("重复关注")
			c.JSON(http.StatusOK, models.Response{StatusCode: 0})
			return
		} else {
			db.Create(&models.Relation{
				FollowerId: curUser.Id,
				UserId:     targetUserId,
				CreatedAt:  time.Time{},
			})
			if db.Error != nil {
				c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "error occur when creating relation"})
				return
			}
			c.JSON(http.StatusOK, models.Response{StatusCode: 0})
			return
		}

	} else if actionType == "2" {
		// 取消关注
		if isRelated {
			db.Where("follower_id = ? AND user_id = ?", curUser.Id, targetUserId).Delete(&models.Relation{})
			if db.Error != nil {
				c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "error occur when creating relation"})
				return
			}
			c.JSON(http.StatusOK, models.Response{StatusCode: 0})
			return

		} else {
			fmt.Println("取消不存在的关注")
			c.JSON(http.StatusOK, models.Response{StatusCode: 0})
			return
		}
	}
}

// FollowList 关注列表
func FollowList(c *gin.Context) {
	token := c.Query("token")
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad user_id"})
		return
	}
	curUser, err := GetUserModelByToken(token)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad token"})
		return
	}
	// &数据库连接
	db, err := gorm.Open(
		mysql.Open(dbdsn),
	)
	// TODO 将错误信息打印至日志中
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "database error"})
		return
	}
	var relations []models.Relation
	// 获取关系对象
	db.Where("follower_id = ?", userId).Find(&relations)
	if db.Error != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "error occur when getting relations"})
		return
	}
	var userList []models.UserRes
	for i := 0; i < len(relations); i++ {
		user, err := GetUserModelById(relations[i].UserId)
		if err != nil {
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "error occur when getting user"})
			return
		}
		userRes, err := GetUserByUserModel(user, curUser.Id)
		if err != nil {
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "error occur when getting user"})
			return
		}
		userList = append(userList, userRes)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{StatusCode: 0},
		UserList: userList,
	})
}

// FollowerList 粉丝列表
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad user_id"})
		return
	}
	curUser, err := GetUserModelByToken(token)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "bad token"})
		return
	}
	// &数据库连接
	db, err := gorm.Open(
		mysql.Open(dbdsn),
	)
	// TODO 将错误信息打印至日志中
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "database error"})
		return
	}
	var relations []models.Relation
	// 获取关系对象
	db.Where("user_id = ?", userId).Find(&relations)
	if db.Error != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "error occur when getting relations"})
		return
	}
	var userList []models.UserRes
	for i := 0; i < len(relations); i++ {
		user, err := GetUserModelById(relations[i].FollowerId)
		if err != nil {
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "error occur when getting user"})
			return
		}
		userRes, err := GetUserByUserModel(user, curUser.Id)
		if err != nil {
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "error occur when getting user"})
			return
		}
		userList = append(userList, userRes)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{StatusCode: 0},
		UserList: userList,
	})
}
