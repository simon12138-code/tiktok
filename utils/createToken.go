/*
* @Author: zgy
* @Date:   2023/7/26 11:32
 */
package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go_gin/middlewares"
	"go_gin/response"
	"time"
)

// 针对对应ID和名称生成的token
func CreateToken(c *gin.Context, Id uint, NickName string, Role int) string {
	//生成token信息
	j := middlewares.NewJWT()
	//可以配合自定义封装想要的信息，在完成token验证后还可以通过上下文进行信息保存和验证其他权限等
	claims := middlewares.CustomClaims{
		//自定义
		ID:          uint(Id),
		NickName:    NickName,
		AuthorityId: uint(Role),
		//官方
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			// TODO 设置token过期时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //token -->30天过期
			Issuer:    "test",
		},
	}
	//生成token
	token, err := j.CreateToken(claims)
	if err != nil {
		response.Err(c, 401, 401, "token生成失败,重新再试", "test")
		return ""
	}
	return token

}
