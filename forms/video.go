/*
* @Author: zgy
* @Date:   2023/8/14 15:28
 */
package forms

import "mime/multipart"

type VideoForm struct {
	Data  *multipart.FileHeader
	Token string `form:"token" json:"token" binding:"required"`
	Title string `form:"title" json:"title" binding:"required"`
}
