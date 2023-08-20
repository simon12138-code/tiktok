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

type VideoListForm struct {
	Token  string `form:"token" json:"token" binding:"required"`
	UserId string `form:"user_id" json:"user_id" binding:"required"`
}
