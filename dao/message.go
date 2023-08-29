package dao

import (
	"go_gin/forms"
	"go_gin/global"
	"go_gin/models"
	"strconv"
	"time"
)

// ChatContentCreate 向数据库中插入一条聊天记录
func ChatContentCreate(message *forms.Message) error {
	var chatContent models.ChatContent
	chatContent.Content = message.Content
	//向ChatContent表插入，并获得自动生成的主键
	rows1 := global.DB.Create(&chatContent)
	if rows1.Error != nil {
		return rows1.Error
	}
	var contentIndex models.ChatContentIndex
	contentIndex.FromUserId = message.FromUserId
	contentIndex.ToUserId = message.ToUserId
	contentIndex.ContentIndex = chatContent.ContentId
	timeNow := time.Now()
	contentIndex.CreateTime = &timeNow
	//将创建好的索引插入数据库
	rows2 := global.DB.Create(&contentIndex)
	if rows2.Error != nil {
		return rows1.Error
	}
	return nil
}

// GetMessageListIndex GetMessageList 查询聊天记录索引列表
func GetMessageListIndex(userid int, touserid int, timePre string) ([]models.ChatContentIndex, error) {
	preMessageTime, _ := strconv.ParseInt(timePre, 10, 64)
	timeNow := time.Now().Unix()
	timeP := time.Unix(preMessageTime, 0)
	if preMessageTime > timeNow {
		timeP = time.Unix(timeNow, 0)
	}
	indexList := []models.ChatContentIndex{}
	// 索引查询
	rows1 := global.DB.Model(models.ChatContentIndex{}).Where("from_user_id in ? AND create_time > ?", []int{userid, touserid}, timeP).Find(&indexList)
	if rows1.Error != nil {
		return nil, rows1.Error
	}

	return indexList, nil
}

// GetMessageList 根据索引列表查询聊天记录列表
func GetMessageList(indexList []models.ChatContentIndex) ([]forms.MessageRes, error) {
	messageList := make([]forms.MessageRes, len(indexList))
	// for循环根据索引查询content，依次插入messageList
	contentList := make([]models.ChatContent, len(indexList))
	for varIndex := range indexList {
		rows := global.DB.Model(models.ChatContent{}).Where("content_id = ?", indexList[varIndex].ContentIndex).Find(&contentList[varIndex])
		messageList[varIndex] = forms.MessageRes{
			Id:         indexList[varIndex].ContentIndex,
			ToUserId:   indexList[varIndex].ToUserId,
			FromUserId: indexList[varIndex].FromUserId,
			Content:    contentList[varIndex].Content,
			CreateTime: indexList[varIndex].CreateTime.Unix(),
		}

		if rows.Error != nil {
			return nil, rows.Error
		}
	}

	return messageList, nil
}
