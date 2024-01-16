/*
消息工具类
*/
package connect

import (
	"LanBridge_Server/logger"
	"encoding/json"
	"net"
)

// 消息定义
type Message struct {
	Cmd            string `json:"Cmd"`
	ServerPassword string `json:"ServerPassword"`
	SrcCode        string `json:"SrcCode"`
	TunnelId       string `json:"TunnelId"`
	DistCode       string `json:"DistCode"`
	DistPassword   string `json:"DistPassword"`
	DistAddr       string `json:"DistAddr"`
}

// 消息对象转json字符串
func NewMessage(cmd string) Message {
	var message Message
	message.Cmd = cmd
	return message
}

// 消息对象转json字符串
func Msg2Str(message Message) string {
	bytes, err := json.Marshal(message)
	if err != nil {
		logger.Debug(err)
		return ""
	}
	return string(bytes)
}

// json字符串转消息对象
func Str2Msg(str string) Message {
	var message Message
	err := json.Unmarshal([]byte(str), &message)
	if err != nil {
		logger.Debug(err)
	}
	return message
}

// 消息拷贝
func CopyMessage(srcMessage Message) (distMessage Message) {
	strMessage := Msg2Str(srcMessage)
	distMessage = Str2Msg(strMessage)
	return
}

// 向指定连接发送数据
func SendMessage(conn net.Conn, message Message) (hasErr bool) {
	bytes, err := json.Marshal(message)
	if err != nil {
		logger.Debug(err)
		return true
	}
	_, err = conn.Write(bytes)
	if err != nil {
		return true
	}
	return false
}
