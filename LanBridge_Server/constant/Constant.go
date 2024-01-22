/*
常量
*/
package constant

import "time"

// 主连接 连接标识
var Cmd_MainConn = "mainconnect"

// 主连接 连接标识
var Cmd_Heartbeat = "heartbeat"

// 主连接 连接成功
var Cmd_MainConnected = "mainconnected"

// 主连接 连接已存在/识别码被占用
var Cmd_MainConnectExist = "mainconnectexist"

// 主连接 指定识别码的机器未连接
var Cmd_NotMainConnected = "notmainconnected"

// 桥接 连接标识 连接申请
var Cmd_BridgeConn_Apply = "bridge_apply"

// 桥接 连接标识 第一段
var Cmd_BridgeConn_Up = "bridge_up"

// 桥接 连接标识 第二段
var Cmd_BridgeConn_Down = "bridge_down"

// 反向代理 连接标识 连接申请
var Cmd_ReverseProxyConn_Apply = "reverseproxyConn_apply"

// 反向代理 连接标识
var Cmd_ReverseProxyConn = "reverseproxyConn"

// 服务器密码有误
var Cmd_BadServerPassword = "badserverpwd"

// 客户端密码有误
var Cmd_BadClientPassword = "badclientpwd"

// 识别码不在白名单中
var Cmd_BadCode = "badcode"

// 连接等待时长
var WaitSecond = time.Second * 3
