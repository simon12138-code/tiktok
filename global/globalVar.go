/*
* @Author: zgy
* @Date:   2023/7/25 15:18
 */
package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
	"github.com/minio/minio-go"
	"go.uber.org/zap"
	"go_gin/config"
	"gorm.io/gorm"
)

var (
	Settings    config.ServerConfig
	Lg          *zap.Logger
	Trans       ut.Translator
	DB          *gorm.DB
	Redis       *redis.Client
	MinioClient *minio.Client
)
