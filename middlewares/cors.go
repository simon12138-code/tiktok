/*
* @Author: pzqu
* @Date:   2023/7/26 13:53
 */
package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取请求方法
		method := c.Request.Method
		//请求头添加相应的允许跨域内容
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		//如果是对应的option方法则直接返回204
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

	}
}
