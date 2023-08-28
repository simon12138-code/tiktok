package controller

import (
	"github.com/gin-gonic/gin"
	"go_gin/dao"
	"go_gin/forms"
	"go_gin/response"
	"strconv"
)

// SendMessage 发送信息方法
func SendMessage(c *gin.Context) {
	mess := forms.Message{}

	atype, _ := strconv.Atoi(c.Query("action_type"))
	if atype != 1 {
		response.Err(c, 400, 400, "未定义消息操作", "")
	}

	//
	strid := c.Query("userId")
	mess.Id, _ = strconv.Atoi(strid)

	strtoid := c.Query("to_user_id")
	mess.ToUserId, _ = strconv.Atoi(strtoid)

	mess.Content = c.Query("content")

	dao.ChatContentCreate(&mess)
	response.Success(c, 200, "success", "")
}

// GetChatMessages 获取消息记录
func GetChatMessages(c *gin.Context) {
	strid := c.Query("userId")
	tmpid, _ := strconv.Atoi(strid)

	strtoid := c.Query("to_user_id")
	tmptoid, _ := strconv.Atoi(strtoid)

	var messageList []forms.Message
	messageList = dao.GetMessageList(tmpid, tmptoid)

	response.Success(c, 200, "success", messageList)
}
