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

	UserRouter := Router.Group("/user")
	{
		UserRouter.POST("/login/", controller.UserLogin)
		UserRouter.POST("/register/", controller.UserRegister)
		UserRouter.GET("/", middlewares.JWTAuth(), controller.GetUserInfo)
	}

}
