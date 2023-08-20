/*
* @Author: zgy
* @Date:   2023/8/14 14:39
 */
package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_gin/dao"
	"go_gin/forms"
	"go_gin/global"
	"go_gin/models"
	"go_gin/utils"
	"image/jpeg"
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"
)

type VideoService struct {
	ctx *gin.Context
}

func NewVideoService(ctx *gin.Context) *VideoService {
	return &VideoService{ctx: ctx}
}

func (this VideoService) Pubish(videoFrom forms.VideoForm) (interface{}, interface{}, error) {
	fileObj, err := videoFrom.Data.Open()
	if err != nil {
		global.Lg.Error(err.Error())
		return "", "", err
	}
	userId := this.ctx.Value("userId").(int)
	finalName := fmt.Sprintf("%d__%s", userId, videoFrom.Data.Filename)
	//设定路径public文件夹下
	saveFile := filepath.Join("../public/", finalName)
	//保存文件
	if err = this.ctx.SaveUploadedFile(videoFrom.Data, saveFile); err != nil {
		global.Lg.Error(err.Error())
		return "保存文件失败", "", err
	}
	snapShotName := finalName + "-cover.jpeg"
	img, err := utils.GetSnapShot(snapShotName, saveFile)
	//获取视频的url
	video_url, err := uploadAndGetUrl("video", finalName, fileObj, videoFrom.Data)
	if err != nil {
		global.Lg.Error(err.Error())
		return "上传失败", "", err
	}
	//创建图片缓存
	var buffer bytes.Buffer
	err = jpeg.Encode(&buffer, img, nil)
	if err != nil {
		global.Lg.Error(err.Error())
		return "图片加载失败", "", err
	}
	imgread := bytes.NewReader(buffer.Bytes())
	//获取封面的url
	cover_url, err := uploadAndGetUrl("cover", snapShotName, imgread, videoFrom.Data)
	if err != nil {
		global.Lg.Error(err.Error())
		return "上传失败", "", err
	}
	//进行Dao操作
	videoDB := dao.NewVideoDB(this.ctx)
	//生成存储对象
	createTime := time.Now()
	video := &models.Video{AuthorId: userId, PlayUrl: video_url, CoverUrl: cover_url, FavoriteCount: 0, CommentCount: 0, Title: videoFrom.Title, CreateTime: &createTime}
	_, err = videoDB.CreateVideoDao(video)

	if err != nil {
		return "db 失败", "", err
	}
	err = videoDB.IncreaseUserVideoInfoWorkCount()
	if err != nil {
		return "db 失败", "", err
	}
	return "成功", "", nil
}

func (this *VideoService) PubishList(form forms.VideoListForm) (interface{}, interface{}, error) {
	videoDB := dao.NewVideoDB(this.ctx)
	userDB := dao.NewUserDB(this.ctx)
	userId, _ := strconv.Atoi(form.UserId)
	videoList, err := videoDB.GetVideoList(userId)
	if err != nil {
		return "video db失败", "", err
	}
	userVideoInfo, err := videoDB.GetUserVideoInfo(userId)
	if err != nil {
		return "video db失败", "", err
	}
	info, err := userDB.GetOneUserInfo(userId)
	if err != nil {
		return "video db失败", "", err
	}
	followerId, _ := this.ctx.Get("userId")
	Followed, err := userDB.IsFollowed(userId, followerId.(int))
	if err != nil {
		return "video db失败", "", err
	}
	favorited := strconv.Itoa(userVideoInfo.FavoritedCount)
	author := forms.Author{Id: userId,
		Name:            info.UserName,
		FollowCount:     info.FollowCount,
		FollowerCount:   info.FollowerCount,
		IsFollow:        Followed,
		Avatar:          info.Avater,
		BackgroundImage: info.BackgroundImage,
		Signature:       info.Signature,
		TotalFavorited:  favorited,
		FavoriteCount:   userVideoInfo.FavoriteCount,
		WorkCount:       userVideoInfo.WorkCount,
	}
	resList := make([]forms.PublishRes, 0, len(videoList))
	for i, v := range videoList {
		resList[i].Author = author
		resList[i].VideoId = v.VideoId
		resList[i].CommentCount = v.CommentCount
		resList[i].FavoriteCount = v.FavoriteCount
		resList[i].Title = v.Title
		resList[i].CoverUrl = v.CoverUrl
		resList[i].PlayUrl = v.PlayUrl
	}
	return "查询成功", resList, nil
}

func (this *VideoService) FavoritedAction(form forms.VideoFavcriteForm) (interface{}, interface{}, error) {
	actionType, _ := strconv.Atoi(form.ActionType)
	videoDB := dao.NewVideoDB(this.ctx)
	userId, _ := this.ctx.Get("userId")

	videoId, _ := strconv.Atoi(form.VideoId)
	favorite := models.Favorite{
		VideoId: videoId,
		UserId:  userId.(int),
	}
	var err error
	var ok bool
	if actionType == 1 {
		ok, err = videoDB.CreateFavorite(&favorite)
	} else if actionType == 2 {
		ok, err = videoDB.DeleteFavorite(&favorite)
	} else {
		return "错误行动类型", "", errors.New("action tpye error")
	}
	if err != nil || !ok {
		return "action fail", "", err
	}
	return "ok", "", nil
}

func (this *VideoService) FavoriteListFormList(form forms.VideoFavoriteListForm) (interface{}, interface{}, error) {
	videoDB := dao.NewVideoDB(this.ctx)
	userDB := dao.NewUserDB(this.ctx)
	userId, _ := strconv.Atoi(form.UserId)
	videoList, err := videoDB.GetFavoriteList(userId)
	if err != nil {
		return "video db失败", "", err
	}
	userVideoInfo, err := videoDB.GetUserVideoInfo(userId)
	if err != nil {
		return "video db失败", "", err
	}
	info, err := userDB.GetOneUserInfo(userId)
	if err != nil {
		return "video db失败", "", err
	}
	followerId, _ := this.ctx.Get("userId")
	Followed, err := userDB.IsFollowed(userId, followerId.(int))
	if err != nil {
		return "video db失败", "", err
	}
	favorited := strconv.Itoa(userVideoInfo.FavoritedCount)
	author := forms.Author{Id: userId,
		Name:            info.UserName,
		FollowCount:     info.FollowCount,
		FollowerCount:   info.FollowerCount,
		IsFollow:        Followed,
		Avatar:          info.Avater,
		BackgroundImage: info.BackgroundImage,
		Signature:       info.Signature,
		TotalFavorited:  favorited,
		FavoriteCount:   userVideoInfo.FavoriteCount,
		WorkCount:       userVideoInfo.WorkCount,
	}
	resList := make([]forms.PublishRes, 0, len(videoList))
	for i, v := range videoList {
		resList[i].Author = author
		resList[i].VideoId = v.VideoId
		resList[i].CommentCount = v.CommentCount
		resList[i].FavoriteCount = v.FavoriteCount
		resList[i].Title = v.Title
		resList[i].CoverUrl = v.CoverUrl
		resList[i].PlayUrl = v.PlayUrl
	}
	return "查询成功", resList, nil

}

func (this VideoService) FeedList(form forms.FeedForm) (interface{}, interface{}, error) {

}

func uploadAndGetUrl(bucketName string, fileName string, fileobj io.Reader, header *multipart.FileHeader) (string, error) {
	// 把文件上传到minio对应的桶中
	ok := utils.UploadFile(bucketName, fileName, fileobj, header.Size)
	if !ok {
		err := errors.New("upload Fail")
		global.Lg.Error(err.Error())
		return "图像上传失败", err
	}
	headerUrl := utils.GetFileUrl(bucketName, fileName, time.Second*24*60*60)
	if headerUrl == "" {
		err := errors.New("getFileUrl fail")
		global.Lg.Error(err.Error())
		return "获取用户头像失败", err
	}
	return headerUrl, nil
}
