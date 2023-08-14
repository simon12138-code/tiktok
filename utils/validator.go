/*
* @Author: zgy
* @Date:   2023/7/25 19:34
 */
package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go_gin/global"
	"go_gin/response"
	"net/http"
	"regexp"
	"strings"
)

// 处理字段校验异常
func HandleValidatorError(c *gin.Context, err error) {
	//如何返回错误信息，先来个断言
	errs, ok := err.(validator.ValidationErrors)
	//如果不是验证类型错误
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	//改为统一回复格式
	msg := removeTopStruct(errs.Translate(global.Trans))
	response.Err(c, http.StatusBadRequest, 400, "字段校验错误", msg)
	//c.JSON(http.StatusBadRequest, gin.H{
	//	//针对校验信息错误类型的方法Translate回调对应的翻译器global.Trans，返回的所有错误对应的翻译，键值对类型
	//	"error": removeTopStruct(errs.Translate(global.Trans))
	//})
}

// removeTopStruct 定义一个去掉结构体名称前缀的自定义方法：
/*
去掉前："PasswordLoginForm.mobile":"mobile为必填字段"
去掉后：					 "mobile":"mobile为必填字段"
*/
func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for filed, err := range fileds {
		//取到结构体名称前缀后的第一个"."字符下标，然后做一个字符串切片作为新的键
		rsp[filed[strings.Index(filed, ".")+1:]] = err
	}
	return rsp
}

func ValidateMobile(fl validator.FieldLevel) bool {
	//利用反射拿到结构体tag含有mobile的key字段
	mobile := fl.Field().String()
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if !ok {
		return false
	}
	return true
}
