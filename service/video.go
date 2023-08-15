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
