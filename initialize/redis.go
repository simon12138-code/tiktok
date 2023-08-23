/*
* @Author: pzqu
* @Date:   2023/7/26 10:18
 */
package initialize

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/redis/go-redis/v9"
	"go_gin/global"
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
