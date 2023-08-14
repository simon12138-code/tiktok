/*
* @Author: zgy
* @Date:   2023/7/26 13:59
 */
package middlewares

import (
	"github.com/gin-gonic/gin"
	"go_gin/response"
)

// 原理:判断token的AuthorityId
func IsAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取token信息
		claims, _ := c.Get("claims")
		//获取现在的用户信息
		currentUser := claims.(*CustomClaims)

		//判断现在的用户权限
		if currentUser.AuthorityId != 1 {
			response.Err(c, 401, 401, "用户没有权限", "")
			//中断下面中间件
			c.Abort()
			return
		}
		//继续执行下面中间件
		c.Next()
	}
}
