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
		if cache.BridgeUpConns[tunnelId] != nil {
			logger.Debug("UpTunnel关闭", tunnelId)
			cache.BridgeUpConns[tunnelId].Close()
			delete(cache.BridgeUpConns, tunnelId)
		}
		// 关闭down
		if cache.BridgeDownConns[tunnelId] != nil {
			logger.Debug("DownTunnel关闭", tunnelId)
			cache.BridgeDownConns[tunnelId].Close()
			delete(cache.BridgeDownConns, tunnelId)
		}
		// 关闭flag
		if cache.BridgeFlags[tunnelId] != nil {
			close(cache.BridgeFlags[tunnelId])
			delete(cache.BridgeFlags, tunnelId)
		}
	}()
	// 如果远端机器没有连接
	if cache.MainConns[message.DistCode] == nil {
		newMessage := CopyMessage(message)
		newMessage.Cmd = constant.Cmd_NotMainConnected
		SendMessage(cache.MainConns[message.SrcCode], newMessage)
		return
	}
	// 向远端机器发送连接请求
	newMessage := CopyMessage(message)
	newMessage.Cmd = constant.Cmd_BridgeConn_Apply
	hasErr := SendMessage(cache.MainConns[message.DistCode], newMessage)
	if hasErr {
		return
	}
	// 等待远端机器连接
	select {
	case <-time.After(constant.WaitSecond):
		logger.Debug("Bridge 远端机器连接超时")
		return
	case connected := <-cache.BridgeFlags[tunnelId]:
		if connected == false {
			return
		}
	}
	// 打通tunnel
	if cache.BridgeUpConns[tunnelId] != nil && cache.BridgeDownConns[tunnelId] != nil {
		go func() {
			io.Copy(cache.BridgeUpConns[tunnelId], cache.BridgeDownConns[tunnelId])
			logger.Debug("DwonTunnel 关闭1", tunnelId)
			cache.BridgeUpConns[tunnelId].Close()
			cache.BridgeDownConns[tunnelId].Close()
		}()
		io.Copy(cache.BridgeDownConns[tunnelId], cache.BridgeUpConns[tunnelId])
		logger.Debug("UpTunnel 关闭1", tunnelId)
		cache.BridgeUpConns[tunnelId].Close()
		cache.BridgeDownConns[tunnelId].Close()
	}
}

// 处理DownTunnel
func OnBridgeDownTunnelConnected(tunnelId string) {
	if cache.BridgeFlags[tunnelId] != nil {
		cache.BridgeFlags[tunnelId] <- true
	}
}
