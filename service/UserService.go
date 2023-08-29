package service

import (
	"github.com/gin-gonic/gin"
	"go_gin/dao"
	"go_gin/forms"
	"go_gin/models"
	"go_gin/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserService struct {
	ctx *gin.Context
}

func NewUserService(ctx *gin.Context) *UserService {
	return &UserService{ctx: ctx}
}
func (this UserService) UserLoginService(userLoginForm forms.UserLoginForm) (map[string]interface{}, interface{}, error) {
	userdb := dao.NewUserDB(this.ctx)
	user, err := userdb.UserLogin(userLoginForm.Username)
	if err != nil {
		return nil, "用户不存在", err
	}
	if comparePasswords(user.Password, []byte(userLoginForm.Password)) {
		//为用户生成token
		token := utils.CreateToken(this.ctx, user.Id)
		//整理用户信息为统一返回封装的map格式（JSON）
		userinfoMap := map[string]interface{}{
			"userId": user.Id,
			"token":  token,
		}
		return userinfoMap, "登录成功", nil
	}
	return nil, "密码错误", nil

}
func (this UserService) CreateUserService(user *models.User) (interface{}, interface{}, error) {
	userdb := dao.NewUserDB(this.ctx)
	_, err := userdb.UserCreate(user)
	if err != nil {
		return "", "db err", err
	}
	//为用户生成token
	token := utils.CreateToken(this.ctx, user.Id)
	//整理用户信息为统一返回封装的map格式（JSON）
	userinfoMap := map[string]interface{}{
		"userid": user.Id,
		"token":  token,
	}
	return userinfoMap, "注册成功", nil

}
func (this UserService) GetOneUserInfoService(userid int) (interface{}, interface{}, error) {
	userdb := dao.NewUserDB(this.ctx)
	userInfo, err := userdb.GetOneUserInfo(userid)
	if err != nil {
		return "", "db 错误", err

	}
	return userInfo, "查询成功", nil
}
func (this UserService) GetFollowInfoService(userids []int, userid int) (interface{}, interface{}, error) {
	userdb := dao.NewUserDB(this.ctx)
	userinfos, err := userdb.GetBranchUsers(userids, userid)
	if err != nil {
		return "", "db 错误", err

	}
	return userinfos, "查询成功", nil
}
func (this UserService) GetFollowerIds(userid int) ([]int, error) {
	userdb := dao.NewUserDB(this.ctx)
	followerIds, err := userdb.GetFollowerIds(userid)
	if err != nil {
		return followerIds, err
	}
	return followerIds, nil
}
func (this UserService) GetFollowedUserIds(userid int) ([]int, error) {
	userdb := dao.NewUserDB(this.ctx)
	followerIds, err := userdb.GetFollowedUserIds(userid)
	if err != nil {
		return followerIds, err
	}
	return followerIds, nil
}
func (this UserService) GetFriendListService(userid int) (interface{}, interface{}, error) {
	userdb := dao.NewUserDB(this.ctx)
	friendInfos, err := userdb.GetFriendList(userid)

	if err != nil {
		return "", "db 错误", err
	}
	if friendInfos == nil {
		return friendInfos, "查询成功(无朋友)", nil
	}
	return *friendInfos, "查询成功", nil

}
func (this UserService) UserFollowActionService(toUserid int, userid int) (interface{}, interface{}, error) {
	userdb := dao.NewUserDB(this.ctx)
	msg, err := userdb.UserActionFollow(toUserid, userid)
	if err != nil {
		return "", msg, err
	}
	return "", msg, err

}
func (this UserService) UserCancelActionService(Touserid int, userid int) (interface{}, interface{}, error) {

	userdb := dao.NewUserDB(this.ctx)
	msg, err := userdb.UserActionCancel(Touserid, userid)
	if err != nil {
		return "", msg, err
	}
	return "", msg, err

}

// 功能函数，解密
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
