/*
* @Author: zgy
* @Date:   2023/7/26 14:41
 */
package utils

import (
	"fmt"
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/policy"
	"go.uber.org/zap"
	"go_gin/global"
	"io"
	"net/url"
	"time"
)

// 根据对象名创建一个对象存储bucket
func CreateMinoBuket(bucketName string) {
	location := "us-east-1"
	//创建对象存储桶，location == relative_path
	err := global.MinioClient.MakeBucket(bucketName, location)
	if err != nil {
		// 检查存储桶是否已经存在。
		exists, err := global.MinioClient.BucketExists(bucketName)
		fmt.Println(exists)
		if err == nil && exists {
			fmt.Printf("We already own %s\n", bucketName)
		} else {
			fmt.Println(err, exists)
			return
		}
	}
	//设置该桶的权限->读写都可
	err = global.MinioClient.SetBucketPolicy(bucketName, policy.BucketPolicyReadWrite)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Successfully created %s\n", bucketName)

}

// UploadFile 上传文件给minio指定的桶中
func UploadFile(bucketName, objectName string, reader io.Reader, objectSize int64) bool {
	//指定方式，桶名，对象名称，io流（存储对象本体），对象大小，对象类型
	n, err := global.MinioClient.PutObject(bucketName, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Successfully uploaded bytes: ", n)
	return true
}

// GetFileUrl 获取文件url
func GetFileUrl(bucketName string, fileName string, expires time.Duration) string {
	//time.Second*24*60*60
	reqParams := make(url.Values)
	//expires 是 获取头像url的过期时间
	presignedURL, err := global.MinioClient.PresignedGetObject(bucketName, fileName, expires, reqParams)
	if err != nil {
		zap.L().Error(err.Error())
		return ""
	}
	return fmt.Sprintf("%s", presignedURL)
}
