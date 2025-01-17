package controller

import "douyin-simple/models"

var DemoVideos = []models.VideoRes{
	{
		Id:     1,
		Author: DemoUser,
		//PlayUrl: "https://www.w3schools.com/html/movie.mp4",
		PlayUrl: "http://172.27.157.250:8880/static/bear.mp4",
		//CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	},
}

var DemoComments = []models.CommentRes{
	{
		Id:         1,
		User:       DemoUser,
		Content:    "Test CommentRes",
		CreateDate: "05-01",
	},
}

var DemoUser = models.UserRes{
	Id:            1,
	Name:          "TestUser",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}
