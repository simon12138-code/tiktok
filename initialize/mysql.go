/*
* @Author: pzqu
* @Date:   2023/7/25 22:09
 */
package initialize

import (
	"fmt"
	"go_gin/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysqlDB() {
	mysqlInfo := global.Settings.Mysqlinfo
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name dsn本质就是访问数据库的对应完整的带设置的连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlInfo.Name, mysqlInfo.Password, mysqlInfo.Host,
		mysqlInfo.Port, mysqlInfo.DBName)
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	global.DB = db
}
