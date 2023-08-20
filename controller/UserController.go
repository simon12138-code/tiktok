package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_gin/forms"
	"go_gin/global"
	"go_gin/models"
	"go_gin/response"
	"go_gin/service"
	"go_gin/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

// 实现了user和relation两个功能模块
// 用户登录
func UserLogin(c *gin.Context) {
	UserLoginForm := forms.UserLoginForm{}
	if err := c.ShouldBind(&UserLoginForm); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}

	userservice := service.NewUserService(c)
	data, msg, err := userservice.UserLoginService(UserLoginForm)
	//用户未注册直接返回
	if err != nil {
		response.Err(c, http.StatusUnauthorized, http.StatusUnauthorized, msg, "")
		return
	}

	response.Success(c, http.StatusOK, msg, data)

}

// 用户注册
func UserRegister(c *gin.Context) {
	UserLoginForm := forms.UserLoginForm{}
	if err := c.ShouldBind(&UserLoginForm); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	Password, err := hashAndSalt([]byte(UserLoginForm.Password))
	if err != nil {
		panic(err)
		return
	}
	UserLoginForm.Password = Password

	User := &models.User{
		UserName: UserLoginForm.Username,
		Password: UserLoginForm.Password,
	}
	userService := service.UserService{}

	data, msg, err := userService.CreateUserService(User)
	if err != nil {
		response.Err(c, 500, 500, msg, data)
		return
	}

	response.Success(c, 200, msg, data)
	return

}

// 获取用户信息
func GetUserInfo(c *gin.Context) {

	getUserInfo := forms.GetUserInfoForm{}

	if err := c.ShouldBind(&getUserInfo); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	UserId := getUserInfo.UserId

	NowUserId := c.Value("userId")
	log.Println(NowUserId, UserId)
	if NowUserId != UserId {
		response.Err(c, http.StatusBadRequest, http.StatusBadRequest, "用户不匹配", "")
		return
	}
	var userid string
	userid = string(UserId)

	data1, err := global.Redis.Get("user_info" + userid).Bytes()
	if err != nil {
		userService := service.UserService{}
		data, msg, err := userService.GetOneUserInfoService(UserId)
		if err != nil {
			response.Err(c, http.StatusInternalServerError, http.StatusInternalServerError, msg, data)
			return
		}
		//序列化结构体
		userJSON, err := json.Marshal(data)
		if err != nil {
			fmt.Println("JSON marshal error:", err)
			return
		}
		//存储结构体到redis
		err1 := global.Redis.Set("user_info"+userid, userJSON, time.Hour)
		if err1 != nil {
			log.Println("缓存失败")
		}
		response.Success(c, http.StatusOK, msg, data)
		return
	} else {
		//反序列化结构体
		var user models.User
		err = json.Unmarshal(data1, &user)
		if err != nil {
			fmt.Println("JSON unmarshal error:", err)
			return
		}
		response.Success(c, http.StatusOK, "查询成功", user)
		log.Println("查询缓存")
		return
	}

}

// 获取粉丝信息
func GetFollowerInfos(c *gin.Context) {
	getUserInfo := forms.GetUserInfoForm{}

	if err := c.ShouldBind(&getUserInfo); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	UserId := getUserInfo.UserId

	NowUserId := c.Value("userId").(int)
	if NowUserId != UserId {
		response.Err(c, http.StatusBadRequest, http.StatusBadRequest, "用户不匹配", "")
		return
	}
	userService := service.UserService{}

	followerIds, err := userService.GetFollowerIds(UserId)
	if err != nil {
		response.Success(c, http.StatusOK, "查询成功(无粉丝)", "")
		return
	}
	data, msg, err := userService.GetFollowInfoService(followerIds, NowUserId)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, http.StatusInternalServerError, msg, data)
		return
	}
	response.Success(c, http.StatusOK, msg, data)
	return
}

// 获取关注者信息
func GetFollowedUserInfos(c *gin.Context) {
	getUserInfo := forms.GetUserInfoForm{}

	if err := c.ShouldBind(&getUserInfo); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	UserId := getUserInfo.UserId

	NowUserId := c.Value("userId").(int)
	if NowUserId != UserId {
		response.Err(c, http.StatusBadRequest, http.StatusBadRequest, "用户不匹配", "")
		return
	}
	userService := service.UserService{}

	followedUserIds, err := userService.GetFollowedUserIds(UserId)
	if err != nil {
		response.Success(c, http.StatusOK, "查询成功（无关注者）", "")
		return
	}
	data, msg, err := userService.GetFollowInfoService(followedUserIds, NowUserId)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, http.StatusInternalServerError, msg, data)
		return
	}
	response.Success(c, http.StatusOK, msg, data)
	return
}

// 获取朋友列表
func GetFriendList(c *gin.Context) {
	getUserInfo := forms.GetUserInfoForm{}

	if err := c.ShouldBind(&getUserInfo); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	UserId := getUserInfo.UserId

	NowUserId := c.Value("userId")
	if NowUserId != UserId {
		response.Err(c, http.StatusBadRequest, http.StatusBadRequest, "用户不匹配", "")
		return
	}
	userService := service.UserService{}
	data, msg, err := userService.GetFriendListService(UserId)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, http.StatusInternalServerError, msg, err)
		return
	}
	response.Success(c, http.StatusOK, msg, data)
	return

}

// 用户执行关注/取关操作
func UserAction(c *gin.Context) {
	userActionForm := forms.ActionForm{}
	if err := c.ShouldBind(&userActionForm); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	ToUserId := userActionForm.ToUserId
	ActionType := userActionForm.ActionType
	UserId := c.Value("userId").(int)
	userService := service.UserService{}
	//关注
	if ActionType == 1 {
		data, msg, err := userService.UserFollowActionService(ToUserId, UserId)
		if err != nil {
			response.Err(c, http.StatusBadRequest, http.StatusBadRequest, msg, data)
			return
		}
		response.Success(c, http.StatusOK, msg, data)
		return
	}
	//取关
	if ActionType == 2 {
		data, msg, err := userService.UserCancelActionService(ToUserId, UserId)
		if err != nil {
			response.Err(c, http.StatusInternalServerError, http.StatusInternalServerError, msg, data)
			return
		}
		response.Success(c, http.StatusOK, msg, data)
		return
	}
	response.Err(c, http.StatusBadRequest, http.StatusBadRequest, "请求类型错误", "")
	return

}

// 功能函数，密码加密存储。

func hashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return string(hash), err
	}
	return string(hash), nil

}
