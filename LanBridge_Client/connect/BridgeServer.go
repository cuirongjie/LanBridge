/*
桥接连接服务端
按照配置文件，启动监听服务
*/
package connect

import (
	"LanBridge_Client/cache"
	"LanBridge_Client/constant"
	"LanBridge_Client/logger"
	"io"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// 按照配置文件，启动监听服务
func StartBridgeServer() {
	for _, mapping := range cache.Conf.Mappings {
		go startListen(mapping)
	}
}

// 启动监听
func startListen(mapping cache.Mapping) {
	LocalPort := strconv.Itoa(mapping.LocalPort)
	server, err := net.Listen("tcp", "0.0.0.0:"+LocalPort)
	if err != nil {
		logger.Info("端口", LocalPort, "被占用！！！")
		ChExit <- true
		return
	}
	logger.Debug("启动监听:", LocalPort)
	for {
		client, err := server.Accept()
		if err == nil {
			go onBridgeVisitConnected(client, mapping)
		}
	}
}

// 接收到连接
func onBridgeVisitConnected(clientConn net.Conn, mapping cache.Mapping) {
	// tunnelId
	rand.Intn(100000)
	tunnelId := strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(rand.Intn(100000))
	//
	defer func() {
		if cache.MainConn != nil {
			_ = clientConn.Close()
		}
	}()
	// 创建与服务器UpTunnel
	upTunnelConn, err := net.Dial("tcp", cache.Conf.ServerAddr)
	if err != nil {
		logger.Debug("UpTunnel连接服务端", cache.Conf.ServerAddr, "出错,错误原因: ", err)
		return
	}
	defer func() {
		if cache.MainConn != nil {
			_ = upTunnelConn.Close()
		}
		logger.Debug("UpTunnel 连接断开", tunnelId)
	}()
	// 发送握手信息
	message := NewMessage(constant.Cmd_BridgeConn_Up)
	message.SrcCode = cache.Conf.MyCode
	message.TunnelId = tunnelId
	message.DistCode = mapping.DistCode
	message.DistPassword = mapping.DistPassword
	message.DistAddr = mapping.DistAddr
	hasErr := SendMessage(&upTunnelConn, message)
	if hasErr {
		logger.Debug("UpTunnel 向服务器发送数据失败", tunnelId)
		return
	}
	logger.Debug("UpTunnel 连接建立", tunnelId)
	// 打通tunnel
	if clientConn != nil && upTunnelConn != nil {
		go func() {
			_, err := io.Copy(clientConn, upTunnelConn)
			if err != nil {
				logger.Debug("onBridgeVisitConnected 1", err)
			}
			logger.Debug("UpTunnel 连接断开1", tunnelId)
			_ = clientConn.Close()
			_ = upTunnelConn.Close()
		}()
		_, err := io.Copy(upTunnelConn, clientConn)
		if err != nil {
			logger.Debug("onBridgeVisitConnected 2", err)
		}
		logger.Debug("源连接断开1", tunnelId)
		_ = clientConn.Close()
		_ = upTunnelConn.Close()
	}
}

func onBadClientPassword(message Message) {
	logger.Info("客户端", message.DistCode, "的密码验证失败！！！")
}
