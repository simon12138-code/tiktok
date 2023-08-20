/*
* @Author: zgy
* @Date:   2023/7/26 13:59
 */
package middlewares

import (
	"github.com/gin-gonic/gin"
)

// 原理:判断token的AuthorityId
func IsAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取token信息
		//获取现在的用户信息
		//判断现在的用户权限
		//继续执行下面中间件
		c.Next()
	}
}
