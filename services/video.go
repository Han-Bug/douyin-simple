package services

import (
	"douyin-simple/models"
	"douyin-simple/utils"
	"gorm.io/gorm"
	"log"
	"time"
)

func CreateVideo(title string, playUrl string, coverUrl string, userId int64) error {
	db, err := ConnectDatabase()
	if err != nil {
		return err
	}

	video := models.Video{
		UserId:    userId,
		Title:     title,
		PlayUrl:   playUrl,
		CoverUrl:  coverUrl,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}
	res := db.Create(&video)
	if res.Error != nil {
		utils.PrintLogError("视频信息存储出错:", "video:", video, "err:", res.Error)
		return res.Error
	}
	return nil
}
func FirstVideo(videoId int64) (models.Video, error) {
	db, err := ConnectDatabase()
	if err != nil {
		return models.Video{}, err
	}

	video := models.Video{}
	res := db.First(&video, videoId)
	if res.Error != nil {
		utils.PrintLogWarn("查询视频信息时出错或不存在该视频:", "videoId=", videoId, " err:", res.Error)
		return models.Video{}, res.Error
	}

	return video, nil
}

func GetVideoResByVideoId(videoId int64, curUserId int64) (models.VideoRes, error) {
	db, err := ConnectDatabase()
	if err != nil {
		return models.VideoRes{}, err
	}

	video := models.Video{}
	res := db.First(&video, videoId)
	if res.Error != nil {
		log.Println("视频查询出错或视频不存在")
		utils.PrintLog(res.Error, "[Warn]")
		return models.VideoRes{}, res.Error
	}
	videoRes, err := GetVideoResByVideoWithDB(video, curUserId, db)
	if err != nil {
		return models.VideoRes{}, err
	}
	return videoRes, nil
}

func GetVideoResByUserId(userId int64, curUserId int64) ([]models.VideoRes, error) {
	db, err := ConnectDatabase()
	if err != nil {
		return nil, err
	}

	var videos []models.Video
	db.Where("user_id = ?", userId).Order("created_at desc").Find(&videos)
	var videoResList []models.VideoRes
	for i := 0; i < len(videos); i++ {

		videoRes, err := GetVideoResByVideoWithDB(videos[i], curUserId, db)
		if err != nil {
			continue
		}
		videoResList = append(videoResList, videoRes)

	}
	return videoResList, nil
}

func GetVideoResList(latestTime time.Time, curUserId int64) ([]models.VideoRes, time.Time, error) {
	db, err := ConnectDatabase()
	if err != nil {
		return nil, time.Now(), err
	}

	// 获取视频数据
	var videos []models.Video
	res := db.Where(" created_at < ?", latestTime).Order("created_at desc").Limit(5).Find(&videos)
	if res.Error != nil {
		log.Println("GetVideoResList:获取视频列表失败")
		utils.PrintLog(res.Error, "[Error]")
		return nil, time.Now(), res.Error
	}

	// 为每个视频填充详细信息
	var videoResList []models.VideoRes
	for i := 0; i < len(videos); i++ {
		videoRes, err := GetVideoResByVideoWithDB(videos[i], curUserId, db)
		if err != nil {
			continue
		}
		videoResList = append(videoResList, videoRes)
	}
	nextTime := time.Now()
	if len(videos) > 0 {
		nextTime = videos[len(videos)-1].CreatedAt
	}
	return videoResList, nextTime, nil
}

func GetVideoResByVideoWithDB(video models.Video, curUserId int64, db *gorm.DB) (models.VideoRes, error) {
	// 查询视频的点赞数、评论数
	var favoriteCount int64
	db.Model(&models.Favorite{}).Where("video_id = ?", video.Id).Count(&favoriteCount)
	var commentCount int64
	db.Model(&models.Comment{}).Where("video_id = ?", video.Id).Count(&commentCount)
	// 查询该用户是否赞了该视频
	var isFavorite = false
	if curUserId > 0 {
		fav, err := FindFavorite(video.Id, curUserId)
		if fav == nil || err != nil {
			isFavorite = true
		}
	}
	// 获取视频作者信息
	authorRes, err := GetUserResByUserId(video.UserId, curUserId)
	if err != nil {
		return models.VideoRes{}, err
	}
	videoRes := models.VideoRes{
		Id:            video.Id,
		Author:        authorRes,
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: favoriteCount,
		CommentCount:  commentCount,
		IsFavorite:    isFavorite,
	}
	return videoRes, nil

}
