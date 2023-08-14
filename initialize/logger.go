/*
* @Author: pzqu
* @Date:   2023/7/25 15:57
 */
package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"go_gin/global"
	"go_gin/utils"
)

// 初始化Logger
func InitLogger() {
	//实例化zap配置
	cfg := zap.NewDevelopmentConfig()
	//配置日志的输出地址
	cfg.OutputPaths = []string{
		fmt.Sprintf("%slog_%s.log", global.Settings.LogAddress, utils.GetNowFormatTodayTime()),
		"stdout",
	}
	logg, _ := cfg.Build()
	zap.ReplaceGlobals(logg)
	//注册到全局中
	global.Lg = logg
}
