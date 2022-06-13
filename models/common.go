package models

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type VideoRes struct {
	Id            int64   `json:"id,omitempty"`
	Author        UserRes `json:"author"`
	PlayUrl       string  `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string  `json:"cover_url,omitempty"`
	FavoriteCount int64   `json:"favorite_count,omitempty"`
	CommentCount  int64   `json:"comment_count,omitempty"`
	IsFavorite    bool    `json:"is_favorite,omitempty"`
}

type CommentRes struct {
	Id         int64   `json:"id,omitempty"`
	User       UserRes `json:"user"`
	Content    string  `json:"content,omitempty"`
	CreateDate string  `json:"create_date,omitempty"`
}

type UserRes struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}
