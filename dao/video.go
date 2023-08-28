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
	"strconv"
	"time"
)
import "go_gin/global"

var video models.Video
var videoInfo models.UserVideoInfo
var favorite models.Favorite

type videoDB struct {
	ctx context.Context
}

func NewVideoDB(ctx context.Context) *videoDB {
	return &videoDB{ctx: ctx}
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

func (db videoDB) GetVideoList(userId int) ([]models.Video, error) {
	videoList := []models.Video{}
	rows := global.DB.Model(video).Where("author_id = ?", userId).Find(&videoList)
	if rows.RowsAffected < 1 {
		return videoList, errors.New("db err")
	}
	return videoList, nil
}

func (db videoDB) GetUserVideoInfo(userId int) (*models.UserVideoInfo, error) {
	userInfo := &models.UserVideoInfo{}
	rows := global.DB.Model(userVideoInfo).Where("user_id = ?", userId).Find(userInfo)
	if rows.RowsAffected < 1 {
		return userInfo, errors.New("db err")
	}
	return userInfo, nil
}

func (db videoDB) GetUserVideoInfoList(userIds []int) ([]models.UserVideoInfo, error) {
	userInfoList := make([]models.UserVideoInfo, len(userIds))
	for i, v := range userIds {
		rows := global.DB.
			Model(userVideoInfo).
			//Clauses(clause.OrderBy{Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{userIds}, WithoutParentheses: true}}).
			Find(&userInfoList[i], v)
		if rows.RowsAffected < 1 {
			return userInfoList, errors.New("db err")
		}
	}

	return userInfoList, nil
}

func (db videoDB) CreateFavorite(favorite *models.Favorite) (int, bool, error) {
	rows := global.DB.Create(favorite)
	if rows.RowsAffected < 1 {
		return 0, false, errors.New("db err")
	}
	//更新作品点按数量
	var FavoriteVideoInfo models.Video
	rows2 := global.DB.Model(video).Where("video_id = ?", favorite.VideoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
	if rows2.RowsAffected < 1 {
		return 0, false, errors.New("db err")
	}
	rows3 := global.DB.Where("video_id = ?", favorite.VideoId).Find(&FavoriteVideoInfo)
	if rows3.RowsAffected < 1 {
		return 0, false, errors.New("db err")
	}
	//更新作者被点赞数量
	rows4 := global.DB.Model(&userVideoInfo).Where("user_id = ?", FavoriteVideoInfo.AuthorId).Update("favorited_count", gorm.Expr("favorited_count + ?", 1))
	if rows4.RowsAffected < 1 {
		return 0, false, errors.New("db err")
	}
	//更新用户点赞数量
	rows5 := global.DB.Model(&userVideoInfo).Where("user_id = ?", favorite.UserId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
	if rows5.RowsAffected < 1 {
		return 0, false, errors.New("db err")
	}
	return FavoriteVideoInfo.AuthorId, true, nil
}

func (db videoDB) DeleteFavorite(userfavorite *models.Favorite) (int, bool, error) {
	rows := global.DB.Where("user_id = ? AND video_id = ?", userfavorite.UserId, userfavorite.VideoId).Delete(favorite)
	if rows.RowsAffected < 1 {
		return 0, false, errors.New("db err")
	}
	//取消点赞的数量
	rows2 := global.DB.Model(&video).Where("video_id = ?", userfavorite.VideoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
	if rows2.RowsAffected < 1 {
		return 0, false, errors.New("db err")
	}

	var FavoriteVideoInfo models.Video

	rows3 := global.DB.Where("video_id = ?", userfavorite.VideoId).Find(&FavoriteVideoInfo)
	if rows3.RowsAffected < 1 {
		return 0, false, errors.New("db err")
	}
	//删除作者的被点赞数量
	rows4 := global.DB.Model(&userVideoInfo).Where("user_id = ?", FavoriteVideoInfo.AuthorId).Update("favorited_count", gorm.Expr("favorited_count - ?", 1))
	if rows4.RowsAffected < 1 {
		return 0, false, errors.New("db err")
	}
	//删除用户的点赞数量
	rows5 := global.DB.Model(&userVideoInfo).Where("user_id = ?", userfavorite.UserId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
	if rows5.RowsAffected < 1 {
		return 0, false, errors.New("db err")
	}
	return FavoriteVideoInfo.AuthorId, true, nil
}

func (db videoDB) GetFavoriteList(userId int) ([]models.Video, error) {
	//定义返回值
	videoList := []models.Video{}
	//先查询用户的喜欢列表
	favoriteList := []models.Favorite{}
	rows := global.DB.Model(favorite).Where("user_id = ? ", userId).Find(&favoriteList)
	if rows.Error != nil {
		return videoList, errors.New("db err")
	}
	//生成主键idlist
	idList := make([]int, len(favoriteList))
	for i, v := range favoriteList {
		idList[i] = v.VideoId
	}
	//根据用户喜欢列表进行查询
	rows2 := global.DB.Find(&videoList, idList)
	if rows2.RowsAffected < 1 {
		return videoList, errors.New("db err")
	}
	return videoList, nil
}

func (db videoDB) GetFeedVideoList(timestr string) ([]models.Video, error) {
	//根据时间戳逆序查询
	res := []models.Video{}
	timestamp, _ := strconv.ParseInt(timestr, 10, 64)
	t := time.Unix(timestamp, 0)
	rows := global.DB.Model(video).Where("create_time < ?", t).Order("create_time desc").Limit(global.MaxFeedCacheNum).Find(&res)
	if rows.Error != nil {
		return nil, errors.New("db err")
	}
	return res, nil
}

func (db videoDB) GetUserIsFavorite(userId string, videoList []int) ([]bool, error) {
	res := make([]bool, len(videoList))
	favorites := []models.Favorite{}
	userid, _ := strconv.Atoi(userId)
	rows := global.DB.
		Model(favorite).Where("user_id = ?", userid).
		//Clauses(clause.OrderBy{Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{userIds}, WithoutParentheses: true}}).
		Find(&favorites)
	if rows.Error != nil {
		return nil, errors.New("db err")
	}
	for i, v := range videoList {
		for _, v2 := range favorites {
			if v == v2.VideoId {
				res[i] = true
			}
		}
	}
	return res, nil

}
