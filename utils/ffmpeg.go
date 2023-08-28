/*
* @Author: zgy
* @Date:   2023/8/14 16:10
 */
package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"go_gin/global"
	"image"
	"image/jpeg"
	"io"
	"os"
)

func exampleReadFrameAsJpeg(inFileName string, frameNum int) (io.Reader, error) {
	buffer := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buffer, os.Stdout).
		Run()
	if err != nil {
		global.Lg.Error(errors.New("获取封面失败").Error())
		return nil, err
	}

	img, err := imaging.Decode(buffer)
	if err != nil {
		global.Lg.Error(errors.New("获取封面失败").Error())
		return nil, err
	}

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)

	return buf, err
}

// 保存截图
func GetSnapShot(snapShotName string, videoFilePath string) (image.Image, error) {
	reader, err := exampleReadFrameAsJpeg(videoFilePath, 96)
	if err != nil {
		global.Lg.Error(errors.New("保存截图失败").Error())
		return nil, err
	}

	img, err := imaging.Decode(reader)
	if err != nil {
		global.Lg.Error(errors.New("保存截图失败").Error())
		return nil, err
	}
	err = imaging.Save(img, "./public/"+snapShotName)
	if err != nil {
		global.Lg.Error(errors.New("保存截图失败").Error())
		return nil, err
	}
	return img, err
}
