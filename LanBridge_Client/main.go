/*
入口
*/
package main

import (
	"LanBridge_Client/cache"
	"LanBridge_Client/connect"
	"LanBridge_Client/logger"
	"LanBridge_Client/webui"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func main() {
	fmt.Print("\n ___        ________   ________    ________   ________   ___   ________   ________   _______      \n|\\  \\      |\\   __  \\ |\\   ___  \\ |\\   __  \\ |\\   __  \\ |\\  \\ |\\   ___ \\ |\\   ____\\ |\\  ___ \\     \n\\ \\  \\     \\ \\  \\|\\  \\\\ \\  \\\\ \\  \\\\ \\  \\|\\ /_\\ \\  \\|\\  \\\\ \\  \\\\ \\  \\_|\\ \\\\ \\  \\___| \\ \\   __/|    \n \\ \\  \\     \\ \\   __  \\\\ \\  \\\\ \\  \\\\ \\   __  \\\\ \\   _  _\\\\ \\  \\\\ \\  \\ \\\\ \\\\ \\  \\  ___\\ \\  \\_|/__  \n  \\ \\  \\____ \\ \\  \\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\|\\  \\\\ \\  \\\\  \\|\\ \\  \\\\ \\  \\_\\\\ \\\\ \\  \\|\\  \\\\ \\  \\_|\\ \\ \n   \\ \\_______\\\\ \\__\\ \\__\\\\ \\__\\\\ \\__\\\\ \\_______\\\\ \\__\\\\ _\\ \\ \\__\\\\ \\_______\\\\ \\_______\\\\ \\_______\\\n    \\|_______| \\|__|\\|__| \\|__| \\|__| \\|_______| \\|__|\\|__| \\|__| \\|_______| \\|_______| \\|_______|\n" +
		"                                                            |_      _   . _ _  _  _  .. _ \n                                                            |_)\\/  (_|_||| (_)| |(_| ||(/_\n                                                               /                  _|L|    \n")
	// 读取、校验配置
	bytes, err1 := ioutil.ReadFile("client.config")
	var conf cache.Config
	err2 := json.Unmarshal(bytes, &conf)
	if err1 != nil || err2 != nil {
		logger.Info("配置有误", err1)
		time.Sleep(2000)
		return
	}
	conf.ServerPassword = strings.TrimSpace(conf.ServerPassword)
	conf.MyPassword = strings.TrimSpace(conf.MyPassword)
	if conf.ServerAddr == "" {
		logger.Info("ServerAddr配置有误")
		time.Sleep(2000)
		return
	}
	if len(conf.MyCode) < 5 || len(conf.MyCode) > 20 {
		logger.Info("MyCode配置有误，长度应在 5-20 个字符")
		time.Sleep(2000)
		return
	}
	//if len(conf.Mappings) == 0 {
	//	logger.Info("Mappings配置有误")
	//	time.Sleep(2000)
	//	return
	//}
	if len(conf.MyPassword) == 0 {
		logger.Info("本机连接密码未启用")
	} else {
		conf.OpenMyPassword = true
		logger.Info("本机连接密码已启用")
	}
	cache.Conf = conf

	// 启动主连接
	go connect.StartMain()

	// 启动监听服务
	go connect.StartBridgeServer()

	// 启动webui
	go webui.StartWebUI()

	var ChExit = make(chan bool, 1)
	<-ChExit
	logger.Info("程序即将退出...")
	time.Sleep(time.Second * 3)
}
