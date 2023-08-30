package cronjob

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"go_gin/dao"
	"go_gin/global"
	"go_gin/models"
	"go_gin/utils"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type CronJobList struct {
	jobList  map[string]Job
	StopChan <-chan struct{}
}

type Job struct {
	Ticker time.Ticker
	Func   JobFunc
}

type JobFunc func(ticker time.Ticker, stop <-chan struct{})

type CronJob interface {
	RigisterJob(string, Job) error
	Run()
	StopCronJob()
}

func NewCronJobList(stop chan struct{}) CronJobList {
	return CronJobList{jobList: map[string]Job{}, StopChan: stop}
}

// 注册定时任务
func (this CronJobList) RigisterJob(jobName string, job Job) error {
	if _, ok := this.jobList[jobName]; ok {
		return errors.New("already rigister this cronjob")
	}
	this.jobList[jobName] = job
	return nil
}

func (this CronJobList) Run() {
	for k, v := range this.jobList {
		global.Lg.Info(fmt.Sprintf("start job %s", k))
		v.Func(v.Ticker, this.StopChan)
	}
}

func (this CronJobList) StopCronJob() {
	<-this.StopChan
}

func HttpReset(ticker time.Ticker, stop <-chan struct{}) {
	for {
		select {
		case <-ticker.C:
			// 等待定时器触发
			//按照global的设置进行重新获取url
			//按照本地缓存读取文件信息和id
			resVideo, err := Reloadfiles("../public/", "record_video.txt")
			resCover, err := Reloadfiles("../public/", "record_cover.txt")
			if err != nil {
				panic(err)
			}
			//获取新的videourl
			videoInfos := make([]models.Video, len(resVideo))
			for i, v := range resVideo {
				headerUrl := utils.GetFileUrl("video", v.VideoFileName, global.MaxUrlExpireTime)
				if headerUrl == "" {
					err := errors.New("getFileUrl fail")
					global.Lg.Error(err.Error())
				}
				videoInfos[i].VideoId = v.VideoId
				videoInfos[i].PlayUrl = headerUrl
			}
			//获取新的videourl
			for i, v := range resCover {
				headerUrl := utils.GetFileUrl("video", v.VideoFileName, global.MaxUrlExpireTime)
				if headerUrl == "" {
					err := errors.New("getFileUrl fail")
					global.Lg.Error(err.Error())
				}
				videoInfos[i].VideoId = v.VideoId
				videoInfos[i].CoverUrl = headerUrl
			}
			//db更新url
			videoDB := dao.NewVideoDB(context.Background())
			err = videoDB.UpdateUrlList(videoInfos)
			if err != nil {
				global.Lg.Error(err.Error())
			}
		case <-stop:
			break
		}
	}
}
func Reloadfiles(dirname string, fileName string) ([]global.VideoInfo, error) {
	res := []string{}
	//删除本地缓存视频和封面
	files, err := os.ReadDir(dirname)
	if err != nil {
		global.Lg.Error(err.Error())
	}
	for _, file := range files {
		if file.Name() != "record_cover.txt" && file.Name() != "record_cover.txt" {
			err := os.Remove(path.Join(dirname, file.Name()))
			if err != nil {
				return nil, err
			}
		}
	}
	//打开文件记录新增文件名 os.O_RDWR（读写模式）|os.O_APPEND（写操作追加）
	file, err := os.OpenFile(path.Join(dirname, fileName), os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//然后读取文件的每一行
	//读原来文件的内容，并且显示在终端
	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		res = append(res, str)
	}
	VideoInfo := Str2Struct(res)
	return VideoInfo, nil
}

func Str2Struct(strList []string) []global.VideoInfo {
	res := make([]global.VideoInfo, len(strList))
	for i, v := range strList {
		tmpList := strings.Split(v, "___")
		videoId, _ := strconv.Atoi(tmpList[0])
		res[i] = global.VideoInfo{VideoId: videoId, VideoFileName: strings.ReplaceAll(tmpList[1], "\n", "")}
	}
	return res
}
