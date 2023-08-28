package router

import (
	"github.com/gin-gonic/gin"
	"go_gin/controller"
	"go_gin/middlewares"
)

func RelationRouter(Router *gin.RouterGroup) {

	RelationRouter := Router.Group("/relation")
	{
		RelationRouter.GET("/follower/list", middlewares.JWTAuth(), controller.GetFollowerInfos)
		RelationRouter.GET("/follow/list", middlewares.JWTAuth(), controller.GetFollowedUserInfos)
		RelationRouter.GET("/friend/list", middlewares.JWTAuth(), controller.GetFriendList)
		RelationRouter.POST("/action", middlewares.JWTAuth(), controller.UserAction)
	}
}
