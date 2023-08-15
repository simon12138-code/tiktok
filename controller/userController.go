/*
* @Author: zgy
* @Date:   2023/7/25 17:28
 */
package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_gin/dao"

	"go_gin/forms"
	"go_gin/models"
	"go_gin/response"
	"go_gin/utils"
	"golang.org/x/crypto/bcrypt"
	"time"
)

//无密码版本
//func PasswordLogin(c *gin.Context) {
//	PasswordLoginForm := forms.PasswordLoginForm{}
//	//不翻译错误信息版本
//	//ShouldBind能够基于请求的不同，自动提取JSON、form表单和QueryString类型的数据，并把值绑定到指定的结构体对象。
//	//if err := c.ShouldBind(&PasswordLoginForm); err != nil {
//	//	color.Blue(err.Error())
//	//
//	//	c.JSON(http.StatusInternalServerError, gin.H{
//	//		"err": err.Error(),
//	//	})
//	//	return
//	//}
//
//	//翻译错误信息类型版本
//	if err := c.ShouldBind(&PasswordLoginForm); err != nil {
//		//统一处理异常
//		utils.HandleValidatorError(c, err)
//		return
//	}
//	response.Success(c, 200, "success", "test")
//}

func PasswordLogin(c *gin.Context) {
	PasswordLoginForm := forms.PasswordLoginForm{}
	//参数校验
	if err := c.ShouldBind(&PasswordLoginForm); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	////验证条形码
	//if !store.Verify(PasswordLoginForm.CaptchaId, PasswordLoginForm.Captcha, true) {
	//	response.Err(c, 400, 400, "验证码错误", []interface{}{
	//		""})
	//}
	//response.Success(c, 200, "成功", "test")
	//encodepassword, _ := EncodePassword(PasswordLoginForm.PassWord)
	user, ok := dao.UsernameFindUserInfo(PasswordLoginForm.UserName, PasswordLoginForm.PassWord)
	//用户未注册直接返回
	if !ok {
		response.Err(c, 401, 401, "该用户未注册", "")
		return
	}
	//为用户生成token
	token := utils.CreateToken(c, user.ID, user.NickName, user.Role)
	//整理用户信息为统一返回封装的map格式（JSON）
	userinfoMap := HandleUserModelToMap(user)
	userinfoMap["token"] = token
	response.Success(c, 200, "success", userinfoMap)
}

// 接口获取用户列表
func GetUserList(c *gin.Context) {
	UserListForm := forms.UserListForm{}
	//首先进行获取参数，进行字段校验,
	if err := c.ShouldBind(&UserListForm); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	//获取数据
	total, userlist := dao.GetUserListDao(UserListForm.Page, UserListForm.PageSize)
	//判断
	if (total + len(userlist)) == 0 {
		response.Err(c, 400, 400, "未获取到数据", map[string]interface{}{
			"total":    total,
			"userlist": userlist,
		})
		return
	}
	response.Success(c, 200, "获取用户列表成功", map[string]interface{}{
		"total":    total,
		"userlist": userlist,
	})
}

// 接口创建用户
func InsertUser(c *gin.Context) {
	UserSubscribeForm := forms.UserSubscribeForm{}
	//参数校验
	if err := c.ShouldBind(&UserSubscribeForm); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	//如果没有填写生日，默认生成
	Birthdaytime := time.Date(2023, 7, 26, 12, 49, 0, 0, time.UTC)
	if UserSubscribeForm.Birthday == nil {
		UserSubscribeForm.Birthday = &Birthdaytime
	}
	//密码加密

	Password, err := EncodePassword(UserSubscribeForm.Password)
	if err != nil {
		panic(err)
		return
	}
	userInsert := &models.User{
		Password: Password,
		NickName: UserSubscribeForm.NickName,
		Birthday: UserSubscribeForm.Birthday,
		Address:  "",
		Desc:     "",
		Gender:   UserSubscribeForm.Gender,
		Role:     UserSubscribeForm.Role,
		Mobile:   UserSubscribeForm.Mobile,
	}
	//传入的是指针所以不需要接受，修改就在userInsert中
	_, bool := dao.UserCreate(userInsert)
	if !bool {
		response.Err(c, 500, 500, "插入失败", "")
		return
	}
	response.Success(c, 200, "创建成功", HandleUserModelToMap(userInsert))
}

func PutHeaderImage(c *gin.Context) {
	//获取对应的文件,key==”file“
	file, _ := c.FormFile("file")
	//File接口继承了流接口
	fileObj, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 把文件上传到minio对应的桶中
	ok := utils.UploadFile("userheader", file.Filename, fileObj, file.Size)
	if !ok {
		response.Err(c, 401, 401, "头像上传失败", "")
		return
	}
	headerUrl := utils.GetFileUrl("userheader", file.Filename, time.Second*24*60*60)
	if headerUrl == "" {
		response.Success(c, 400, "获取用户头像失败", "headerMap")
		return
	}
	//TODO 把用户的头像地址存入到对应user表中head_url 中
	response.Success(c, 200, "头像上传成功", map[string]interface{}{
		"userheaderUrl": headerUrl,
	})

}

// 功能函数:处理登录后的返回值
func HandleUserModelToMap(user *models.User) map[string]interface{} {
	birthday := ""
	if user.Birthday == nil {
		birthday = ""
	} else {
		birthday = user.Birthday.Format("2006-01-02")
	}
	userItemMap := map[string]interface{}{
		"id":        user.ID,
		"nick_name": user.NickName,
		"head_url":  user.HeadUrl,
		"birthday":  birthday,
		"address":   user.Address,
		"desc":      user.Desc,
		"gender":    user.Gender,
		"role":      user.Role,
		"mobile":    user.Mobile,
	}
	return userItemMap
}

// 功能函数，密码加密存储
func EncodePassword(password string) (string, error) {
	//加密方式：hash(password+randomnum)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		return password, err
	}
	encodePWD := string(hash)
	return encodePWD, nil
}
