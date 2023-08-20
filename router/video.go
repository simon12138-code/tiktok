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
		//登录功能
		VideoRouter.POST("/publish/action", middlewares.JWTAuth(), controller.Publish)

	}
}
