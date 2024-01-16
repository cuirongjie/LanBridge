/*
维护客户端的主连接
*/
package connect

import (
	"LanBridge_Server/cache"
	"LanBridge_Server/constant"
	"LanBridge_Server/logger"
	"net"
	"strconv"
	"time"
)

// 启动主服务
func StartServer() {
	Port := strconv.Itoa(cache.Conf.Port)
	server, err := net.Listen("tcp", "0.0.0.0:"+Port)
	if err != nil {
		logger.Info("端口", Port, "被占用！！！")
		return
	}
	logger.Info("启动客户端监听:", Port)
	for {
		client, err := server.Accept()
		if err == nil {
			go onMainConnect(client)
		}
	}
}

// 处理新建立的连接
func onMainConnect(client net.Conn) {
	var message Message
	//
	ch := make(chan bool, 1)
	go func() {
		// 读取握手信息
		buf := make([]byte, 2048)
		n, err := client.Read(buf)
		if err != nil {
			ch <- false
			return
		}
		data := string(buf[:n])
		logger.Debug(data)
		message = Str2Msg(data)
		// 如果握手失败
		if len(message.Cmd) == 0 {
			ch <- false
			return
		}
		if message.Cmd == constant.Cmd_MainConn { // 主连接
			srcCode := message.SrcCode
			serverPassword := message.ServerPassword
			// 校验密码
			if cache.Conf.OpenServerPassword && cache.Conf.ServerPassword != serverPassword {
				// 给client发送密码错误
				newMessage := NewMessage(constant.Cmd_BadServerPassword)
				SendMessage(client, newMessage)
				ch <- false
				return
			}
			// 校验白名单
			if cache.Conf.OpenWhitelist {
				rightCode := false
				for i := range cache.Conf.Whitelist {
					if cache.Conf.Whitelist[i] == srcCode {
						rightCode = true
					}
				}
				// 给client发送识别码有误
				if rightCode == false {
					newMessage := NewMessage(constant.Cmd_BadCode)
					SendMessage(client, newMessage)
					ch <- false
					return
				}
			}
			// 放入连接池
			if cache.MainConns[srcCode] != nil {
				// 给client发送连接已存在
				newMessage := NewMessage(constant.Cmd_MainConnectExist)
				SendMessage(client, newMessage)
				ch <- false
				return
			}
			cache.MainConns[srcCode] = client
			// 给client发送连接成功
			newMessage := NewMessage(constant.Cmd_MainConnected)
			SendMessage(client, newMessage)
			go onMainConnected(srcCode)
			logger.Debug("接收到客户端", srcCode, "连接")
			ch <- true
		} else if message.Cmd == constant.Cmd_BridgeConn_Up { // BridgeUpTunnel连接
			// 放入连接池
			cache.BridgeUpConns[message.TunnelId] = client
			cache.BridgeFlags[message.TunnelId] = make(chan bool, 1)
			// 开始处理UpTunnel
			go OnBridgeUpTunnelConnected(message)
			logger.Debug("UpTunnel建立", message.TunnelId)
			ch <- true
		} else if message.Cmd == constant.Cmd_BridgeConn_Down { // BridgeDownTunnel连接
			tunnelId := message.TunnelId
			// 放入连接池
			cache.BridgeDownConns[tunnelId] = client
			// 开始处理DownTunnel
			go OnBridgeDownTunnelConnected(tunnelId)
			logger.Debug("DownTunnel建立", tunnelId)
			ch <- true
		} else if message.Cmd == constant.Cmd_ReverseProxyConn { // ReverseProxy连接
			tunnelId := message.TunnelId
			// 放入连接池
			cache.ReverseProxyDownConns[tunnelId] = client
			// 开始处理DownTunnel
			go OnReverseProxyDownTunnelConnected(tunnelId)
			logger.Debug("ReverseProxyDownTunnel建立", tunnelId)
			ch <- true
		} else {
			ch <- false
		}
	}()

	select {
	case <-time.After(constant.WaitSecond):
		logger.Debug("超时退出", message.Cmd)
		client.Close()
		return
	case connected := <-ch:
		if connected == false {
			client.Close()
			return
		}
	}
}

// 处理客户端主连接
func onMainConnected(clientCode string) {
	defer func() {
		if cache.MainConns[clientCode] != nil {
			cache.MainConns[clientCode].Close()
			delete(cache.MainConns, clientCode)
			logger.Debug("与客户端", clientCode, "连接断开")
		}
	}() // 读取握手信息
	for {
		buf := make([]byte, 1024)
		n, err := cache.MainConns[clientCode].Read(buf)
		if err != nil {
			return
		}
		data1 := string(buf[:n])
		logger.Debug("收到客户端", clientCode, "数据", data1)
		// 处理接收的数据
		message := Str2Msg(data1)
		switch message.Cmd {
		case constant.Cmd_BadClientPassword: // 客户端密码有误
			if cache.BridgeFlags[message.TunnelId] != nil {
				cache.BridgeFlags[message.TunnelId] <- false
			}
			newMessage := CopyMessage(message)
			newMessage.Cmd = constant.Cmd_BadClientPassword
			SendMessage(cache.MainConns[message.SrcCode], newMessage)
		}
	}
}