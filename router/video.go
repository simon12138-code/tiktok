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
		//发布功能,目前采用侵入式鉴权解决方案
		VideoRouter.POST("/publish/action/", middlewares.JWTAuthInBody(), controller.Publish)
		//列表功能
		VideoRouter.GET("/publish/list/", middlewares.JWTAuth(), controller.PublishList)
		//点赞功能
		VideoRouter.POST("/favorite/action/", middlewares.JWTAuth(), controller.FavoriteAction)
		//喜欢列表
		VideoRouter.GET("/favorite/list/", middlewares.JWTAuth(), controller.FavoriteList)
		//视频推流
		VideoRouter.GET("/feed/", middlewares.JWTAuth(), controller.Feed)
	}
}
