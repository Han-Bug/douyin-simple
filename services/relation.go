package services

import (
	"douyin-simple/models"
	"douyin-simple/utils"
	"log"
	"time"
)

func CreateRelation(followerId int64, userId int64) error {

	db, err := ConnectDatabase()
	if err != nil {
		return err
	}
	relation := models.Relation{
		FollowerId: followerId,
		UserId:     userId,
		CreatedAt:  time.Now(),
	}
	// 插入关注对象
	res := db.Create(&relation)
	if res.Error != nil {
		log.Println("关注数据创建出错")
		utils.PrintLog(res.Error, "[Warn]")
		return res.Error
	}
	return nil
}

func DeleteRelation(followerId int64, userId int64) error {
	db, err := ConnectDatabase()
	if err != nil {
		return err
	}
	res := db.Where("follower_id = ? AND user_id = ?", followerId, userId).Delete(&models.Relation{})
	if res.Error != nil {
		log.Println("关注数据删除出错")
		utils.PrintLog(res.Error, "[Warn]")
		return res.Error
	}
	return nil
}

// FindFollowUserList 获取关注列表
func FindFollowUserList(followerId int64, curUserId int64) ([]models.UserRes, error) {
	db, err := ConnectDatabase()
	if err != nil {
		return nil, err
	}
	var relations []models.Relation

	// 获取关系对象
	db.Where("follower_id = ?", followerId).Find(&relations)
	if db.Error != nil {
		log.Println("获取关系对象出错")
		utils.PrintLog(err, "[Error]")
		return nil, db.Error
	}

	// 遍历每个关系对象，获取用户数据
	var userList []models.UserRes
	for i := 0; i < len(relations); i++ {

		userRes, err := GetUserResByUserId(relations[i].UserId, curUserId)
		if err != nil {
			continue
		}
		userList = append(userList, userRes)
	}

	return userList, nil

}

// FindFollowerUserList 获取粉丝列表
func FindFollowerUserList(followId int64, curUserId int64) ([]models.UserRes, error) {
	db, err := ConnectDatabase()
	if err != nil {
		return nil, err
	}
	var relations []models.Relation

	// 获取关系对象
	db.Where("user_id = ?", followId).Find(&relations)
	if db.Error != nil {
		log.Println("获取关系对象出错")
		utils.PrintLog(err, "[Error]")
		return nil, db.Error
	}

	// 遍历每个关系对象，获取用户数据
	var userList []models.UserRes
	for i := 0; i < len(relations); i++ {

		userRes, err := GetUserResByUserId(relations[i].FollowerId, curUserId)
		if err != nil {
			continue
		}
		userList = append(userList, userRes)
	}

	return userList, nil

}
