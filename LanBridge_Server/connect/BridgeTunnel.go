/*
桥接连接 处理桥接tunnel
*/
package connect

import (
	"LanBridge_Server/cache"
	"LanBridge_Server/constant"
	"LanBridge_Server/logger"
	"io"
	"time"
)

// 处理UpTunnel
func OnBridgeUpTunnelConnected(message Message) {
	tunnelId := message.TunnelId
	defer func() {
		// 关闭up
		if CloseAndDelConn(cache.BridgeUpConns, tunnelId) {
			logger.Debug("UpTunnel关闭", tunnelId)
		}
		// 关闭down
		if CloseAndDelConn(cache.BridgeDownConns, tunnelId) {
			logger.Debug("DownTunnel关闭", tunnelId)
		}
		// 关闭flag
		CloseAndDelFlag(cache.BridgeFlags, tunnelId)
	}()
	// 如果远端机器没有连接
	if GetConn(cache.MainConns, message.DistCode) == nil {
		newMessage := CopyMessage(message)
		newMessage.Cmd = constant.Cmd_NotMainConnected
		SendMessage(GetConn(cache.MainConns, message.SrcCode), newMessage)
		return
	}
	// 向远端机器发送连接请求
	newMessage := CopyMessage(message)
	newMessage.Cmd = constant.Cmd_BridgeConn_Apply
	hasErr := SendMessage(GetConn(cache.MainConns, message.DistCode), newMessage)
	if hasErr {
		return
	}
	// 等待远端机器连接
	select {
	case <-time.After(constant.WaitSecond):
		logger.Debug("Bridge 远端机器连接超时")
		return
	case connected := <-*GetFlag(cache.BridgeFlags, tunnelId):
		if connected == false {
			return
		}
	}
	// 打通tunnel
	bridgeUpConn := GetConn(cache.BridgeUpConns, tunnelId)
	bridgeDownConn := GetConn(cache.BridgeDownConns, tunnelId)
	if bridgeUpConn != nil && bridgeDownConn != nil {
		go func() {
			_, err := io.Copy(*bridgeUpConn, *bridgeDownConn)
			if err != nil {
				logger.Debug("OnBridgeUpTunnelConnected 1", err)
			}
			logger.Debug("DwonTunnel 关闭1", tunnelId)
			CloseAndDelConn(cache.BridgeUpConns, tunnelId)
			CloseAndDelConn(cache.BridgeDownConns, tunnelId)
		}()
		_, err := io.Copy(*bridgeDownConn, *bridgeUpConn)
		if err != nil {
			logger.Debug("OnBridgeUpTunnelConnected 2", err)
		}
		logger.Debug("UpTunnel 关闭1", tunnelId)
		CloseAndDelConn(cache.BridgeUpConns, tunnelId)
		CloseAndDelConn(cache.BridgeDownConns, tunnelId)
	}
}

// 处理DownTunnel
func OnBridgeDownTunnelConnected(tunnelId string) {
	bridgeFlag := GetFlag(cache.BridgeFlags, tunnelId)
	if bridgeFlag != nil {
		*bridgeFlag <- true
	}
}
