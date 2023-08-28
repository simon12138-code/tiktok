/*
* @Author: pzqu
* @Date:   2023/7/25 22:07
 */

package models

type User struct {
	Id              int    `json:"id" gorm:"primaryKey"`
	UserName        string `json:"user_name"`
	Avater          string `json:"avater"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	FollowerCount   int    `json:"follower_count"`
	FollowCount     int    `json:"follow_count"`
	Password        string `json:"password"`
}

func (User) TableName() string {
	return "user"
}

type Relation struct {
	Id         int `json:"id" gorm:"primaryKey"`
	FollowerId int `json:"follower_id"`
	UserId     int `json:"user_id"`
	FriendFlag int `json:"friend_flag"`
}

func (Relation) TableName() string {
	return "relation"
}
