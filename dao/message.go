package dao

import (
	"go_gin/forms"
	"go_gin/global"
	"go_gin/models"
)

// ChatContentCreate 向数据库中插入一条聊天记录
func ChatContentCreate(message *forms.Message) {
	var chatContent models.ChatContent
	chatContent.Content = message.Content
	//向ChatContent表插入，并获得自动生成的主键
	global.DB.Create(&chatContent)

	var contentIndex models.ChatContentIndex
	contentIndex.UserId = message.FromUserId
	contentIndex.ToUserId = message.ToUserId
	contentIndex.ContentIndex = chatContent.ContentId
	//将创建好的索引插入数据库
	global.DB.Create(&contentIndex)
}

// GetMessageList 查询聊天记录列表
func GetMessageList(userid int, touserid int) []forms.Message {

	messageList := make([]forms.Message, 0)
	indexList := make([]models.ChatContentIndex, 0)
	// 索引查询
	global.DB.Where("ID = ? AND ToUserId = ?", userid, touserid).Find(&indexList)

	// for循环根据索引查询content，依次插入messageList
	contentList := make([]models.ChatContent, 0)
	for varIndex := range indexList {
		global.DB.Where("ID = ?", indexList[varIndex].ContentIndex).Find(&contentList[varIndex])
		messageList[varIndex] = forms.Message{
			Id:         indexList[varIndex].ContentIndex,
			ToUserId:   userid,
			FromUserId: touserid,
			Content:    contentList[varIndex].Content}

	}

	return messageList
}
