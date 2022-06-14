package services

import (
	"douyin-simple/models"
	"douyin-simple/utils"
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

func CreateFavorite(userId int64, videoId int64) error {
	db, err := ConnectDatabase()
	if err != nil {
		return err
	}
	favorite := models.Favorite{
		UserId:    userId,
		VideoId:   videoId,
		CreatedAt: time.Now(),
	}
	// 查询是否已点赞
	favoriteObj := models.Favorite{}
	res := db.Where("user_id = ? AND video_id = ?", userId, videoId).First(&favoriteObj)
	if res.Error != nil {
		// 如果未点赞
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// 插入点赞对象
			res := db.Create(&favorite)
			if res.Error != nil {
				log.Println("点赞数据创建出错")
				utils.PrintLog(res.Error, "[Warn]")
				return res.Error
			}
			return nil
		}
		// 如果出现意外错误
		log.Println("查询是否点赞时出错")
		utils.PrintLog(res.Error, "[Error]")
		return res.Error
	}
	// 如果已点赞
	err = errors.New("重复点赞")
	utils.PrintLog(err, "[Warn]")
	return err
}

func DeleteFavorite(userId int64, videoId int64) error {
	db, err := ConnectDatabase()
	if err != nil {
		return err
	}
	res := db.Where("video_id = ? AND user_id = ?", videoId, userId).Delete(&models.Favorite{})
	if res.Error != nil {
		log.Println("点赞数据删除出错")
		utils.PrintLog(res.Error, "[Warn]")
		return res.Error
	}
	return nil
}

func FindFavorite(videoId int64, userId int64) ([]models.Favorite, error) {
	db, err := ConnectDatabase()
	if err != nil {
		return nil, err
	}

	var favorite []models.Favorite
	res := db.Where("video_id = ? AND user_id = ?", videoId, userId).Limit(1).Find(&favorite)
	if res.Error != nil {
		log.Println("点赞数据查询出错")
		utils.PrintLog(res.Error, "[Warn]")
		return nil, res.Error
	}
	if len(favorite) == 0 {
		return nil, nil
	} else {
		return favorite, nil
	}
}

func FindFavoriteVideoList(targetUserId int64, curUserId int64) ([]models.VideoRes, error) {
	db, err := ConnectDatabase()
	if err != nil {
		return nil, err
	}
	var videos []models.VideoRes

	// 获取点赞视频列表
	var favorites []models.Favorite
	db.Where("user_id = ?", targetUserId).Order("created_at desc").Find(&favorites)
	for i := 0; i < len(favorites); i++ {
		// 获取视频信息
		videoRes, err := GetVideoResByVideoId(favorites[i].VideoId, curUserId)
		if err != nil {
			continue
		}
		videos = append(videos, videoRes)

	}
	return videos, nil
}
