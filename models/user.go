/*
* @Author: pzqu
* @Date:   2023/7/25 22:07
 */
package models

import "time"

type User struct {
	ID       uint       `json:"id" gorm:"primaryKey"`
	Password string     `json:"password" `
	NickName string     `json:"nick_name"`
	HeadUrl  string     `json:"head_url"`
	Birthday *time.Time `json:"birthday" gorm:"type:date"`
	Address  string     `json:"address"`
	Desc     string     `json:"desc"`
	Gender   string     `json:"gender"`
	Role     int        `json:"role"`
	Mobile   string     `json:"mobile"`
}

func (User) TableName() string {
	return "user"
}
