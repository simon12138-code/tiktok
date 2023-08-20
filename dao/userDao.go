package dao

import (
	"context"
	"errors"
	"fmt"
	"go_gin/forms"
	"go_gin/global"
	"go_gin/models"
	"log"
)

var user models.User
var userVideoInfo models.UserVideoInfo
var chatRelation models.ChatContentIndex
var userRelation models.Relation

type userDB struct {
	ctx context.Context
}

func NewUserDB(ctx context.Context) userDB {
	return userDB{ctx: ctx}
}

func (db userDB) UserLogin(username string, password string) (*models.User, error) {
	//调用DB
	rows := global.DB.Where("user_name =  ?  ", username).Find(&user)
	fmt.Println(&user)
	//查询失败
	if rows.RowsAffected < 1 {
		return &user, errors.New("db 错误")
	}
	return &user, nil
}
func (db userDB) UserCreate(user *models.User) (*models.User, error) {
	//调用DB
	user.FollowerCount = 0
	user.FollowCount = 0
	user.Signature = "点击添加介绍，让大家认识你..."
	user.Avater = "defaultAvater.jpg"
	user.BackgroundImage = "defaultBackGround.jpg"
	rows := global.DB.Create(user)

	userVideoInfo.UserId = user.Id
	userVideoInfo.WorkCount = 0
	userVideoInfo.FavoriteCount = 0
	userVideoInfo.FavoritedCount = 0
	rows1 := global.DB.Create(userVideoInfo)
	if rows1.RowsAffected < 1 {
		log.Println("用户视频关系表创建失败")
		return nil, errors.New("db 错误")
	}
	fmt.Println(user)
	//查询失败
	if rows.RowsAffected < 1 {
		log.Println("用户表创建失败")
		return user, errors.New("db 错误")
	}
	return user, nil
}
func (db userDB) GetOneUserInfo(userId int) (*forms.UserRes, error) {
	var userInfo forms.UserRes

	rows := global.DB.Where("id = ?", userId).Find(&user)
	fmt.Println(&user)
	if rows.RowsAffected < 1 {
		return &userInfo, errors.New("db 错误")
	}

	rows1 := global.DB.Where("user_id = ?", userId).Find(&userVideoInfo)
	fmt.Println(&userVideoInfo)
	if rows1.RowsAffected < 1 {
		return &userInfo, errors.New("db 错误")
	}

	userInfo.Id = userId
	userInfo.UserName = user.UserName
	userInfo.Signature = user.Signature
	userInfo.FollowCount = user.FollowCount
	userInfo.FollowerCount = user.FollowerCount
	userInfo.Avater = user.Avater
	userInfo.BackgroundImage = user.BackgroundImage
	userInfo.TotalFavorited = userVideoInfo.FavoritedCount
	userInfo.WorkCount = userVideoInfo.WorkCount
	userInfo.FavoriteCount = userVideoInfo.FavoriteCount

	return &userInfo, nil

}
func (db userDB) GetBranchUsers(userIds []int, userid int) (*[]forms.FollowRes, error) {
	var followInfos []forms.FollowRes

	var relation models.Relation

	nowUser := userid
	for i := 0; i < len(userIds); i++ {
		rows := global.DB.Where("id = ?", userIds[i]).Find(&user)
		fmt.Println(&user)
		if rows.RowsAffected < 1 {
			return &followInfos, errors.New("db 错误")
		}

		rows1 := global.DB.Where("user_id = ?", userIds[i]).Find(&userVideoInfo)
		fmt.Println(&userVideoInfo)
		if rows1.RowsAffected < 1 {
			return &followInfos, errors.New("db 错误")
		}
		var followFlag bool
		rows2 := global.DB.Where("user_id = ? and follower_id = ?", userIds[i], nowUser).Find(&relation)
		if rows2.RowsAffected < 1 {
			followFlag = false
		} else {
			followFlag = true
		}

		var followInfo forms.FollowRes

		followInfo.Id = userIds[i]
		followInfo.UserName = user.UserName
		followInfo.Signature = user.Signature
		followInfo.FollowCount = user.FollowCount
		followInfo.FollowerCount = user.FollowerCount
		followInfo.Avater = user.Avater
		followInfo.IsFollow = followFlag
		followInfo.BackgroundImage = user.BackgroundImage
		followInfo.TotalFavorited = userVideoInfo.FavoritedCount
		followInfo.WorkCount = userVideoInfo.WorkCount
		followInfo.FavoriteCount = userVideoInfo.FavoriteCount
		followInfos = append(followInfos, followInfo)
	}
	return &followInfos, nil
}
func (db userDB) GetFollowerIds(userId int) ([]int, error) {
	var relations []models.Relation

	var Ids []int
	rows := global.DB.Where("user_id = ?", userId).Find(&relations)
	if rows.RowsAffected < 1 {
		return Ids, errors.New("粉丝数为0")
	}
	for _, relation := range relations {
		Ids = append(Ids, relation.FollowerId)
	}
	return Ids, nil

}
func (db userDB) GetFollowedUserIds(userId int) ([]int, error) {
	var relations []models.Relation

	var Ids []int
	rows := global.DB.Where("follower_id = ?", userId).Find(&relations)
	if rows.RowsAffected < 1 {
		return Ids, errors.New("粉丝数为0")
	}
	for _, relation := range relations {
		Ids = append(Ids, relation.UserId)
	}
	return Ids, nil

}
func (db userDB) GetFriendList(userId int) (*[]forms.FriendRes, error) {
	var relation models.Relation

	var friendsInfos []forms.FriendRes

	var relations []models.Relation

	var chatContent models.ChatContent

	rows := global.DB.Where("user_id = ? and friend_flag = ?", userId, 1).Find(&relations)
	if rows.RowsAffected < 1 {
		return &friendsInfos, nil
	}

	var friendids []int
	for _, relation1 := range relations {
		friendids = append(friendids, relation1.FollowerId)
	}

	nowUser := userId
	for i := 0; i < len(friendids); i++ {
		rows := global.DB.Where("id = ?", friendids[i]).Find(&user)
		fmt.Println(&user)
		if rows.RowsAffected < 1 {
			return &friendsInfos, errors.New("db 错误")
		}

		rows1 := global.DB.Where("user_id = ?", friendids[i]).Find(&userVideoInfo)
		fmt.Println(&userVideoInfo)
		if rows1.RowsAffected < 1 {
			return &friendsInfos, errors.New("db 错误")
		}
		var followFlag bool
		rows2 := global.DB.Where("user_id = ? and follower_id = ?", friendids[i], nowUser).Find(&relation)
		if rows2.RowsAffected < 1 {
			followFlag = false
		} else {
			followFlag = true
		}

		var friendInfo forms.FriendRes

		friendInfo.Id = friendids[i]
		friendInfo.UserName = user.UserName
		friendInfo.Signature = user.Signature
		friendInfo.FollowCount = user.FollowCount
		friendInfo.FollowerCount = user.FollowerCount
		friendInfo.Avater = user.Avater
		friendInfo.IsFollow = followFlag
		friendInfo.BackgroundImage = user.BackgroundImage
		friendInfo.TotalFavorited = userVideoInfo.FavoritedCount
		friendInfo.WorkCount = userVideoInfo.WorkCount
		friendInfo.FavoriteCount = userVideoInfo.FavoriteCount

		raws := global.DB.Where("user_id = ? and to_user_id = ?", nowUser, friendids[i]).Find(&chatRelation)
		if raws.RowsAffected < 1 {
			friendInfo.Message = ""
		}
		raws1 := global.DB.Where("content_id = ?", chatRelation.ContentIndex).Last(&chatContent)
		if raws1.RowsAffected < 1 {
			friendInfo.Message = ""

		}
		friendInfo.Message = chatContent.Content

		friendsInfos = append(friendsInfos, friendInfo)
	}

	return &friendsInfos, nil

}

func (db userDB) UserActionFollow(toUserId int, userId int) (string, error) {
	var relation models.Relation

	// 思路，如果目标用户未关注，则先在relation表中添加关系，再检查目标用户是否是我的粉丝，如果是friendFlag变未1（双方），最后对user表操作，自己的关注人数加一，对方的粉丝数加一
	// 取关思路相反
	nowUserId := userId

	raws := global.DB.Where("user_id = ? and follower_id = ?", toUserId, nowUserId).Find(&relation)
	if raws.RowsAffected < 1 {

		//没有对应关系，则创建关系
		relation.UserId = toUserId
		relation.FollowerId = nowUserId
		relation.FriendFlag = 0
		//操作关注和被关注者的数据库，分别使其关注数加一和粉丝数加一
		//粉丝

		var user models.User
		raws3 := global.DB.Where("id = ?", nowUserId).Find(&user)
		if raws3.RowsAffected < 1 {
			return "未查询到用户", errors.New("sql错误")
		}
		user.FollowCount++
		raws3 = global.DB.Save(&user)
		if raws3.RowsAffected < 1 {
			return "保存信息失败", errors.New("sql错误")
		}
		//被关注者
		var user1 models.User
		raws4 := global.DB.Where("id = ?", toUserId).Find(&user1)
		if raws4.RowsAffected < 1 {
			return "未查询到用户", errors.New("sql错误")
		}
		user1.FollowerCount++
		raws4 = global.DB.Save(&user1)
		if raws4.RowsAffected < 1 {
			return "保存信息失败", errors.New("sql错误")
		}

		//检查被关注者是否是关注者的粉丝，如果是，增加对应的朋友关系，如果不是，则进行下一步
		var relation1 models.Relation
		raws2 := global.DB.Where("user_id = ? and follower_id = ?", nowUserId, toUserId).Find(&relation1)
		if raws2.RowsAffected >= 1 {
			relation1.FriendFlag = 1
			relation.FriendFlag = 1
			raws2 = global.DB.Save(&relation1)
			if raws2.RowsAffected < 1 {
				return "设置朋友失败", errors.New("sql错误")
			}
		}

		raws1 := global.DB.Create(&relation)
		if raws1.RowsAffected < 1 {
			return "关注失败", errors.New("关注失败")
		}
		return "关注成功", nil
	}

	return "用户已关注", nil
}

func (db userDB) UserActionCancel(toUserId int, userid int) (string, error) {
	var relation models.Relation

	nowUserId := userid
	raws := global.DB.Where("user_id = ? and follower_id = ?", toUserId, nowUserId).Find(&relation)
	if raws.RowsAffected < 1 {
		return "用户未关注，无法取关", errors.New("sql错误")
	}

	//如果有对应关系，则执行取关注操作
	//删除对应关系并删除对应的朋友关系
	if relation.FriendFlag == 1 {

		var relation1 models.Relation
		raws2 := global.DB.Where("user_id = ? and follower_id = ?", nowUserId, toUserId).Find(&relation1)
		if raws2.RowsAffected >= 1 {
			relation1.FriendFlag = 0
			raws2 = global.DB.Save(&relation1)
			if raws2.RowsAffected < 1 {
				return "朋友关系更新失败", errors.New("sql错误")
			}
		}

	}
	raws5 := global.DB.Delete(&relation)
	if raws5.RowsAffected < 1 {
		return "取关失败", errors.New("sql错误")
	}
	var user models.User

	raws3 := global.DB.Where("id = ?", nowUserId).Find(&user)
	if raws3.RowsAffected < 1 {
		return "未查询到用户", errors.New("sql错误")
	}
	user.FollowCount--
	raws3 = global.DB.Save(&user)
	if raws3.RowsAffected < 1 {
		return "保存信息失败", errors.New("sql错误")
	}
	//被关注者
	var user1 models.User
	raws4 := global.DB.Where("id = ?", toUserId).Find(&user1)
	if raws4.RowsAffected < 1 {
		return "未查询到用户", errors.New("sql错误")
	}

	user1.FollowerCount--
	raws4 = global.DB.Save(&user1)
	if raws4.RowsAffected < 1 {
		return "保存信息失败", errors.New("sql错误")
	}

	return "取关成功", nil
}

func (db userDB) IsFollowed(userId int, followerId int) (bool, error) {
	rows := global.DB.Model(userRelation).Where("follower_id = ? AND user_id = ?", followerId, userId)
	if rows.Error != nil {
		return false, errors.New("sql错误")
	}
	if rows.RowsAffected < 1 {
		return false, nil
	}
	return true, nil
}
