/*
* @Author: zgy
* @Date:   2023/7/25 17:20
 */

package forms

type UserLoginForm struct {
	Password string `form:"password" json:"password" binding:"required,max=32"`
	Username string `json:"username" form:"username" binding:"required,max=32"`
}
type UserRes struct {
	Id              int    `json:"id" gorm:"primaryKey"`
	UserName        string `json:"user_name"`
	Avater          string `json:"avater"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	FollowerCount   int    `json:"follower_count"`
	FollowCount     int    `json:"follow_count"`
	TotalFavorited  int    `json:"total_favorited"`
	WorkCount       int    `json:"work_count"`
	FavoriteCount   int    `json:"favorite_count"`
}
type GetUserInfoForm struct {
	UserId int    `json:"user_id" form:"user_id"  binding:"required" `
	Token  string `json:"token" form:"token" binding:"required"`
}
type FollowRes struct {
	Id              int    `json:"id" gorm:"primaryKey"`
	UserName        string `json:"user_name"`
	Avater          string `json:"avater"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	IsFollow        bool   `json:"is_follow"`
	FollowerCount   int    `json:"follower_count"`
	FollowCount     int    `json:"follow_count"`
	TotalFavorited  int    `json:"total_favorited"`
	WorkCount       int    `json:"work_count"`
	FavoriteCount   int    `json:"favorite_count"`
}
type FriendRes struct {
	Id              int    `json:"id" gorm:"primaryKey"`
	UserName        string `json:"name"`
	FollowCount     int    `json:"follow_count"`
	FollowerCount   int    `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  string `json:"total_favorited"`
	WorkCount       int    `json:"work_count"`
	FavoriteCount   int    `json:"favorite_count"`
}
type ActionForm struct {
	ToUserId   int `json:"to_user_id" form:"to_user_id"`
	ActionType int `json:"action_type" form:"action_type"` //1关注，2取消关注
}
