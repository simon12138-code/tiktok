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
	"github.com/redis/go-redis/v9"
	"go_gin/dao"
	"go_gin/forms"
	"go_gin/global"
	"go_gin/models"
	redis_db "go_gin/redis-db"
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
	userDB := dao.NewUserDB(this.ctx)
	//生成存储对象
	createTime := time.Now()
	video := &models.Video{AuthorId: userId, PlayUrl: video_url, CoverUrl: cover_url, FavoriteCount: 0, CommentCount: 0, Title: videoFrom.Title, CreateTime: &createTime}

	if err != nil {
		global.Lg.Error(err.Error())
		return "上传失败", "", err
	}
	//存入DB
	videoId, err := videoDB.CreateVideoDao(video)
	if err != nil {
		return "db 失败", "", err
	}
	err = videoDB.IncreaseUserVideoInfoWorkCount()
	if err != nil {
		return "db 失败", "", err
	}
	var redisCacheErrChan chan error
	//插入缓存
	go func() {
		//必须DB
		info, err := userDB.GetOneUserInfo(userId)
		if err != nil {
			redisCacheErrChan <- err
			return
		}
		videoRedis := redis_db.NewVideoRdis(this.ctx)
		//如果对应的用户缓存不存在则则更新（一个用户可能有或者没有发布视频，或者上一条视频还在缓存中）
		if videoRedis.CounterExists("user_id", userId) {

			//存入缓存
			err = videoRedis.InsertUserCounter(info)
			if err != nil {
				redisCacheErrChan <- err
				return
			}
		}
		//判断relation表是否存在
		if videoRedis.RelationExists(userId) {
			//获取用户关系
			userRelation, err := userDB.GetFollowedUserIds(userId)
			if err != nil {
				redisCacheErrChan <- err
				return
			}
			//存入缓存
			err = videoRedis.InsertUserRelation(userRelation, userId)
			if err != nil {
				redisCacheErrChan <- err
				return
			}
		}

		//插入必插入项
		author := forms.Author{Id: userId,
			Name:            info.UserName,
			FollowCount:     info.FollowCount,
			FollowerCount:   info.FollowerCount,
			IsFollow:        true,
			Avatar:          info.Avater,
			BackgroundImage: info.BackgroundImage,
			Signature:       info.Signature,
		}
		videoCache := forms.PublishRes{}
		videoCache.Author = author
		videoCache.VideoId = videoId
		videoCache.CommentCount = 0
		videoCache.FavoriteCount = 0
		videoCache.Title = video.Title
		videoCache.CoverUrl = video.CoverUrl
		videoCache.PlayUrl = video.PlayUrl
		err = videoRedis.InsertVideoAndVideoCounter(videoCache, createTime)
		if err != nil {
			redisCacheErrChan <- err
			return
		}
		redisCacheErrChan <- nil
	}()
	if err = <-redisCacheErrChan; err != nil {
		global.Lg.Error(err.Error())
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
		var authorId int
		authorId, ok, err = videoDB.CreateFavorite(&favorite)
		//插入缓存（如果存在的话）
		go func() {
			videoRedis := redis_db.NewVideoRdis(this.ctx)
			videoRedis.IncreaseFavorite(favorite, authorId)
		}()
	} else if actionType == 2 {
		var authorId int
		authorId, ok, err = videoDB.DeleteFavorite(&favorite)
		//插入缓存（如果存在的话）
		go func() {
			videoRedis := redis_db.NewVideoRdis(this.ctx)
			videoRedis.DecreaseFavorite(favorite, authorId)
		}()
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

func (this VideoService) FeedList(form forms.FeedForm) (interface{}, interface{}, int, error) {
	//预分配内存
	res := make([]forms.FeedRes, global.MaxFeedCacheNum)
	userId, _ := this.ctx.Get("userId")
	//先访问redis
	videoRedis := redis_db.NewVideoRdis(this.ctx)
	videoList, err := videoRedis.GetFeed(form)
	if err != redis.Nil && err != nil {
		return "redis err", "", 0, err
	} else if err == redis.Nil {
		//缓存失效，访问db
		videoDB := dao.NewVideoDB(this.ctx)
		//预分配用户内存
		userList := make([]int, global.MaxFeedCacheNum)
		videoList := make([]models.Video, global.MaxFeedCacheNum)
		if form.LatestTime == "" {
			//获取当前时间戳
			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
			videoList, err = videoDB.GetFeedVideoList(timestamp)
			if err != nil {
				return "db err", "", 0, err
			}
		} else {
			videoList, err = videoDB.GetFeedVideoList(form.LatestTime)
			if err != nil {
				return "db err", "", 0, err
			}
		}
		//获取喜爱列表的用户idlist
		for i, v := range videoList {
			userList[i] = v.AuthorId
		}
		//获取用户信息
		userDB := dao.NewUserDB(this.ctx)
		userInfoList, err := userDB.GetUserList(userList)
		if err != nil {
			return "db err", "", 0, err
		}
		//获取用户视频信息
		userVideoInfoList, err := videoDB.GetUserVideoInfoList(userList)
		if err != nil {
			return "db err", "", 0, err
		}
		//获取用户关注信息
		userFollowerList, err := userDB.GetUsersFollowerIds(userList)
		if err != nil {
			return "db err", "", 0, err
		}
		//装填返回值
		for i := 0; i < len(videoList); i++ {
			//userInfo
			res[i].Author.Id = userInfoList[i].Id
			res[i].Author.Name = userInfoList[i].UserName
			res[i].Author.Signature = userInfoList[i].Signature
			res[i].Author.BackgroundImage = userInfoList[i].BackgroundImage
			res[i].Author.Avatar = userInfoList[i].Avater
			res[i].Author.FollowCount = userInfoList[i].FollowCount
			res[i].Author.FollowerCount = userInfoList[i].FollowerCount
			//二分查找
			res[i].Author.IsFollow = isFollow(userFollowerList[i], userId.(int))
			//userVideoInfo
			res[i].Author.FavoriteCount = userVideoInfoList[i].FavoriteCount
			res[i].Author.TotalFavorited = strconv.Itoa(userVideoInfoList[i].FavoritedCount)
			res[i].Author.WorkCount = userVideoInfoList[i].WorkCount
			//VideoInfo
			res[i].VideoId = videoList[i].VideoId
			res[i].FavoriteCount = videoList[i].FavoriteCount
			res[i].CommentCount = videoList[i].CommentCount
			res[i].Title = videoList[i].Title
			res[i].PlayUrl = videoList[i].PlayUrl
			res[i].CoverUrl = videoList[i].CoverUrl
		}
		//创建时间戳序列,默认推流个数30
		timeList := make([]int64, global.MaxFeedCacheNum)
		for i := 0; i < global.MaxFeedCacheNum; i++ {
			timeList[i] = videoList[i].CreateTime.Unix()
		}
		//将查询到的内容插入到redis中更新
		err = videoRedis.InsertFeedList(res, timeList, userFollowerList)
		if err != nil {
			return "", "", 0, err
		}
		next_time := timeList[30]
		return "success", res[:30], int(next_time), nil
	}

	return "", videoList, int(time.Now().Unix()), nil
}

func isFollow(nums []int, target int) bool {
	low, high := 0, len(nums)-1
	mid := 0
	for low <= high {
		mid = low + (high-low)/2
		if nums[mid] == target {
			return true
		} else if nums[mid] > target {
			high = mid - 1
		} else if nums[mid] < target {
			low = mid + 1
		}
	}
	return false
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
