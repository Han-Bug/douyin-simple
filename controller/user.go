package controller

import (
	models "douyin-simple/models"
	"douyin-simple/services"
	"douyin-simple/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]models.UserRes{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	models.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	models.Response
	User models.UserRes `json:"user"`
}

// Register 注册接口
/*
	需考虑问题：
	用户名已存在
	用户名非法
	密码非法
	防SQL注入
*/
func Register(c *gin.Context) {
	// 获取输入数据
	username := c.Query("username")
	password := c.Query("password")

	// 查询是否已经注册
	_, err := services.FindUserByPwd(username, password)
	if err == nil {
		// 如果已经注册
		OutPutGeneralResponse(c, 1, "已存在该用户")
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果查询出错
		OutPutGeneralResponse(c, 1, "数据库异常")
		return
	}

	// 创建用户
	user, err := services.CreateUser(username, password)
	if err != nil {
		OutPutGeneralResponse(c, 1, "用户注册失败")
		return
	}
	// 生成token
	token, err := services.CreateToken(user)
	if err != nil {
		OutPutGeneralResponse(c, 1, "token创建失败")
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: models.Response{StatusCode: 0},
		UserId:   user.Id,
		Token:    token,
	})
}

// Login 登录接口
func Login(c *gin.Context) {
	// 获取输入信息
	username := c.Query("username")
	password := c.Query("password")
	// TODO 查询数据库前先做基本的数据合法性检查

	// 登录
	user, err := services.FindUserByPwd(username, password)
	if err != nil {
		// 如果未找到匹配用户
		if errors.Is(err, gorm.ErrRecordNotFound) {
			OutPutGeneralResponse(c, 1, "帐号与密码不匹配或不存在该帐号")
			return
		}
		// 如果查询出错
		OutPutGeneralResponse(c, 1, "用户信息获取出错")
		return
	}
	// 生成token
	token, err := services.CreateToken(user)
	if err != nil {
		OutPutGeneralResponse(c, 1, "token生成错误")
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: models.Response{StatusCode: 0},
		UserId:   user.Id,
		Token:    token,
	})

}

//GetUserModelByToken 根据token获取用户信息
// 当token解析失败或无效时err不为空
// 当err为空时必定返回解析后的用户
//func GetUserModelByToken(token string) (user models.User, err error) {
//
//	// 获取用户信息
//	user, err = GetUserModelByPwd(tokenData.Name, tokenData.Password)
//	if err != nil {
//		fmt.Println("GetUserModelByToken函数中获取用户信息时出错：", err)
//		utils.PrintLog(err, "[Error]")
//		return models.User{}, err
//	}
//	return user, err
//}

func UserInfo(c *gin.Context) {
	// 获取当前用户编号
	curUserId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		log.Println("UserInfo中解析user_id时出错：", err)
		utils.PrintLog(err, "[Error]")
		OutPutGeneralResponse(c, 1, "无效的user_id")
		return
	}
	// 获取token
	token := c.Query("token")
	// 解析token，获取用户信息
	tokenData, err := services.ResolveToken(token)
	if err != nil {
		OutPutGeneralResponse(c, 1, "token解析出错")
		return
	}
	// 获取用户详细信息
	user, err := services.GetUserResByUser(tokenData.User, curUserId)
	if err != nil {
		OutPutGeneralResponse(c, 1, "用户信息获取失败")
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: models.Response{StatusCode: 0},
		User:     user,
	})

}
