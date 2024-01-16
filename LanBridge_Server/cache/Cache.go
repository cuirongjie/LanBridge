/*
缓存
*/
package cache

import (
	"net"
)

// 主连接池
var MainConns map[string]net.Conn

// 桥接Up连接池
var BridgeUpConns = make(map[string]net.Conn)

// 桥接Down连接池
var BridgeDownConns = make(map[string]net.Conn)

// 桥接连接状态池
var BridgeFlags = make(map[string]chan bool)

// 反向代理Up连接池
var ReverseProxyUpConns = make(map[string]net.Conn)

// 反向代理Down连接池
var ReverseProxyDownConns = make(map[string]net.Conn)

// 反向代理连接池状态
var ReverseProxyFlags = make(map[string]chan bool)

// 配置
type Mapping struct {
	ServerPort int    `json:"ServerPort"`
	RemoteCode string `json:"RemoteCode"`
	DistAddr   string `json:"DistAddr"`
}
type Config struct {
	Port               int       `json:"Port"`
	OpenServerPassword bool      `json:"OpenServerPassword"`
	ServerPassword     string    `json:"ServerPassword"`
	Debug              bool      `json:"Debug"`
	OpenWhitelist      bool      `json:"OpenWhitelist"`
	Whitelist          []string  `json:"Whitelist"`
	Mappings           []Mapping `json:"Mappings"`
}

var Conf Config
