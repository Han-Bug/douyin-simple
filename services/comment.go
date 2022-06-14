package services

import (
	"douyin-simple/models"
	"douyin-simple/utils"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

func CreateComment(user models.User, videoId int64, content string) (models.CommentRes, error) {
	// 数据库连接
	db, err := ConnectDatabase()
	if err != nil {
		return models.CommentRes{}, err
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
		log.Println("创建评论时出错：", db.Error)
		utils.PrintLog(db.Error, "[Error]")
		return models.CommentRes{}, db.Error
	}
	userRes, err := GetUserResByUser(user, user.Id)
	if err != nil {
		return models.CommentRes{}, err
	}
	commentRes := models.CommentRes{
		Id:         comment.Id,
		User:       userRes,
		Content:    content,
		CreateDate: comment.CreatedAt.String(),
	}
	return commentRes, nil
}

func DeleteComment(commentId int64) error {
	// 数据库连接
	db, _ := ConnectDatabase()
	// 删除评论
	db.Delete(&models.Comment{}, commentId)
	if db.Error != nil {
		fmt.Println("删除评论时出错：", db.Error)
		utils.PrintLog(db.Error, "[Error]")
		return db.Error
	}
	return nil
}

func GetCommentList(userId int64, videoId int64) ([]models.CommentRes, error) {
	db, err := ConnectDatabase()
	if err != nil {
		return nil, err
	}
	// 获取评论数据
	var comments []models.Comment
	db.Where("video_id = ?", videoId).Order("created_at desc").Find(&comments)
	if db.Error != nil {
		log.Println("评论获取出错")
		utils.PrintLog(db.Error, "[Error]")
		return nil, db.Error
	}
	// 加工评论数据
	var commentRes []models.CommentRes
	for i := 0; i < len(comments); i++ {
		// 获取评论作者信息
		var author models.User
		res := db.Where("user_id = ?", comments[i].UserId).First(&author)
		if res.Error != nil {
			log.Println("获取评论作者信息时出错")
			utils.PrintLog(res.Error, "[Error]")
		}
		authorRes, err := GetUserResByUser(author, userId)
		if err != nil {
			continue
		}
		commentRes = append(commentRes, models.CommentRes{
			Id:         comments[i].Id,
			User:       authorRes,
			Content:    comments[i].Content,
			CreateDate: comments[i].CreatedAt.String(),
		})
	}
	return commentRes, nil
}
