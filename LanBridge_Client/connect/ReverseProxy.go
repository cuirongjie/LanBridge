/*
反向代理连接客户端
*/
package connect

import (
	"LanBridge_Client/cache"
	"LanBridge_Client/constant"
	"LanBridge_Client/logger"
	"io"
	"net"
)

// 收到连接申请
func onReverseProxyApply(message Message) {
	// 连接目标机器
	distConn, err := net.Dial("tcp", message.DistAddr)
	if err != nil {
		logger.Debug("ReverseProxyTunnel 连接目标机器", message.DistAddr, "出错,错误原因: ", err)
		return
	}
	defer func() {
		distConn.Close()
	}()
	// 创建与服务器Tunnel
	reverseProxyTunnelConn, err := net.Dial("tcp", cache.Conf.ServerAddr)
	if err != nil {
		logger.Debug("ReverseProxyTunnel 连接服务端", cache.Conf.ServerAddr, "出错,错误原因: ", err)
		return
	}
	defer func() {
		reverseProxyTunnelConn.Close()
		logger.Debug("ReverseProxyTunnel 连接关闭", message.TunnelId)
	}()
	// 发送握手信息
	newMsg := CopyMessage(message)
	newMsg.Cmd = constant.Cmd_ReverseProxyConn
	hasErr := SendMessage(reverseProxyTunnelConn, newMsg)
	if hasErr {
		logger.Debug("ReverseProxyTunnel 向服务器发送数据失败", message.TunnelId)
		return
	}
	logger.Debug("ReverseProxyTunnel 连接建立", message.TunnelId)
	// 打通tunnel
	go func() {
		io.Copy(distConn, reverseProxyTunnelConn)
		logger.Debug("ReverseProxyTunnel 连接 关闭1", message.TunnelId)
		distConn.Close()
		reverseProxyTunnelConn.Close()
	}()
	io.Copy(reverseProxyTunnelConn, distConn)
	logger.Debug("ReverseProxyTunnel 目标连接 关闭1", message.TunnelId)
	distConn.Close()
	reverseProxyTunnelConn.Close()
}
