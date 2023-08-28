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

		// 判断消息操作类型，并调用相关函数
		MessageRouter.POST("action", middlewares.JWTAuth(), controller.ActionChoice)
	}
}
