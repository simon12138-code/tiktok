/*
* @Author: zgy
* @Date:   2023/7/26 14:33
 */
package initialize

import (
	"github.com/minio/minio-go"
	"go_gin/global"
	"go_gin/utils"
	"log"
)

func InitMinio() {
	minioInfo := global.Settings.MinioInfo
	// 初使化 minio client对象。false是关闭https证书校验
	minioClient, err := minio.New(minioInfo.Endpoint, minioInfo.AccessKeyId, minioInfo.SecretAccessKey, false)
	if err != nil {
		log.Fatalln(err)
	}
	//客户端注册到全局变量中
	global.MinioClient = minioClient
	//创建一个叫userheader的存储桶。
	utils.CreateMinoBuket("video")
	utils.CreateMinoBuket("cover")
}
