/*
桥接连接客户端
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
func onBridgeApply(message Message) {
	// 校验密码
	if cache.Conf.OpenMyPassword && message.DistPassword != cache.Conf.MyPassword {
		newMsg := CopyMessage(message)
		newMsg.Cmd = constant.Cmd_BadClientPassword
		SendMessage(cache.MainConn, newMsg)
		return
	}
	// 连接目标机器
	distConn, err := net.Dial("tcp", message.DistAddr)
	if err != nil {
		logger.Debug("Bridge 连接目标机器", message.DistAddr, "出错,错误原因: ", err)
		return
	}
	defer func() {
		distConn.Close()
	}()
	// 创建与服务器DownTunnel
	downTunnelConn, err := net.Dial("tcp", cache.Conf.ServerAddr)
	if err != nil {
		logger.Debug("DownTunnel连接服务端", cache.Conf.ServerAddr, "出错,错误原因: ", err)
		return
	}
	defer func() {
		downTunnelConn.Close()
		logger.Debug("DwonTunnel连接 关闭", message.TunnelId)
	}()
	// 发送握手信息
	newMsg := CopyMessage(message)
	newMsg.Cmd = constant.Cmd_BridgeConn_Down
	hasErr := SendMessage(downTunnelConn, newMsg)
	if hasErr {
		logger.Debug("DwonTunnel 向服务器发送数据失败", message.TunnelId)
		return
	}
	logger.Debug("DwonTunnel 连接建立", message.TunnelId)
	// 打通tunnel
	go func() {
		io.Copy(distConn, downTunnelConn)
		logger.Debug("DwonTunnel连接 关闭1", message.TunnelId)
		distConn.Close()
		downTunnelConn.Close()
	}()
	io.Copy(downTunnelConn, distConn)
	logger.Debug("目标连接 关闭1", message.TunnelId)
	distConn.Close()
	downTunnelConn.Close()
}
