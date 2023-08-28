/*
* @Author: pzqu
* @Date:   2023/8/24 19:16
 */
package dao

import (
	"context"
	"go_gin/models"
	"testing"
)

func Init() {
	InitConfig()
	InitMysqlDB()
}
func TestUserDB_UserCreate(t *testing.T) {
	Init()
	userdb := NewUserDB(context.Background())
	user := models.User{
		UserName:        "test",
		Avater:          "asdasd",
		BackgroundImage: "asdasd",
		Signature:       "asasd直接欧派那个v",
		FollowerCount:   0,
		FollowCount:     0,
		Password:        "sdadasdasd",
	}
	res, err := userdb.UserCreate(&user)
	if err != nil {
		panic(err)
	}
	println(res)
}
