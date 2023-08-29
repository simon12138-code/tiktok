package controller

import (
	"github.com/gin-gonic/gin"
	"go_gin/forms"
	"go_gin/models"
	"go_gin/response"
	"go_gin/service"
	"go_gin/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
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

		response.Err(c, http.StatusUnauthorized, struct {
			response.Response
			data map[string]interface{}
		}{
			response.Response{StatusCode: http.StatusUnauthorized, StatusMsg: msg.(string)},
			map[string]interface{}{},
		})
		return
	}

	response.Success(c, struct {
		response.Response
		Id    int    `json:"user_id"`
		Token string `json:"token"`
	}{
		response.Response{StatusCode: 0, StatusMsg: msg.(string)},
		data["userId"].(int),
		data["token"].(string),
	})

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
		response.Err(c, 500, struct {
			response.Response
			data map[string]interface{}
		}{
			response.Response{StatusCode: http.StatusInternalServerError, StatusMsg: msg.(string)},
			map[string]interface{}{},
		})
		return
	}

	response.Success(c, struct {
		response.Response
		data map[string]interface{}
	}{
		response.Response{StatusCode: http.StatusUnauthorized, StatusMsg: msg.(string)},
		data.(map[string]interface{}),
	})
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
		response.Err(c, http.StatusBadRequest, struct {
			response.Response
		}{
			response.Response{StatusCode: http.StatusUnauthorized, StatusMsg: "用户不匹配"}})
		return
	}
	userService := service.UserService{}
	data, msg, err := userService.GetOneUserInfoService(UserId)

	if err != nil {
		response.Err(c, http.StatusBadRequest, struct {
			response.Response
			forms.UserRes `json:"user"`
		}{
			response.Response{StatusCode: http.StatusUnauthorized, StatusMsg: msg.(string)},
			forms.UserRes{},
		})
		return
	}
	response.Success(c, struct {
		response.Response
		forms.UserRes `json:"user"`
	}{
		response.Response{StatusCode: 0, StatusMsg: msg.(string)},
		*data.(*forms.UserRes),
	})

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
		response.Err(c, http.StatusBadRequest, struct {
			response.Response
		}{
			response.Response{StatusCode: http.StatusUnauthorized, StatusMsg: "用户不匹配"}})
		return
	}
	userService := service.UserService{}

	followerIds, err := userService.GetFollowerIds(UserId)
	if err != nil {
		response.Err(c, http.StatusBadRequest, struct {
			response.Response
		}{
			response.Response{StatusCode: 0, StatusMsg: "无粉丝"}})
		return
	}
	data, msg, err := userService.GetFollowInfoService(followerIds, NowUserId)
	if err != nil {
		response.Err(c, http.StatusBadRequest, struct {
			response.Response
		}{
			response.Response{StatusCode: http.StatusUnauthorized, StatusMsg: "用户不匹配"},
		})
		return
	}
	response.Success(c, struct {
		response.Response
		Data []forms.UserRes `json:"user_list"`
	}{
		response.Response{StatusCode: 0, StatusMsg: msg.(string)},
		data.([]forms.UserRes),
	})
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
		response.Err(c, http.StatusBadRequest, struct {
			response.Response
		}{
			response.Response{StatusCode: http.StatusUnauthorized, StatusMsg: "用户不匹配"},
		})
		return
	}
	userService := service.UserService{}

	followedUserIds, err := userService.GetFollowedUserIds(UserId)
	if err != nil {
		response.Success(c, struct {
			response.Response
		}{
			response.Response{StatusCode: 0, StatusMsg: "无关注者"},
		})
		return
	}
	data, msg, err := userService.GetFollowInfoService(followedUserIds, NowUserId)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, struct {
			response.Response
		}{
			response.Response{StatusCode: http.StatusInternalServerError, StatusMsg: msg.(string)},
		})
		return
	}
	response.Success(c, struct {
		response.Response
		Data []forms.UserRes `json:"user_list"`
	}{
		response.Response{StatusCode: http.StatusUnauthorized, StatusMsg: msg.(string)},
		data.([]forms.UserRes),
	})
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
	//该接口支持查询别人的好友列表，所以直接选取userId作为查询key
	//NowUserId := c.Value("userId")
	//if NowUserId != UserId {
	//	response.Err(c, http.StatusBadRequest, struct {
	//		response.Response
	//	}{
	//		response.Response{StatusCode: http.StatusUnauthorized, StatusMsg: "用户不匹配"},
	//	})
	//	return
	//}
	userService := service.UserService{}
	data, msg, err := userService.GetFriendListService(UserId)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, struct {
			response.Response
		}{
			response.Response{StatusCode: http.StatusInternalServerError, StatusMsg: msg.(string)},
		})
		return
	}
	response.Success(c, struct {
		response.Response
		Data []forms.FriendRes `json:"user_list"`
	}{
		response.Response{StatusCode: 0, StatusMsg: msg.(string)},
		data.([]forms.FriendRes),
	})
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
		_, msg, err := userService.UserFollowActionService(ToUserId, UserId)
		if err != nil {
			response.Err(c, http.StatusBadRequest, struct {
				response.Response
			}{
				response.Response{StatusCode: http.StatusUnauthorized, StatusMsg: msg.(string)},
			})
			return
		}
		response.Success(c, struct {
			response.Response
		}{
			response.Response{StatusCode: 0, StatusMsg: msg.(string)},
		})
		return
	}
	//取关
	if ActionType == 2 {
		_, msg, err := userService.UserCancelActionService(ToUserId, UserId)
		if err != nil {
			response.Err(c, http.StatusBadRequest, struct {
				response.Response
			}{
				response.Response{StatusCode: http.StatusUnauthorized, StatusMsg: msg.(string)},
			})
			return
		}
		response.Success(c, struct {
			response.Response
		}{
			response.Response{StatusCode: 0, StatusMsg: msg.(string)},
		})
		return
	}
	response.Err(c, http.StatusBadRequest, struct {
		response.Response
	}{
		response.Response{StatusCode: http.StatusUnauthorized, StatusMsg: "请求类型错误"},
	})
	return

}

// 功能函数，密码加密存储

func hashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return string(hash), err
	}
	return string(hash), nil

}
