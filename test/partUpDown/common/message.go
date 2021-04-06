package common

type Message struct {
	Type string `json:"type"`  // 消息类型
	Data string `json:"data"`  // 消息内容
	AddTime string `sql:"add_time"`  // 演示 sql tag的用法
}