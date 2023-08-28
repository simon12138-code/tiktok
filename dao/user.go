/*
* @Author: zgy
* @Date:   2023/7/25 22:15
 */
package dao

//var users []models.User
//
//// GetUserList 获取用户列表(page第几页,page_size每页几条数据)
//func GetUserListDao(page int, page_size int) (int, []interface{}) {
//	//分页用户列表数据
//	userList := make([]interface{}, 0, len(users))
//	//计算偏移量
//	offset := (page - 1) * page_size
//	//查询所有user，gorm将查询的内容放入user变量中
//	result := global.DB.Offset(offset).Limit(page_size).Find(&users)
//	//查不到数据
//	if result.RowsAffected == 0 {
//		return 0, userList
//	}
//	//获取user总数
//	total := len(users)
//
//	//将user变量内的部分传入返回列表中
//	for _, useSingle := range users {
//		birthday := ""
//		if useSingle.Birthday == nil {
//			birthday = ""
//		} else {
//			// 给未设置生日的初始值
//			birthday = useSingle.Birthday.Format("2006-01-02")
//		}
//		userItemMap := map[string]interface{}{
//			"id":        useSingle.ID,
//			"password":  useSingle.Password,
//			"nick_name": useSingle.NickName,
//			"head_url":  useSingle.HeadUrl,
//			"birthday":  birthday,
//			"address":   useSingle.Address,
//			"desc":      useSingle.Desc,
//			"gender":    useSingle.Gender,
//			"role":      useSingle.Role,
//			"mobile":    useSingle.Mobile,
//		}
//		userList = append(userList, userItemMap)
//	}
//	return total, userList
//}
//
//var user models.User
//
//// 根据密码和用户名查找user表所有信息
//func UsernameFindUserInfo(username string, password string) (*models.User, bool) {
//	//调用DB
//	rows := global.DB.Where("nick_name = ? AND password = ?", username, password).Find(&user)
//	fmt.Println(&user)
//	//查询失败
//	if rows.RowsAffected < 1 {
//		return &user, false
//	}
//	return &user, true
//}
//
//func UserCreate(user *models.User) (*models.User, bool) {
//	//调用DB
//	rows := global.DB.Create(user)
//	fmt.Println(user)
//	//查询失败
//	if rows.RowsAffected < 1 {
//		return user, false
//	}
//	return user, true
//}
