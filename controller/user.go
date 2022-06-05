package controller

import (
	"douyin-simple/models"
	"douyin-simple/units"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
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
var (
	dbdsn = "douyin:665577733_douYIN@tcp(119.23.68.131:3306)/douyin"
)

const (
	TOKEN_KEY           = "66557773"
	TOKEN_PREKEY        = "DOUYIN"
	TOKEN_PREKEY_LENGTH = 6
)

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

	// 先尝试登录
	user, err := GetUserModelByPwd(username, password)
	if err != nil {
		// TODO 将错误信息打印至日志中
		fmt.Println(err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "unexpected error occur when try login "},
		})
		return
	}
	if user != (models.User{}) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "User_Login already exist"},
		})
		return
	}
	// 进行注册
	newUser, err := DoRegister(username, password)
	if err != nil {
		// TODO 将错误信息打印至日志中
		fmt.Println(err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "unexpected error occur when register "},
		})
		return
	}

	token, err := MakeToken(username, password)
	if err != nil {
		// TODO 将错误信息打印至日志中
		fmt.Println(err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "unexpected error occur when making token "},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: models.Response{StatusCode: 0},
		UserId:   newUser.Id,
		Token:    token,
	})
}

// Login 登录接口
func Login(c *gin.Context) {

	// 获取输入信息
	username := c.Query("username")
	password := c.Query("password")

	//
	token := username + password

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "User_Login doesn't exist"},
		})
	}
}

func DoRegister(username string, password string) (newUser models.User, err error) {

	// &数据库连接
	db, err := gorm.Open(
		mysql.Open(dbdsn),
	)
	// TODO 将错误信息打印至日志中
	if err != nil {
		fmt.Println(err)
		return models.User{}, err
	}

	// TODO 数据有效性验证

	// &密码加密
	encPwd := units.EncryptMd5(password)
	// &创建用户对象，
	newUser = models.User{
		Name:     username,
		Password: encPwd,
	}
	// 执行Create
	res := db.Create(&newUser)
	if res.Error != nil {
		// TODO 将错误信息打印至日志中
		fmt.Println(res.Error)
		return models.User{}, err
	}

	return newUser, nil
}

// GetUserModelByPwd 通过帐号密码在数据库中查找相关用户，并返回用户数据
func GetUserModelByPwd(username string, password string) (user models.User, err error) {
	// 连接数据库
	db, err := gorm.Open(
		mysql.Open(dbdsn),
	)

	// TODO 将错误打印至日志中
	if err != nil {
		fmt.Println(err)
		return models.User{}, err
	}
	// 密码加密
	encPwd := units.EncryptMd5(password)

	user = models.User{}
	// 查询用户基本信息
	res := db.Model(&models.User{}).Where("user_name = ? AND password = ?", username, encPwd).First(&user)
	if res.Error != nil {
		// TODO 将错误打印至日志中
		fmt.Println(res.Error)
		return models.User{}, res.Error
	}

	return user, nil
}
func GetUserByUserModel(userModel models.User, curUserId int64) (user models.UserRes, err error) {
	// 连接数据库
	db, err := gorm.Open(
		mysql.Open(dbdsn),
	)
	// 查询用户关注总数
	var followCount int64
	db.Model(&models.Relation{}).Where("follower_id = ?", user.Id).Count(&followCount)
	// 查询用户粉丝总数
	var followerCount int64
	db.Model(&models.Relation{}).Where("user_id = ?", user.Id).Count(&followerCount)
	// 查询是否已关注
	var isRelatedN int64
	db.Model(&models.Relation{}).Where("user_id = ? AND follower_id = ?", user.Id, curUserId).Count(&isRelatedN)
	isRelated := false
	if isRelatedN > 0 {
		isRelated = true
	}
	user = models.UserRes{
		Id:            userModel.Id,
		Name:          userModel.Name,
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      isRelated,
	}
	return user, nil
}

func GetUserModelById(userId int64) (user models.User, err error) {
	// 连接数据库
	db, err := gorm.Open(
		mysql.Open(dbdsn),
	)

	if err != nil {
		// TODO 将错误打印至日志中
		fmt.Println(err)
		return models.User{}, err
	}

	user = models.User{}
	// 执行查询
	db.First(&user, userId)

	return user, nil
}

//GetUserModelByToken 根据token获取用户信息
func GetUserModelByToken(token string) (user models.User, err error) {
	// 解码token
	des, err := units.DecryptDes(token, []byte(TOKEN_KEY))
	if err != nil {
		// TODO 将错误打印至日志
		fmt.Println(err)
		return models.User{}, err
	}

	// 校验token前缀符
	sli := des[0:TOKEN_PREKEY_LENGTH]
	if strings.Compare(sli, TOKEN_PREKEY) != 0 {
		// TODO 将错误打印至日志
		fmt.Println("无效的token:", des)
		return models.User{}, errors.New("无效的token:" + des)
	}
	// 解析token
	var tokenData models.TokenData
	err = json.Unmarshal([]byte(des), &tokenData)
	if err != nil {
		// TODO 将错误打印至日志
		fmt.Println("token解析失败", des)
		return models.User{}, errors.New("token解析失败")
	}

	// 获取用户信息
	user, err = GetUserModelByPwd(tokenData.Name, tokenData.Password)
	if err != nil {
		return models.User{}, err
	}
	return user, err
}

func MakeToken(username string, password string) (token string, err error) {
	var tokenData = models.TokenData{
		Name:     username,
		Password: password,
		CreateAt: time.Now(),
	}
	data, err := json.Marshal(tokenData)
	if err != nil {
		return "", err
	}
	deToken := strings.Join([]string{TOKEN_PREKEY, string(data)}, "")
	des, err := units.EncryptDes(deToken, []byte(TOKEN_KEY))
	if err != nil {
		return "", err
	}
	return des, nil
}
func UserInfo(c *gin.Context) {
	curUserId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "bad id"},
		})
		return
	}
	token := c.Query("token")
	userModel, err := GetUserModelByToken(token)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "bad token"},
		})
		return
	}
	user, err := GetUserByUserModel(userModel, curUserId)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "error occurred when getting user"},
		})
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: models.Response{StatusCode: 0},
		User:     user,
	})

}
