/*
* @Author: zgy
* @Date:   2023/8/14 15:28
 */
package forms

import "mime/multipart"

type VideoForm struct {
	Data  *multipart.FileHeader
	Token string `form:"token" json:"token" binding:"required"`
	Title string `form:"title" json:"title" binding:"required"`
}

type VideoListForm struct {
	Token  string `form:"token" json:"token" binding:"required"`
	UserId string `form:"user_id" json:"user_id" binding:"required"`
}

type VideoFavcriteForm struct {
	Token      string `form:"token" json:"token" binding:"required"`
	VideoId    string `form:"video_id" json:"video_id" binding:"required"`
	ActionType string `form:"action_type" json:"action_type" binding:"required"`
}

type VideoFavoriteListForm struct {
	UserId string `form:"user_id" json:"user_id" binding:"required"` //用户id
	Token  string `form:"token" json:"token" binding:"required"`     //用户鉴权token
}

type FeedForm struct {
	LatestTime string `form:"latest_time" json:"latest_time"` //可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	Token      string `form:"token" json:"token"`             //用户登录状态下设置
}

type VideoFavoriteListRes PublishRes

type Author struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	FollowCount     int    `json:"follow_count"`
	FollowerCount   int    `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  string `json:"total_favorited"`
	WorkCount       int    `json:"work_count"`
	FavoriteCount   int    `json:"favorite_count"`
}
type PublishRes struct {
	VideoId       int    `json:"video_id" `
	Author        Author `json:"author"`
	PlayUrl       string `json:"play_url" `
	CoverUrl      string `json:"cover_url" `
	FavoriteCount int    `json:"favorite_count" `
	CommentCount  int    `json:"comment_count" `
	Title         string `json:"title"`
}

type FeedRes PublishRes
