/*
* @Author: pzqu
* @Date:   2023/7/25 15:48
 */
package initialize

import (
	"github.com/gin-gonic/gin"
	"go_gin/middlewares"
	"go_gin/router"
)

func Routers() *gin.Engine {
	//获取固定路由
	Router := gin.Default()
	//分配路由。装配顺序： / + v1/ + user/ + list
	ApiGroup := Router.Group("/douyin/")
	//路由分组
	router.VideoRouter(ApiGroup)
	router.UserRouter(ApiGroup)
	router.RelationRouter(ApiGroup)

	//添加中间件（全局中间件）
	Router.Use(middlewares.Ginlogger(), middlewares.GinRecovery(true))
	Router.Use(middlewares.Cors())
	return Router
}
