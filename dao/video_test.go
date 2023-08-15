/*
* @Author: pzqu
* @Date:   2023/8/14 13:35
 */
package dao

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-playground/assert/v2"
	"github.com/spf13/viper"
	"go_gin/config"
	"go_gin/global"
	"go_gin/models"
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

func InitConfig() {
	//实例化Viper
	v := viper.New()
	//文件路径设置
	v.SetConfigFile("../settings-dev.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//创建配置变量
	serverConfig := config.ServerConfig{}
	//初始化 将绑定文件的配置信息反序列化到变量当中，完成文件信息到配置变量的转换
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	//传递全局变量 调用关系 settings-dev.yaml -> viper -> serverConfig(局部变量) -> global.Settings(全局变量)
	global.Settings = serverConfig
	color.Blue("11111111", global.Settings.LogAddress)
}
func InitMysqlDB() {
	mysqlInfo := global.Settings.Mysqlinfo
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name dsn本质就是访问数据库的对应完整的带设置的连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlInfo.Name, mysqlInfo.Password, mysqlInfo.Host,
		mysqlInfo.Port, mysqlInfo.DBName)
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	global.DB = db
}
func TestCreateVideoDao(t *testing.T) {
	//启动配置
	InitConfig()
	//启动数据库
	InitMysqlDB()
	//测试Dao接口
	time := time.Now()
	video := &models.Video{AuthorId: 1, PlayUrl: "http://test", CoverUrl: "https://", FavoriteCount: 0, CommentCount: 0, Title: "test", CreateTime: &time}
	videoDB := NewVideoDB(context.Background())

	id, err := videoDB.CreateVideoDao(video)
	videoDB.ctx = context.WithValue(videoDB.ctx, "user_id", id)
	err = videoDB.IncreaseUserVideoInfoWorkCount()
	assert.Equal(t, err, nil)
}
