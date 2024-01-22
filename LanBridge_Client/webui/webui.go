package webui

import (
	"LanBridge_Client/cache"
	"LanBridge_Client/connect"
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
	http.HandleFunc("/", netstatus)
	http.ListenAndServe("127.0.0.1:8204", nil)
}
