/*
* @Author: pzqu
* @Date:   2023/7/25 22:07
 */
package models

type User struct {
}

func (User) TableName() string {
	return "user"
}
