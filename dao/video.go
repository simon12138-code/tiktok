/*
* @Author: zgy
* @Date:   2023/8/14 12:16
 */
package dao

import (
	"errors"
	"go_gin/models"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)
import "go_gin/global"

var video models.Video
var videoInfo models.UserVideoInfo

type videoDB struct {
	ctx context.Context
}

func NewVideoDB(ctx context.Context) videoDB {
	return videoDB{ctx: ctx}
}
func (db videoDB) CreateVideoDao(video *models.Video) (int, error) {
	rows := global.DB.Create(video)
	if rows.RowsAffected < 1 {
		return video.VideoId, errors.New("db err")
	}
	return video.AuthorId, nil
}

func (db videoDB) CreateUserVideoInfoDao(userVideoInfo *models.UserVideoInfo) (int, error) {
	rows := global.DB.Create(userVideoInfo)
	if rows.RowsAffected < 1 {
		return userVideoInfo.UserId, errors.New("db err")
	}
	return userVideoInfo.UserId, nil
}
func (db videoDB) IncreaseUserVideoInfoWorkCount() error {
	userId := (db.ctx).Value("userId").(int)
	rows := global.DB.Model(&videoInfo).Where("user_id = ?", userId).Update("work_count", gorm.Expr("work_count + ?", 1))
	if rows.RowsAffected < 1 {
		return errors.New("db err")
	}
	return nil
}
func QueryUserVideoInfoDao(userId int) (*models.UserVideoInfo, error) {
	return &models.UserVideoInfo{}, nil
}
