/*
* @Author: zgy
* @Date:   2023/8/14 11:40
 */
package models

import "time"

type Video struct {
	VideoId       int        `json:"video_id" `
	AuthorId      int        `json:"author_id" `
	PlayUrl       string     `json:"play_url" `
	CoverUrl      string     `json:"cover_url" `
	FavoriteCount int        `json:"favorite_count" `
	CommentCount  int        `json:"comment_count" `
	Title         string     `json:"title"`
	CreateTime    *time.Time `json:"create_time" gorm:"type:date"`
}

func (Video) TableName() string {
	return "video"
}

type Favorite struct {
	ID      int `json:"id" gorm:"primaryKey"`
	UserId  int `json:"user_id" `
	VideoId int `json:"video_id" `
}

func (Favorite) TableName() string {
	return "favorite"
}

type UserVideoInfo struct {
	UserId         int `json:"user_id" gorm:"prmaryKey"`
	FavoriteCount  int `json:"favorite_count"`
	FavoritedCount int `json:"favorited_count"`
	WorkCount      int `json:"work_count"`
}

func (UserVideoInfo) TableName() string {
	return "user_video_info"
}
