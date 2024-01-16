/*
维护主连接
*/
package connect

import (
	"LanBridge_Client/cache"
	"LanBridge_Client/constant"
	"LanBridge_Client/logger"
	"net"
	"time"
)

var connected = false

// 退出程序
var reqExit = false
var ChExit = make(chan bool, 1)

// 连接Server
func StartMain() {
	defer func() {
		if cache.MainConn != nil {
			cache.MainConn.Close()
		}
		connected = false
	}()
	var err error
	for {
		time.Sleep(time.Millisecond * 500)
		if reqExit {
			ChExit <- true
			return
		}
		if connected {
			continue
		}
		cache.MainConn, err = net.Dial("tcp", cache.Conf.ServerAddr)
		if err != nil {
			logger.Debug("与服务端连接失败，正在重连...")
			continue
		}
		connected = true
		go connectRequest()
		time.Sleep(time.Second * 1)
	}
}

// 处理新建立的连接
func connectRequest() {
	defer func() {
		cache.MainConn.Close()
		connected = false
		logger.Info("与服务端连接断开")
	}()

	// 发送握手信息
	message := NewMessage(constant.Cmd_MainConn)
	message.SrcCode = cache.Conf.MyCode
	message.ServerPassword = cache.Conf.ServerPassword
	SendMessage(cache.MainConn, message)

	// 获取服务器的应答
	for {
		var buf = make([]byte, 1024)
		n, err := cache.MainConn.Read(buf)
		if err != nil {
			return
		}
		data := string(buf[:n])
		go onData(data)
	}
}

// 处理来自服务器的数据
func onData(data string) {
	message := Str2Msg(data)
	switch message.Cmd {
	case constant.Cmd_MainConnected: // 连接成功
		logger.Info("与服务端连接成功")
	case constant.Cmd_MainConnectExist: // 识别码被占用
		logger.Info("本机识别码", cache.Conf.MyCode, "已存在（被占用）！！！")
		reqExit = true
		ChExit <- true
	case constant.Cmd_BadCode: // 识别码有误（不在服务器白名单）
		logger.Info("本机识别码不被认可！！！")
		reqExit = true
		ChExit <- true
	case constant.Cmd_NotMainConnected: // 指定识别码的机器未连接
		logger.Info("设备", message.DistCode, "不存在或不在线")
	case constant.Cmd_BadServerPassword: // 指定识别码的机器未连接
		logger.Info("服务器密码有误！！！")
		reqExit = true
		ChExit <- true
	case constant.Cmd_BadClientPassword: // 提供的密码有误
		onBadClientPassword(message)
	case constant.Cmd_BridgeConn_Apply: // 桥接连接申请
		onBridgeApply(message)
	case constant.Cmd_ReverseProxyConn_Apply: // 桥接连接申请
		onReverseProxyApply(message)
	}
}
