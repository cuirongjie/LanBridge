/*
入口
*/
package main

import (
	"LanBridge_Server/cache"
	"LanBridge_Server/connect"
	"LanBridge_Server/logger"
	"encoding/json"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

// 启动主服务
func main() {
	ch := make(chan bool)
	// 初始化
	cache.MainConns = make(map[string]net.Conn)
	cache.BridgeUpConns = make(map[string]net.Conn)
	// 读取、校验配置
	bytes, err1 := ioutil.ReadFile("server.config")
	var conf = cache.Config{
		Port: 28010,
	}
	err2 := json.Unmarshal(bytes, &conf)
	if err1 != nil || err2 != nil {
		logger.Info("采用默认配置")
	} else {
		conf.ServerPassword = strings.TrimSpace(conf.ServerPassword)
		if conf.Port < 1 || conf.Port > 65535 {
			logger.Info("请正确配置端口")
			time.Sleep(2000)
			return
		}
		if conf.Mappings != nil {
			for i := range conf.Mappings {
				mapping := conf.Mappings[i]
				if mapping.ServerPort < 1 || mapping.ServerPort > 65535 {
					logger.Info("请正确配置端口")
					time.Sleep(2000)
					return
				}
			}
		}
		if len(conf.ServerPassword) == 0 {
			logger.Info("连接密码未启用")
		} else if len(conf.ServerPassword) > 20 {
			logger.Info("连接密码有误，长度不能超过 20 个字符")
			time.Sleep(2000)
			return
		} else {
			conf.OpenServerPassword = true
			logger.Info("连接密码已启用")
		}
		if len(conf.Whitelist) > 0 {
			conf.OpenWhitelist = true
			logger.Info("白名单已启用")
		} else {
			logger.Info("白名单未启用")
		}
	}
	cache.Conf = conf

	// 启动主连接
	go connect.StartServer()

	// 启动反向代理服务连接
	go connect.StartReverseProxyServer()

	<-ch
}
