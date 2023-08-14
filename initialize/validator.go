/*
* @Author: pzqu
* @Date:   2023/7/25 19:20
 */
package initialize

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"go_gin/global"
	"go_gin/utils"
	"reflect"
	"strings"
)

// 翻译器：给InitTrans传递一个参数,判断加载什么语言包,然后获取到语言包赋值给全局翻译器
func InitTrans(locale string) (err error) {
	color.Red("test")
	//修改gin框架的validator引擎属性，实现定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//注册一个获取jsontag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		//第一个入参为备用翻译器
		uni := ut.New(enT, enT, zhT)
		global.Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		}

		switch locale {
		case "en":
			_ = en_translations.RegisterDefaultTranslations(v, global.Trans)
		case "zh":
			_ = zh_translations.RegisterDefaultTranslations(v, global.Trans)
		default:
			_ = en_translations.RegisterDefaultTranslations(v, global.Trans)
		}
		//注册自定义校验器 入参(v 校验器实例，tag 对应的tag，提示信息 string， 自定义的校验函数)
		RegisterValidatorFunc(v, "mobile", "手机号码非法", utils.ValidateMobile)
		return
	}
	return
}

type Func func(fl validator.FieldLevel) bool

func RegisterValidatorFunc(v *validator.Validate, tag string, msgStr string, fn Func) {
	//注册tag自定义校验,validator.Func(fn)进行类型转换
	_ = v.RegisterValidation(tag, validator.Func(fn))
	//自定义错误内容，该部分修改子validator官方文档
	_ = v.RegisterTranslation(tag, global.Trans, func(ut ut.Translator) error {
		return ut.Add(tag, "0"+msgStr, true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
	return
}
