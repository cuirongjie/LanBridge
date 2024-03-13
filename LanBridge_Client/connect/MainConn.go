/*
维护主连接
*/
package connect

import (
	"LanBridge_Client/cache"
	"LanBridge_Client/constant"
	"LanBridge_Client/logger"
	"errors"
	"net"
	"time"
)

// 连接Server
func StartMain() {
	defer func() {
		if cache.MainConn != nil {
			_ = cache.MainConn.Close()
		}
	}()
	logger.Info("正在连接服务器...")
	var err error
	for {
		time.Sleep(time.Second * 1)
		cache.MainConn, err = net.Dial("tcp", cache.Conf.ServerAddr)
		if err != nil {
			continue
		}
		connectRequest()
		logger.Debug("正在重连...")
	}
}

// 处理新建立的连接
func connectRequest() {
	defer func() {
		if cache.MainConn != nil {
			_ = cache.MainConn.Close()
		}
		logger.Info("与服务端连接断开")
	}()

	// 发送握手信息
	message := NewMessage(constant.Cmd_MainConn)
	message.SrcCode = cache.Conf.MyCode
	message.ServerPassword = cache.Conf.ServerPassword
	SendMessage(&cache.MainConn, message)

	// 获取服务器的应答
	for {
		var buf = make([]byte, 1024)
		cache.MainConn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := cache.MainConn.Read(buf)
		if err != nil {
			return
		}
		data := string(buf[:n])
		err1 := onData(data)
		if err1 != nil {
			return
		}
	}
}

// 处理来自服务器的数据
func onData(data string) (err error) {
	message := Str2Msg(data)
	switch message.Cmd {
	case constant.Cmd_MainConnectExist: // 识别码被占用
		logger.Info("本机识别码", cache.Conf.MyCode, "已存在（被占用）！！！")
		return errors.New("err")
	case constant.Cmd_BadCode: // 识别码有误（不在服务器白名单）
		logger.Info("本机识别码不被认可！！！")
		return errors.New("err")
	case constant.Cmd_BadServerPassword: // 提供的服务器密码有误
		logger.Info("提供的服务器密码有误！！！")
		return errors.New("err")
	case constant.Cmd_MainConnected: // 连接成功
		logger.Info("与服务端连接成功")
	case constant.Cmd_NotMainConnected: // 指定识别码的机器未连接
		logger.Info("设备", message.DistCode, "不存在或不在线")
	case constant.Cmd_BadClientPassword: // 提供的密码有误
		onBadClientPassword(message)
	case constant.Cmd_BridgeConn_Apply: // 桥接连接申请
		onBridgeApply(message)
	case constant.Cmd_ReverseProxyConn_Apply: // 代理连接申请
		onReverseProxyApply(message)
	case constant.Cmd_Heartbeat: // 收到心跳
		onHeartbeat(message)
	}
	return nil
}

// 收到心跳
func onHeartbeat(message Message) {
	cache.AllStatus = message
	logger.Debug("onHeartbeat", message)
}
