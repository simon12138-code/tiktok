package models

type ChatContentIndex struct {
	UserId       int `json:"user_id" gorm:"primaryKey"`
	ToUserId     int `json:"to_user_id"`
	ContentIndex int `json:"content_index"`
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
