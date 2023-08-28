package router

import (
	"github.com/gin-gonic/gin"
	"go_gin/controller"
	"go_gin/middlewares"
)

func MessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message")
	{
		// 查询聊天记录
		MessageRouter.GET("chat", middlewares.JWTAuth(), controller.GetChatMessages)
		// 发送消息
		MessageRouter.POST("action", middlewares.JWTAuth(), controller.SendMessage)
	}
}
