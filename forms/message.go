package forms

type Message struct {
	Id         int    `json:"id"`
	ToUserId   int    `json:"to_user_id"`
	FromUserId int    `json:"from_user_id"`
	Content    string `json:"content"`
}
type MessageRes struct {
	Id int `json:"id"`
	//消息id
	ToUserId int `json:"to_user_id"`
	//消息接收者id
	FromUserId int `json:"from_user_id"`
	//消息发送者id
	Content string `json:"content"`
	//消息内容
	CreateTime int64 `json:"create_time"`
	//消息发送时间 yyyy-MM-dd HH:MM
}
