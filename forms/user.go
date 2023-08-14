/*
* @Author: zgy
* @Date:   2023/7/25 17:20
 */
package forms

import "time"

//无条形码结构体
//type PasswordLoginForm struct {
//	//结构体tag自带校验参数
//	/*
//		通过 reflect.Type 获取结构体成员信息 reflect.StructField 结构中的 Tag 被称为结构体标签
//		tag都是静态的，无须实例化结构体，可以通过反射获取到
//		结构体标签由一个或多个键值对组成。键与值使用冒号分隔，值用双引号括起来。键值对之间使用一个空格分隔
//		编写 Tag 时，必须严格遵守键值对的规则，值内部不能添加空格
//	*/
//	PassWord string `form:"password" json:"password" binding:"required,min=3,max=20"`
//	UserName string `form:"name" json:"name" binding:"required,mobile"`
//}

// 有条形码信息的结构体
type PasswordLoginForm struct {
	PassWord  string `form:"password" json:"password" binding:"required,min=3,max=20"`
	UserName  string `form:"name" json:"name" binding:"required"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"` //验证码
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`       //验证码id
}

type UserListForm struct {
	Page     int `form:"page" json:"page" binding:"required"`
	PageSize int `form:"page_size" json:"page_size" binding:"required"`
}
type UserSubscribeForm struct {
	Password string     `form:"password" json:"password" binding:"required,min=3,max=20"`
	NickName string     `form:"name" json:"name" binding:"required"`
	Birthday *time.Time `form:"birthday" json:"birthday"`
	Address  string     `form:"address" json:"address" binding:"required"`
	Desc     string     `form:"desc" json:"desc"`
	Gender   string     `form:"gender" json:"gender" binding:"required"`
	Role     int        `form:"role" json:"role" binding:"required"`
	Mobile   string     `form:"mobile" json:"mobile" binding:"required,mobile"`
}
