package models

import "time"

type ChatContentIndex struct {
	Id           int        `json:"id" gorm:"primaryKey"`
	FromUserId   int        `json:"from_user_id"`
	ToUserId     int        `json:"to_user_id"`
	ContentIndex int        `json:"content_index"`
	CreateTime   *time.Time `json:"create_time"`
}

func (ChatContentIndex) TableName() string {
	return "chat_content_index"
}

type ChatContent struct {
	ContentId int    `json:"content_id" gorm:"primaryKey"`
	Content   string `json:"content"`
}

func (ChatContent) TableName() string {
	return "chat_content"
}
