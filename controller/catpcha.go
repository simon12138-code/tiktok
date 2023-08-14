/*
* @Author: zgy
* @Date:   2023/7/26 10:26
 */
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"go_gin/response"
	"net/http"
)

// base64的默认缓存对象
var store = base64Captcha.DefaultMemStore

func GetCaptcha(ctx *gin.Context) {
	//生成对应验证码类型的驱动,NewDriverDigit可以根据参数调节验证码参数
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	//b64s是图片的base64编码
	id, b64s, err := cp.Generate()
	if err != nil {
		//日志记录
		zap.S().Errorf("生成验证码错误,:%s ", err.Error())
		//统一错误处理
		response.Err(ctx, http.StatusInternalServerError, 500, "生成验证码错误", "")
		return
	}
	response.Success(ctx, 200, "生成验证码成功", gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})
}
