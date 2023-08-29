package router

import (
	"github.com/gin-gonic/gin"
	"go_gin/controller"
)

func UtilRouter(Router *gin.RouterGroup) {

	UserRouter := Router.Group("")
	{
		UserRouter.POST("/head", controller.SendHead)
	}

}
