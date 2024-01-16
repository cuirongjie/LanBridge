/*
缓存
*/
package cache

import (
	"net"
)

// 主连接
var MainConn net.Conn

// 配置
type Mapping struct {
	LocalPort    int    `json:"LocalPort"`
	DistCode     string `json:"RemoteCode"`
	DistPassword string `json:"RemotePassword"`
	DistAddr     string `json:"DistAddr"`
}
type Config struct {
	ServerAddr     string    `json:"ServerAddr"`
	ServerPassword string    `json:"ServerPassword"`
	MyCode         string    `json:"MyCode"`
	OpenMyPassword bool      `json:"OpenMyPassword"`
	MyPassword     string    `json:"MyPassword"`
	Debug          bool      `json:"Debug"`
	Mappings       []Mapping `json:"Mappings"`
}

var Conf Config
