/*
* @Author: zgy
* @Date:   2023/8/14 19:59
 */
package router

import (
	"github.com/gin-gonic/gin"
	"go_gin/controller"
	"go_gin/middlewares"
)

func VideoRouter(router *gin.RouterGroup) {
	VideoRouter := router.Group("")
	{
		//发布功能
		VideoRouter.POST("/publish/action", middlewares.JWTAuth(), controller.Publish)
		//列表功能
		VideoRouter.POST("/publish/list", middlewares.JWTAuth(), controller.PublishList)
		//点赞功能
		VideoRouter.POST("/favorite/action", middlewares.JWTAuth(), controller.FavoriteAction)
		//喜欢列表
		VideoRouter.POST("/favorite/list", middlewares.JWTAuth(), controller.FavoriteAction)
		//视频推流
		VideoRouter.POST("/douyin/feed", controller.FavoriteAction)
	}
}
