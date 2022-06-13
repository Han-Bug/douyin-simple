package services

import (
	"douyin-simple/models"
	"douyin-simple/utils"
	"errors"
	"gorm.io/gorm"
	"log"
	"strings"
)

const (
	MD5_PREKEY = "douyin"
)

func CreateUser(username string, password string) (models.User, error) {

	db, err := ConnectDatabase(DSN)
	if err != nil {
		return models.User{}, err
	}

	// 密码加密
	encPwd := utils.EncryptMd5(strings.Join([]string{MD5_PREKEY, password}, ""))
	// 创建用户对象，
	newUser := models.User{
		Name:     username,
		Password: encPwd,
	}
	// 执行Create
	// 一个bug点：newUser在插入后获得的newUser返回值中其Id未更新
	res := db.Create(&newUser)
	if res.Error != nil {
		log.Println("创建用户信息时出错：")
		utils.PrintLog(res.Error, "[Error]")
		return models.User{}, res.Error
	}
	// 重新再获取用户信息
	res = db.Where("name = ? AND password = ?", username, encPwd).First(&newUser)
	if res.Error != nil {
		log.Println("重获取用户信息时出错：")
		utils.PrintLog(res.Error, "[Error]")
		return models.User{}, res.Error

	}
	return newUser, nil
}

// FindUserByPwd 通过帐号密码在数据库中查找相关用户，并返回用户数据
func FindUserByPwd(username string, password string) (user models.User, err error) {
	db, err := ConnectDatabase(DSN)
	if err != nil {
		return models.User{}, err
	}
	// 密码加密
	encPwd := utils.EncryptMd5(strings.Join([]string{MD5_PREKEY, password}, ""))
	user = models.User{}

	// 获取用户模型
	res := db.Model(&models.User{}).Where("name = ? AND password = ?", username, encPwd).First(&user)
	if res.Error != nil {
		// 如果不存在该用户
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return models.User{}, res.Error
		}
		log.Println("查询用户信息时出现错误:", res.Error)
		utils.PrintLog(err, "[Error]")
		return models.User{}, res.Error
	}

	return user, nil
}

// GetUserResByUser 获取加工后的用户信息，若当前未登录则curUserId值应为-1
func GetUserResByUser(user models.User, curUserId int64) (models.UserRes, error) {
	// 连接数据库
	db, err := ConnectDatabase(DSN)
	if err != nil {
		return models.UserRes{}, err
	}
	// 查询用户关注总数
	var followCount int64
	db.Model(&models.Relation{}).Where("follower_id = ?", user.Id).Count(&followCount)
	// 查询用户粉丝总数
	var followerCount int64
	db.Model(&models.Relation{}).Where("user_id = ?", user.Id).Count(&followerCount)
	// 查询是否已关注
	var isRelatedN int64
	// 如果用户未登录（curUser编号小于0）
	if curUserId < 0 {
		isRelatedN = -1
	} else {
		db.Model(&models.Relation{}).Where("user_id = ? AND follower_id = ?", user.Id, curUserId).Count(&isRelatedN)
	}
	isRelated := false
	if isRelatedN > 0 {
		isRelated = true
	}
	userRes := models.UserRes{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      isRelated,
	}
	return userRes, nil
}

func FirstUserById(userId int64) (models.User, error) {
	db, err := ConnectDatabase(DSN)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{}
	// 执行查询
	res := db.First(&user, userId)
	if res.Error != nil {
		log.Println("查询用户信息时出错或不存在该用户")
		utils.PrintLog(res.Error, "[Error]")
		return models.User{}, res.Error
	}

	return user, nil
}
func GetUserResByUserId(userId int64, curUserId int64) (models.UserRes, error) {
	// 连接数据库
	db, err := ConnectDatabase(DSN)
	if err != nil {
		return models.UserRes{}, err
	}

	// 获取用户数据
	user := models.User{}
	res := db.First(&user, userId)
	if res.Error != nil {
		log.Println("查询用户信息时出错或不存在该用户")
		utils.PrintLog(res.Error, "[Error]")
		return models.UserRes{}, res.Error
	}

	// 查询用户关注总数
	var followCount int64
	db.Model(&models.Relation{}).Where("follower_id = ?", user.Id).Count(&followCount)
	// 查询用户粉丝总数
	var followerCount int64
	db.Model(&models.Relation{}).Where("user_id = ?", user.Id).Count(&followerCount)
	// 查询是否已关注
	isRelated := false
	// 如果用户已登录（curUser编号不小于0）
	if curUserId >= 0 {
		var isRelatedN int64
		db.Model(&models.Relation{}).Where("user_id = ? AND follower_id = ?", user.Id, curUserId).Count(&isRelatedN)
		if isRelatedN > 0 {
			isRelated = true
		}

	}
	userRes := models.UserRes{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      isRelated,
	}
	return userRes, nil
}
