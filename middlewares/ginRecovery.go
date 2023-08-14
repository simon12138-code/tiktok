/*
* @Author: pzqu
* @Date:   2023/7/25 16:31
 */
package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_gin/global"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
)

// 使用recovery接收调用链上的panic（防止系统宕机，并且输出到日志当中）
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		//defer调用链是先进后出,到整个调用链使用后，执行匿名函数
		defer func() {
			//recover接受panic错误
			err := recover()
			if err != nil {
				var brokenPipe bool
				//是否是网络调用错误
				if ne, ok := err.(*net.OpError); ok {
					//是否是系统调用错误
					if se, ok := ne.Err.(*os.SyscallError); ok {
						//检查系统差错是否包含 broken pipe和 connection reset by peer
						/*
							broken pipe:写入端出现的时候，另一端却休息或退出了，因此造成没有及时取走管道中的数据，从而系统异常退出
							connect reset by peer:连接被关闭,可能是在连接建立过程中三次握手失效等等，要具体分析，一般报错应用层不知道传输层的具体出错
						*/
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				//获取对应的resq
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				//如果出现系统错误
				if brokenPipe {
					global.Lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)

					//出现系统错误后，代表对短链接一般消失掉了，无法将错误信息写入status
					c.Error(err.(error))
					c.Abort()
					return
				}
				//非系统错误，不会造成崩溃的
				if stack {
					//记录请求
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				//终止请求，并且返回对应的状态码500
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		//在调用链之前拦截
		c.Next()
	}
}
