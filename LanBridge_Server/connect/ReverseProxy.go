/*
反向代理连接服务端
*/
package connect

import (
	"LanBridge_Server/cache"
	"LanBridge_Server/constant"
	"LanBridge_Server/logger"
	"io"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// 按照配置文件，启动监听服务
func StartReverseProxyServer() {
	for _, mapping := range cache.Conf.Mappings {
		go startListen(mapping)
	}
}

// 启动监听
func startListen(mapping cache.Mapping) {
	ServerPort := strconv.Itoa(mapping.ServerPort)
	server, err := net.Listen("tcp", "0.0.0.0:"+ServerPort)
	if err != nil {
		logger.Info("端口", ServerPort, "被占用！！！")
		return
	}
	logger.Debug("启动监听:", ServerPort)
	for {
		client, err := server.Accept()
		if err == nil {
			go onReverseClientConnect(&client, mapping)
		}
	}
}

// 接收到连接
func onReverseClientConnect(clientConn *net.Conn, mapping cache.Mapping) {
	// tunnelId
	rand.Intn(100000)
	tunnelId := strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(rand.Intn(100000))
	//
	defer func() {
		// 关闭up
		if CloseAndDelConn(cache.ReverseProxyUpConns, tunnelId) {
			logger.Debug("ReverseProxyUpTunnel关闭", tunnelId)
		}
		// 关闭down
		if CloseAndDelConn(cache.ReverseProxyDownConns, tunnelId) {
			logger.Debug("ReverseProxyDownTunnel关闭", tunnelId)
		}
		// 关闭flag
		CloseAndDelFlag(cache.ReverseProxyFlags, tunnelId)
	}()
	// 如果远端机器没有连接
	if GetConn(cache.MainConns, mapping.RemoteCode) == nil {
		return
	}
	// 放入连接池
	StoreConn(cache.ReverseProxyUpConns, tunnelId, clientConn)
	reverseCh := make(chan bool, 1)
	StoreFlag(cache.ReverseProxyFlags, tunnelId, &reverseCh)
	// 向远端机器发送连接请求
	message := NewMessage(constant.Cmd_ReverseProxyConn_Apply)
	message.TunnelId = tunnelId
	message.DistCode = mapping.RemoteCode
	message.DistAddr = mapping.DistAddr
	hasErr := SendMessage(GetConn(cache.MainConns, mapping.RemoteCode), message)
	if hasErr {
		return
	}
	// 等待远端机器连接
	select {
	case <-time.After(constant.WaitSecond):
		logger.Debug("ReverseProxy 远端机器连接超时")
		return
	case connected := <-*GetFlag(cache.ReverseProxyFlags, tunnelId):
		if connected == false {
			return
		}
	}
	// 打通tunnel
	reverseProxyUpConn := GetConn(cache.ReverseProxyUpConns, tunnelId)
	reverseProxyDownConn := GetConn(cache.ReverseProxyDownConns, tunnelId)
	if reverseProxyUpConn != nil && reverseProxyDownConn != nil {
		go func() {
			_, err := io.Copy(*reverseProxyUpConn, *reverseProxyDownConn)
			if err != nil {
				logger.Debug("onReverseClientConnect 1", err)
			}
			logger.Debug("ReverseProxyUpTunnel 关闭1", tunnelId)
			CloseAndDelConn(cache.ReverseProxyUpConns, tunnelId)
			CloseAndDelConn(cache.ReverseProxyDownConns, tunnelId)
		}()
		_, err := io.Copy(*reverseProxyDownConn, *reverseProxyUpConn)
		if err != nil {
			logger.Debug("onReverseClientConnect 2", err)
		}
		logger.Debug("ReverseProxyDownTunnel 关闭1", tunnelId)
		CloseAndDelConn(cache.ReverseProxyUpConns, tunnelId)
		CloseAndDelConn(cache.ReverseProxyDownConns, tunnelId)
	}
}

// 处理DownTunnel
func OnReverseProxyDownTunnelConnected(tunnelId string) {
	reverseProxyFlag := GetFlag(cache.ReverseProxyFlags, tunnelId)
	if reverseProxyFlag != nil {
		*reverseProxyFlag <- true
	}
}
