/*
* @Author: pzqu
* @Date:   2023/7/25 21:50
 */
package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/* http状态码和code自定义状态码的区别:
http状态码是200,204,400,404,500 这样http规范定义的状态码
code只是个int类型数字,是和前端一起商量的状态码,比如 100010代表字段校验错误等等,是对http状态码一种详细的补充 */
// 封装c.JSON函数
func Success(c *gin.Context, res interface{}) {

	c.JSON(http.StatusOK, res)
	return
}

func Err(c *gin.Context, httpCode int, data interface{}) {
	c.JSON(httpCode, data)
}

type Response struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
type TestResponse struct {
	Response
	VideoList []int `json:"video_list"`
}
type ChatResponse struct {
	StatusCode string `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
