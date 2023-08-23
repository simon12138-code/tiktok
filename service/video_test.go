/*
* @Author: zgy
* @Date:   2023/8/20 16:57
 */
package service

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go_gin/config"
	"go_gin/forms"
	"go_gin/global"
	"golang.org/x/net/context"
	"testing"
)

func InitRedis() {
	//拼接redis地址
	addr := fmt.Sprintf("%s:%d", global.Settings.Redisinfo.Host, global.Settings.Redisinfo.Port)
	//生成redis客户端
	global.Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0, //使用默认数据库
	})
	//连接redis数据库
	_, err := global.Redis.Ping(context.Background()).Result()
	//打印错误
	if err != nil {
		color.Red("[InitRedis] 链接redis异常:")
		color.Yellow(err.Error())
	}
}
func InitConfig() {
	//实例化Viper
	v := viper.New()
	//文件路径设置
	v.SetConfigFile("../settings-dev.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//创建配置变量
	serverConfig := config.ServerConfig{}
	//初始化 将绑定文件的配置信息反序列化到变量当中，完成文件信息到配置变量的转换
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	//传递全局变量 调用关系 settings-dev.yaml -> viper -> serverConfig(局部变量) -> global.Settings(全局变量)
	global.Settings = serverConfig
	color.Blue("11111111", global.Settings.LogAddress)
}

func Test(t *testing.T) {
	//启动配置
	InitConfig()
	//启动redis
	InitRedis()
	ctx := gin.Context{}
	video := NewVideoService(&ctx)
	form := forms.FeedForm{LatestTime: "", Token: ""}
	msg, videoList, err := video.FeedList(form)
	if err != nil {
		panic(err)
	}
	println(videoList, msg)
	//assert.Equal(t, res, true)
}
