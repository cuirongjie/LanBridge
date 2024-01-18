/*
缓存
*/
package cache

import (
	"sync"
)

// 主连接池
var MainConns = &sync.Map{}

// 桥接Up连接池
var BridgeUpConns = &sync.Map{}

// 桥接Down连接池
var BridgeDownConns = &sync.Map{}

// 桥接连接状态池
var BridgeFlags = &sync.Map{}

// 反向代理Up连接池
var ReverseProxyUpConns = &sync.Map{}

// 反向代理Down连接池
var ReverseProxyDownConns = &sync.Map{}

// 反向代理连接池状态
var ReverseProxyFlags = &sync.Map{}

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
