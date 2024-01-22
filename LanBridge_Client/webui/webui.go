package webui

import (
	"LanBridge_Client/cache"
	"LanBridge_Client/connect"
	"LanBridge_Client/logger"
	"fmt"
	"net/http"
)

// 当前网络状态
func netstatus(writer http.ResponseWriter, request *http.Request) {
	ConnStatus := cache.AllStatus.(connect.Message).ConnStatus
	statusMap := map[string][]string{}
	statusMap["在线客户端"] = ConnStatus.MainCodes
	fmt.Fprintf(writer, connect.Obj2Str(statusMap))
}

func StartWebUI() {
	serverURI := "127.0.0.1:28011"
	logger.Info("WebUI:", serverURI)
	http.HandleFunc("/", netstatus)
	http.ListenAndServe(serverURI, nil)
}
