package controller

import (
	"github.com/gin-gonic/gin"
	"go_gin/dao"
	"go_gin/forms"
	"go_gin/models"
	"go_gin/response"
	"strconv"
)

func ActionChoice(c *gin.Context) {
	atype, _ := strconv.Atoi(c.Query("action_type"))
	switch atype {
	case 1:
		SendMessage(c)
	default:
		response.Err(c, 400, struct {
			response.Response
		}{
			response.Response{StatusCode: 400, StatusMsg: "未定义操作"},
		})
	}
}

// SendMessage 发送信息方法
func SendMessage(c *gin.Context) {
	mess := forms.Message{}

	//
	strid, _ := c.Get("userId")
	mess.FromUserId = strid.(int)

	strtoid := c.Query("to_user_id")
	mess.ToUserId, _ = strconv.Atoi(strtoid)

	mess.Content = c.Query("content")
	err := dao.ChatContentCreate(&mess)
	if err != nil {
		response.Err(c, 500, response.Response{StatusCode: 500, StatusMsg: err.Error()})
		return
	}
	response.Success(c, response.Response{StatusCode: 0, StatusMsg: "Success"})
}

// GetChatMessages 获取消息记录
func GetChatMessages(c *gin.Context) {
	userId, _ := c.Get("userId")

	strtoid := c.Query("to_user_id")
	tmptoid, _ := strconv.Atoi(strtoid)
	pre_msg_time := c.Query("pre_msg_time")

	var listIndex []models.ChatContentIndex
	listIndex, err1 := dao.GetMessageListIndex(userId.(int), tmptoid, pre_msg_time)
	if err1 != nil {
		response.Err(c, 500, response.ChatResponse{StatusCode: "500", StatusMsg: err1.Error()})
		return
	}

	var messageList []forms.MessageRes
	messageList, err2 := dao.GetMessageList(listIndex)
	if err2 != nil {
		response.Err(c, 500, response.ChatResponse{StatusCode: "500", StatusMsg: err2.Error()})
		return
	}

	response.Success(c, struct {
		response.ChatResponse
		MessageList []forms.MessageRes `json:"message_list"`
	}{
		response.ChatResponse{StatusCode: "0", StatusMsg: "Success"},
		messageList,
	})
}
