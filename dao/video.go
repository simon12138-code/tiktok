/*
* @Author: zgy
* @Date:   2023/8/14 12:16
 */
package dao

import "go_gin/models"
import "go_gin/global"

var Video models.Video

func CreateVideoDao(video *models.Video) error {
	global.DB.Create(video)
}
