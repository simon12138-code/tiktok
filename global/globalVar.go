/*
* @Author: zgy
* @Date:   2023/7/25 15:18
 */
package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"
	"time"

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
	VideoChan   chan VideoInfo
	CoverChan   chan VideoInfo
)

type VideoInfo struct {
	VideoId       int
	VideoFileName string
}

const RedisFeedKey = "redisFeedKey"
const MaxFeedCacheNum = 60
const RedisUserCountKey = "RedisUserCountKey"
const RedisVideoCountKey = "RedisVideoCountKey"
const DBMaxInitRelationSliceNum = 30
const VideoInfoCollectorMaxNum = 100
const CollectorRetryTime = 3
const CollectorRetryTimeDuration = time.Second * 60

// max 7days
const MaxUrlExpireTime = time.Second * 7 * 24 * 60 * 60
