/*
消息工具类
*/
package connect

import (
	"LanBridge_Client/logger"
	"encoding/json"
	"net"
	"sync"
)

// 消息定义
type Message struct {
	Cmd            string     `json:"a"`
	ServerPassword string     `json:"b"`
	SrcCode        string     `json:"c"`
	TunnelId       string     `json:"d"`
	DistCode       string     `json:"e"`
	DistPassword   string     `json:"f"`
	DistAddr       string     `json:"g"`
	ConnStatus     ConnStatus `json:"h"`
}

// 消息定义
type ConnStatus struct {
	MainCodes []string `json:"i"`
	Clients   []string `json:"g"`
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

// 消息对象转json字符串
func Obj2Str(jsonstr any) string {
	bytes, err := json.Marshal(jsonstr)
	if err != nil {
		logger.Debug(err)
		return ""
	}
	return string(bytes)
}

// 消息拷贝
func CopyMessage(srcMessage Message) (distMessage Message) {
	strMessage := Msg2Str(srcMessage)
	distMessage = Str2Msg(strMessage)
	return
}

// 向指定连接发送数据
func SendMessage(conn *net.Conn, message Message) (hasErr bool) {
	bytes, err := json.Marshal(message)
	if err != nil {
		logger.Debug(err)
		return true
	}
	_, err = (*conn).Write(bytes)
	if err != nil {
		return true
	}
	return false
}

// 获取指定连接
func GetConn(connMap *sync.Map, key string) *net.Conn {
	conn, ok := connMap.Load(key)
	if ok {
		return conn.(*net.Conn)
	} else {
		return nil
	}
}

// 保存指定连接
func StoreConn(connMap *sync.Map, key string, conn *net.Conn) {
	CloseAndDelConn(connMap, key)
	connMap.Store(key, conn)
}

// 删除指定连接
func CloseAndDelConn(connMap *sync.Map, key string) (ok bool) {
	conn, ok := connMap.LoadAndDelete(key)
	if ok {
		_ = (*conn.(*net.Conn)).Close()
		conn = nil
	}
	return ok
}

// 获取指定状态
func GetFlag(connMap *sync.Map, key string) *chan bool {
	ch, ok := connMap.Load(key)
	if ok {
		return ch.(*chan bool)
	} else {
		return nil
	}
}

// 保存指定状态
func StoreFlag(connMap *sync.Map, key string, flag *chan bool) {
	connMap.Store(key, flag)
}

// 删除指定状态
func CloseAndDelFlag(connMap *sync.Map, key string) {
	connMap.Delete(key)
}
