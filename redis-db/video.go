/*
* @Author: pzqu
* @Date:   2023/8/20 17:26
 */
package redis_db

import (
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go_gin/forms"
	"go_gin/global"
	"go_gin/models"
	"go_gin/utils"
	"golang.org/x/net/context"
	"strconv"
	"sync"
	"time"
)

// redis存储数据结构
type Author struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	FollowCount     int    `json:"follow_count"`
	FollowerCount   int    `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  string `json:"total_favorited"`
	WorkCount       int    `json:"work_count"`
	FavoriteCount   int    `json:"favorite_count"`
}
type PublishRes struct {
	VideoId       int    `json:"video_id" `
	Author        Author `json:"author"`
	PlayUrl       string `json:"play_url" `
	CoverUrl      string `json:"cover_url" `
	FavoriteCount int    `json:"favorite_count" `
	CommentCount  int    `json:"comment_count" `
	Title         string `json:"title"`
}

// video_const_info
type VideoAuthor struct {
	AuthorId        int    `json:"author_id"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
}

// 固定数据采用sset的方式，以时间戳为score进行排序
type Video struct {
	Author   VideoAuthor `json:"author"`
	VideoId  int         `json:"video_id" `
	PlayUrl  string      `json:"play_url" `
	CoverUrl string      `json:"cover_url" `
	Title    string      `json:"title"`
}

// 计数器类型改用hset进行存储，key为对应的ID，
// video_var_info
type VideoCount struct {
	FavoriteCount int `json:"favorite_count" `
	CommentCount  int `json:"comment_count" `
}

// user_var_info
type UserVideoCount struct {
	TotalFavorited string `json:"total_favorited"`
	WorkCount      int    `json:"work_count"`
	FavoriteCount  int    `json:"favorite_count"`
}
type UserCount struct {
	FollowCount   int `json:"follow_count"`
	FollowerCount int `json:"follower_count"`
}

// 采用set进行存储key:userId Value:userId
type UserFollow map[string][]string

type VideoRedis struct {
	ctx context.Context
}

func NewVideoRdis(ctx context.Context) *VideoRedis {
	return &VideoRedis{ctx: ctx}
}

// 视频流缓存获取
func (this VideoRedis) GetFeed(form forms.FeedForm) ([]forms.FeedRes, error) {
	//创建请求参数
	res := make([]Video, 0, 30)
	var redisRangeArg redis.ZRangeArgs
	if form.LatestTime == "" {
		now := time.Now().Unix()
		nowStr := strconv.FormatInt(now, 10)
		redisRangeArg = redis.ZRangeArgs{
			Key:     global.RedisFeedKey,
			ByScore: true,
			Rev:     true,
			Start:   "0",
			Stop:    nowStr,
		}
	} else {
		redisRangeArg = redis.ZRangeArgs{
			Key:     global.RedisFeedKey,
			ByScore: true,
			Rev:     true,
			Start:   "-inf",
			Stop:    form.LatestTime,
		}
	}
	//1、获取对应的videoConstInfo
	resList, err := global.Redis.ZRangeArgs(context.Background(), redisRangeArg).Result()
	if err == redis.Nil {
		return nil, redis.Nil
	} else if err != nil {
		return nil, err
	}
	if len(resList) == 0 {
		return nil, redis.Nil
	}
	//json解析
	for i, j := len(resList)-1, 0; j < 30; {
		var temp Video
		_ = json.Unmarshal([]byte(resList[i]), &temp)
		i--
		res = append(res, temp)
		if i == 0 {
			break
		}

		j++
	}
	authorIds := make([]int, len(res))
	videoIds := make([]int, len(res))
	for i, v := range res {
		authorIds[i] = v.Author.AuthorId
		videoIds[i] = v.VideoId
	}

	//2、获取counter
	userCounter, err := getCounterList(this.ctx, "user", authorIds)
	if err != nil {
		return nil, err
	}
	videoCounter, err := getCounterList(this.ctx, "video", videoIds)
	if err != nil {
		return nil, err
	}
	//3、获取follower信息
	userId := this.ctx.Value("userId")

	isFollows, err := getFollowerList(this.ctx, authorIds, strconv.Itoa(userId.(int)))
	AllRes := make([]forms.FeedRes, len(res))

	//组装返回结果
	for i, v := range res {
		AllRes[i] = Cache2Res(v, userCounter[i], videoCounter[i], isFollows[i])
	}
	return AllRes, nil
}

func Cache2Res(v Video, userCounter map[string]interface{}, videoCounter map[string]interface{}, isFollow bool) forms.FeedRes {
	favoriteCount, _ := strconv.Atoi(videoCounter["favorite_count"].(string))
	commentCount, _ := strconv.Atoi(videoCounter["favorite_count"].(string))
	followCount, _ := strconv.Atoi(userCounter["follow_count"].(string))
	followerCount, _ := strconv.Atoi(userCounter["follower_count"].(string))
	totaFavorited, _ := userCounter["total_favorited"].(string)
	return forms.FeedRes{
		VideoId: v.VideoId,
		Author: forms.Author{
			Id:              v.Author.AuthorId,
			Name:            v.Author.Name,
			FollowCount:     followCount,
			FollowerCount:   followerCount,
			IsFollow:        isFollow,
			Avatar:          v.Author.Avatar,
			BackgroundImage: v.Author.BackgroundImage,
			Signature:       v.Author.Signature,
			TotalFavorited:  totaFavorited,
			WorkCount:       0,
			FavoriteCount:   0,
		},
		PlayUrl:       v.PlayUrl,
		CoverUrl:      v.CoverUrl,
		FavoriteCount: favoriteCount,
		CommentCount:  commentCount,
		Title:         v.Title,
	}
}

// 判定counter是否存在
func (this VideoRedis) CounterExists(CounterType string, Id int) bool {
	_, err := getCounter(this.ctx, CounterType, Id)
	if err != nil {
		return false
	}
	return true
}

// 判定relation是否存在
func (this VideoRedis) RelationExists(userId int) bool {
	return relationExists(this.ctx, userId)
}

// 插入feedList
func (this VideoRedis) InsertFeedList(feedForm []forms.FeedRes, score []int64, followList [][]int) error {
	//重新整理一下思路，插入对应的缓存信息为以下三个部分
	//1、sset，保存固有信息，score为时间戳
	//2、计数器类型，采用hset，分两个部分，一部分为以videoId为key的count，一部分为以userId为key的count
	//3、关注关系类型，采用set，以对应videoAuthorId为key进行存储，value为对应的followerID
	//所以该方法的步骤为
	//1、进行输入数据到缓存形式的pack，输入数据为res形式，缓存数据为定义的数据格式
	feedList, videoCountList, userCountList, userFollow := res2Cache(feedForm, score, followList)
	//2、开启四个协程完成内容插入
	var wg sync.WaitGroup
	wg.Add(4)
	err := make(chan error, 4)
	//协程处理
	go func() {
		defer wg.Done()
		_, inErr := global.Redis.ZAdd(this.ctx, global.RedisFeedKey, *feedList...).Result()
		if inErr != nil {
			err <- inErr
			return
		}
		err <- nil
	}()
	//插入用户计数器
	go func() {

		inErr := ListCounter(this.ctx, "user_id", *userCountList)
		if inErr != nil {
			err <- inErr
			return
		}
		err <- nil
	}()
	//插入video计数器
	go func() {

		inErr := ListCounter(this.ctx, "video_id", *videoCountList)
		if inErr != nil {
			err <- inErr
			return
		}
		err <- nil
	}()
	//插入用户关系
	go func() {

		inErr := addFollower(this.ctx, *userFollow)
		if inErr != nil {
			err <- inErr
			return
		}
		err <- nil

	}()
	//阻塞等待四个协程执行完毕
	//收集错误信息
	for i := 0; i < 4; i++ {
		stepErr := <-err
		if stepErr != nil {
			return stepErr
		}
	}
	return nil
}

func relationExists(ctx context.Context, userId int) bool {
	//空集合也算存在
	_, err := global.Redis.SMembers(ctx, strconv.Itoa(userId)).Result()
	if err != nil {
		return false
	}
	return true
}

// 添加关系setList
func addFollower(ctx context.Context, follow UserFollow) error {
	pipe := global.Redis.Pipeline()
	for k, v := range follow {
		uid, _ := strconv.Atoi(k)
		key := getFollowerKey(uid)
		_, err := pipe.Del(ctx, key).Result()
		if err != nil {
			return err
		}
		_, err = pipe.SAdd(ctx, key, v).Result()
		if err != nil {
			return err
		}
		global.Lg.Info(fmt.Sprintf("设置 uid=%d, key=%s\n", uid, key))
	}
	//批量执行
	_, err := pipe.Exec(ctx)
	if err != nil { // 报错后进行一次额外尝试
		_, err = pipe.Exec(ctx)
		if err != nil {
			global.Lg.Error(err.Error())
			return err
		}
	}
	return nil

}

func getFollowerKey(id int) string {
	return fmt.Sprintf("followerInfo_%d", id)
}

// 发布时的单条video缓存信息插入
func (this VideoRedis) InsertVideoAndVideoCounter(VideoForm forms.PublishRes, timeStamp time.Time) error {
	//需要查润video，userInfo,userVideoInfo,userFollower
	//前者直接插入，后三者判定是否存在，如果不存在则插入，如果存在则保持
	videoZ := redis.Z{}
	video, videoCounter := publish2video(VideoForm)
	byteList, _ := json.Marshal(*video)
	//调用高效转接器
	videoZ.Member = utils.Bytes2Str(byteList)
	videoZ.Score = float64(timeStamp.Unix())
	_, err := global.Redis.ZAdd(this.ctx, global.RedisFeedKey, videoZ).Result()
	if err != nil {
		return err
	}
	err = ListCounter(this.ctx, "video_id", []map[string]interface{}{*videoCounter})
	if err != nil {
		return err
	}
	return nil
}

// 发布视频时存入Counter
func (this VideoRedis) InsertUserCounter(info *forms.UserRes) error {

	userCount := []map[string]interface{}{{
		"user_id":         strconv.Itoa(info.Id),
		"follower_count":  strconv.Itoa(info.FollowerCount),
		"follow_count":    strconv.Itoa(info.FollowCount),
		"total_favorited": info.TotalFavorited,
		"work_count":      info.WorkCount,
	}}
	//videoCount := map[string]interface{}{
	//	"video_id":       strconv.Itoa(res.VideoId),
	//	"favorite_count": strconv.Itoa(res.FavoriteCount),
	//	"comment_count":  strconv.Itoa(res.CommentCount),
	//}
	err := ListCounter(this.ctx, "user_id", userCount)
	return err
}

// 存入用户的关系list
func (this VideoRedis) InsertUserRelation(relation []int, id int) error {
	follow := UserFollow{strconv.Itoa(id): List2StringList(relation)}
	err := addFollower(this.ctx, follow)
	return err
}

func (this VideoRedis) IncreaseFavorite(favorite models.Favorite, authorId int) {
	userId := favorite.UserId
	videoId := favorite.VideoId
	//判断videocount是否存在
	if ok, _ := global.Redis.Exists(this.ctx, GetVideoCounterKey(videoId)).Result(); ok == 1 {
		incrByUserLikeInVideoInfo(this.ctx, videoId)
	}
	//判断用户count是否存在
	if ok, _ := global.Redis.Exists(this.ctx, GetUserCounterKey(userId)).Result(); ok == 1 {
		incrByUserLikeInUserInfo(this.ctx, userId)
	}
	//判断作者count是否存在
	if ok, _ := global.Redis.Exists(this.ctx, GetUserCounterKey(authorId)).Result(); ok == 1 {
		incrByUserLikeInUserInfo(this.ctx, authorId)
	}
}

func (this VideoRedis) DecreaseFavorite(favorite models.Favorite, authorId int) {
	userId := favorite.UserId
	videoId := favorite.VideoId
	//判断videocount是否存在
	if ok, _ := global.Redis.Exists(this.ctx, GetVideoCounterKey(videoId)).Result(); ok == 1 {
		decrByUserLikeInVideoInfo(this.ctx, videoId)
	}
	//判断用户count是否存在
	if ok, _ := global.Redis.Exists(this.ctx, GetUserCounterKey(userId)).Result(); ok == 1 {
		decrByUserLikeInUserInfo(this.ctx, userId)
	}
	//判断作者count是否存在
	if ok, _ := global.Redis.Exists(this.ctx, GetUserCounterKey(authorId)).Result(); ok == 1 {
		decrByUserLikeInUserInfo(this.ctx, authorId)
	}
}

// 批量获取视频数据
func getCounterList(ctx context.Context, counterType string, IDS []int) ([]map[string]interface{}, error) {
	pipe := global.Redis.Pipeline()
	res := make([]map[string]interface{}, len(IDS))
	var key string
	for _, ID := range IDS {
		if counterType == "user" {
			key = GetUserCounterKey(ID)
		} else if counterType == "video" {
			key = GetVideoCounterKey(ID)
		}
		pipe.HGetAll(ctx, key)
	}
	//执行批量获取
	cmders, err := pipe.Exec(ctx)
	if err != nil {
		global.Lg.Error(err.Error())
		return nil, err
	}
	for i, cmder := range cmders {
		counterMap, err := cmder.(*redis.MapStringStringCmd).Result()
		if err != nil {
			global.Lg.Error(err.Error())
			return nil, err
		}
		temp := map[string]interface{}{}
		for field, value := range counterMap {
			temp[field] = value
		}
		res[i] = temp
	}
	return res, nil
}

// 批量获取用户是否关注
func getFollowerList(ctx context.Context, IDS []int, userId string) ([]bool, error) {
	res := make([]bool, len(IDS))
	pipe := global.Redis.Pipeline()
	for _, id := range IDS {
		pipe.SIsMember(ctx, strconv.Itoa(id), userId)
	}
	//批量执行
	cmders, err := pipe.Exec(ctx)
	//访问返回值
	for i, cmder := range cmders {
		res[i], err = cmder.(*redis.BoolCmd).Result()
		if err != nil {
			global.Lg.Error(err.Error())
			return nil, err
		}
	}
	return res, nil

}

// ListCounter用于重建计数器，FeedListcacheMiss CounterType 对应计数器的主键名称
func ListCounter(ctx context.Context, CounterType string, Counter []map[string]interface{}) error {
	pipe := global.Redis.Pipeline()
	for _, counter := range Counter {
		//获取Id
		id, err := strconv.Atoi(counter[CounterType].(string))
		var key string
		if CounterType == "user_id" {
			key = GetUserCounterKey(id)
		} else if CounterType == "video_id" {
			key = GetVideoCounterKey(id)
		}
		_, err = pipe.Del(ctx, key).Result()
		if err != nil {
			return err
		}
		_, err = pipe.HMSet(ctx, key, counter).Result()
		if err != nil {
			return err
		}
		global.Lg.Info(fmt.Sprintf("设置 type=%s,id=%d, key=%s\n", CounterType, id, key))
	}
	// 批量执行上面for循环设置好的hmset命令
	_, err := pipe.Exec(ctx)
	if err != nil { // 报错后进行一次额外尝试
		_, err = pipe.Exec(ctx)
		if err != nil {
			global.Lg.Error(err.Error())
			return err
		}
	}
	return nil
}

func GetUserCounterKey(userID int) string {
	return fmt.Sprintf("%s_%d", global.RedisUserCountKey, userID)
}

func getCounter(ctx context.Context, counterType string, ID int) (map[string]interface{}, error) {
	pipe := global.Redis.Pipeline()
	res := make(map[string]interface{})
	var key string
	if counterType == "user" {
		key = GetUserCounterKey(ID)
	} else if counterType == "video" {
		key = GetVideoCounterKey(ID)
	}
	pipe.HGetAll(ctx, key)
	cmders, err := pipe.Exec(ctx)
	if err != nil {
		global.Lg.Error(err.Error())
		return nil, err
	}
	for _, cmder := range cmders {
		counterMap, err := cmder.(*redis.MapStringStringCmd).Result()
		if err != nil {
			global.Lg.Error(err.Error())
			return nil, err
		}
		for field, value := range counterMap {
			res[field] = value
		}
	}
	return res, nil
}

// IncrByUserLike 点赞数+1
func incrByUserLikeInUserInfo(ctx context.Context, userID int) {
	incrByUserField(ctx, userID, "total_favorited")
}
func incrByUserLikeInVideoInfo(ctx context.Context, videoID int) {
	incrByVideoField(ctx, videoID, "favorite_count")
}

// DecrByUserLike 点赞数-1
func decrByUserLikeInUserInfo(ctx context.Context, userID int) {
	decrByUserField(ctx, userID, "total_favorited")
}
func decrByUserLikeInVideoInfo(ctx context.Context, videoID int) {
	decrByVideoField(ctx, videoID, "favorite_count")
}

// DecrByUserCollect 收藏数-1
//
//	func DecrByUserCollect(ctx context.Context, userID int) {
//		decrByUserField(ctx, userID, "follow_collect_set_count")
//	}
func incrByVideoField(ctx context.Context, videoID int, field string) {
	change(ctx, "video", videoID, field, 1)
}
func incrByUserField(ctx context.Context, userID int, field string) {
	change(ctx, "user", userID, field, 1)
}

func decrByUserField(ctx context.Context, userID int, field string) {
	change(ctx, "user", userID, field, 1)
}
func decrByVideoField(ctx context.Context, videoID int, field string) {
	change(ctx, "video", videoID, field, -1)
}

func change(ctx context.Context, counterType string, ID int, field string, incr int64) {
	var redisKey string
	if counterType == "user" {
		redisKey = GetUserCounterKey(ID)
	} else if counterType == "video" {
		redisKey = GetVideoCounterKey(ID)
	}
	before, err := global.Redis.HGet(ctx, redisKey, field).Result()
	if err != nil {
		panic(err)
	}
	beforeInt, err := strconv.ParseInt(before, 10, 64)
	if err != nil {
		panic(err)
	}
	if beforeInt+incr < 0 {
		global.Lg.Info(fmt.Sprintf("禁止变更计数，计数变更后小于0. %d + (%d) = %d\n", beforeInt, incr, beforeInt+incr))
		return
	}
	global.Lg.Info(fmt.Sprintf("user_id: %d\n更新前\n%s = %s\n--------\n", ID, field, before))
	_, err = global.Redis.HIncrBy(ctx, redisKey, field, incr).Result()
	if err != nil {
		panic(err)
	}
	// fmt.Printf("更新记录[%d]:%d\n", userID, num)
	count, err := global.Redis.HGet(ctx, redisKey, field).Result()
	if err != nil {
		panic(err)
	}
	global.Lg.Info(fmt.Sprintf("user_id: %d\n更新后\n%s = %s\n--------\n", ID, field, count))
}

func res2Cache(feedForm []forms.FeedRes, score []int64, followList [][]int) (*[]redis.Z, *[]map[string]interface{}, *[]map[string]interface{}, *UserFollow) {
	feedList := make([]redis.Z, len(feedForm))
	videoCountList := make([]map[string]interface{}, len(feedForm))
	userCountList := make([]map[string]interface{}, len(feedForm))
	userFollow := UserFollow{}
	//encode操作
	for i, v := range feedForm {
		video, userCount, videoCount := feedForm2video(v)
		byteList, _ := json.Marshal(video)
		//调用高效转接器
		feedList[i].Member = utils.Bytes2Str(byteList)
		feedList[i].Score = float64(score[i])
		videoCountList[i] = videoCount
		userCountList[i] = userCount
		stringList := List2StringList(followList[i])
		userFollow[strconv.Itoa(v.Author.Id)] = stringList
	}
	return &feedList, &videoCountList, &userCountList, &userFollow
}

func List2StringList(ints []int) []string {
	res := make([]string, len(ints))
	for i, v := range ints {
		res[i] = strconv.Itoa(v)
	}
	return res
}

func publish2video(res forms.PublishRes) (*Video, *map[string]interface{}) {
	video := &Video{
		Author: VideoAuthor{
			AuthorId:        res.Author.Id,
			Name:            res.Author.Name,
			Avatar:          res.Author.Avatar,
			BackgroundImage: res.Author.BackgroundImage,
			Signature:       res.Author.Signature,
		},
		Title:    res.Title,
		VideoId:  res.VideoId,
		CoverUrl: res.CoverUrl,
		PlayUrl:  res.PlayUrl,
	}
	videoCount := map[string]interface{}{
		"video_id":       strconv.Itoa(res.VideoId),
		"favorite_count": strconv.Itoa(res.FavoriteCount),
		"comment_count":  strconv.Itoa(res.CommentCount),
	}
	return video, &videoCount
}
func feedForm2video(res forms.FeedRes) (*Video, map[string]interface{}, map[string]interface{}) {
	video := &Video{
		Author: VideoAuthor{
			AuthorId:        res.Author.Id,
			Name:            res.Author.Name,
			Avatar:          res.Author.Avatar,
			BackgroundImage: res.Author.BackgroundImage,
			Signature:       res.Author.Signature,
		},
		Title:    res.Title,
		VideoId:  res.VideoId,
		CoverUrl: res.CoverUrl,
		PlayUrl:  res.PlayUrl,
	}
	userCount := map[string]interface{}{
		"user_id":         strconv.Itoa(res.Author.Id),
		"follower_count":  strconv.Itoa(res.Author.FollowerCount),
		"follow_count":    strconv.Itoa(res.Author.FollowCount),
		"total_favorited": res.Author.TotalFavorited,
		"work_count":      res.Author.WorkCount,
	}
	videoCount := map[string]interface{}{
		"video_id":       strconv.Itoa(res.VideoId),
		"favorite_count": strconv.Itoa(res.FavoriteCount),
		"comment_count":  strconv.Itoa(res.CommentCount),
	}
	return video, userCount, videoCount
}

func GetVideoCounterKey(videoId int) string {
	return fmt.Sprintf("%s_%d", global.RedisVideoCountKey, videoId)
}
