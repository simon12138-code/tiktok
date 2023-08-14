/*
* @Author: pzqu
* @Date:   2023/7/25 15:44
 */
package router

import (
	"github.com/gin-gonic/gin"
	"go_gin/controller"
	"go_gin/middlewares"
)

func UserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	{
		UserRouter.GET("list", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message": "pong",
			})
		})
		//登录功能
		UserRouter.POST("login", controller.PasswordLogin)
		//列表获取功能
		UserRouter.POST("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), controller.GetUserList)
		//创建用户功能
		UserRouter.POST("create", controller.InsertUser)
		//添加头像功能
		UserRouter.POST("UploadUserHeaderImage", controller.PutHeaderImage)
	}
}
