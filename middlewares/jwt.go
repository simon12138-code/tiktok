/*
* @Author: zgy
* @Date:   2023/7/26 11:07
 */
package middlewares

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"go_gin/global"
	"go_gin/response"
	"net/http"
)

type CustomClaims struct {
	ID int
	jwt.StandardClaims
}

// 无侵入式解析body-token
func JWTAuthInBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1、获取token
		token := c.PostForm("token")
		if token == "" {
			response.Err(c, http.StatusUnauthorized, response.Response{StatusMsg: "请登录", StatusCode: http.StatusUnauthorized})
			//先终止中间件调用
			c.Abort()
			return
		}
		//创建一个新的验证key
		j := NewJWT()
		//解析token
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				//token过期
				response.Err(c, http.StatusUnauthorized, response.Response{StatusMsg: "请登录", StatusCode: http.StatusUnauthorized})
				c.Abort()
				return
			}
			//其他错误
			response.Err(c, http.StatusUnauthorized, response.Response{StatusMsg: "未登录", StatusCode: http.StatusUnauthorized})
			c.Abort()
		}
		//打印上下文
		fmt.Println(c)
		// gin的上下文记录claims和userId的值
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}

}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1、获取token
		var token string

		token = c.Query("token")

		//color.Yellow(token)
		if token == "" {

			response.Err(c, http.StatusUnauthorized, response.Response{StatusMsg: "请登录", StatusCode: http.StatusUnauthorized})
			//先终止中间件调用
			c.Abort()
			return
		}
		//创建一个新的验证key
		j := NewJWT()
		//解析token
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				//token过期
				response.Err(c, http.StatusUnauthorized, response.Response{StatusMsg: "请登录", StatusCode: http.StatusUnauthorized})
				c.Abort()
				return
			}
			//其他错误
			response.Err(c, http.StatusUnauthorized, response.Response{StatusMsg: "未登录", StatusCode: http.StatusUnauthorized})
			c.Abort()
			return
		}
		//打印上下文
		//fmt.Println(c)
		// gin的上下文记录claims和userId的值
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

// 错误类型变量
var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.Settings.JWTKey.SigningKey),
	}
}

// 创建一个token，针对j的内部方法
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 创建key的解析方法
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		//校验错误类型。断言如果错误是token非法
		if ve, ok := err.(*jwt.ValidationError); ok {
			//根据校验错误类型判断,类似与字符串匹配
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		//token合法且是对应的声明
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}

}
