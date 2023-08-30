/*
* @Author: zgy
* @Date:   2023/7/25 14:49
 */
package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"go_gin/global"
	"go_gin/initialize"
	"time"
)

func main() {
	//1、初始化配置
	initialize.InitConfig()
	//2、初始Routers
	Router := initialize.Routers()
	//3、初始化日志信息
	initialize.InitLogger()
	//4、初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		//该错误表示获取对应翻译器失败
		panic(err)
	}
	// 5.初始化mysql
	initialize.InitMysqlDB()
	//6. 初始化redis
	initialize.InitRedis()
	color.Cyan("go-gin服务开始了")
	global.Redis.Set(context.Background(), "test", "testValue", time.Second)
	////延迟两秒
	value := global.Redis.Get(context.Background(), "test")
	color.Blue(value.Val())
	// 7. 初始化minIO
	initialize.InitMinio()
	// 8、其中异步收集器
	initialize.InitColloctor()
	//9、启动url刷新定时任务
	initialize.InitCronJob()
	err := Router.Run(fmt.Sprintf(":%d", global.Settings.Port))
	if err != nil {
		zap.L().Info("this is hello func", zap.String("error", "启动错误!"))
	}
	//第二节内容
	//r := gin.Default()
	//r.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "pong",
	//	})
	//})
	//r.Run(fmt.Sprintf(":%d", global.Settings.Port))
}
