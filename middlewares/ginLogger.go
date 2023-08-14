/*
* @Author: pzqu
* @Date:   2023/7/25 16:13
 */
package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// 传递闭包(匿名函数)
func Ginlogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		//链式调用前
		start := time.Now()
		//请求路径
		path := c.Request.URL.Path
		//请求参数
		query := c.Request.URL.RawQuery
		//链式调用
		c.Next()
		//链式调用后
		cost := time.Now().Sub(start)
		if c.Writer.Status() != 200 {
			//记录异常信息
			//这里zap.L()本质上调用了对应的已经替换的全局变量zap.Logger实例
			zap.L().Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				//一般指代使用的web端浏览器
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("cost", cost),
			)
		}
	}
}
