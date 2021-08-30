package socket

import (
	"down_tip/core"
	"embed"
	_ "embed"
	logs "github.com/danbai225/go-logs"
	"github.com/getlantern/systray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/skratchdot/open-golang/open"
	"strings"
)

var socket *core.Module

func ExportModule() *core.Module {
	socket = core.NewModule("socket", "socket", "socket", onReady, exit, router)
	return socket
}
func onReady(item *systray.MenuItem) {
	for {
		select {
		case <-item.ClickedCh:
			open.Run(socket.GetRootUrl())
		}
	}
}
func exit() {

}

////go:embed web
var files embed.FS

func router(group *ghttp.RouterGroup) {
	group.GET("/*", func(r *ghttp.Request) {
		path := r.Request.URL.Path
		if path == "/socket" {
			path = "web/index.html"
		} else {
			path = strings.Replace(path, "/socket", "web", 1)
		}
		file, err := files.ReadFile(path)
		if err == nil {
			r.Response.Write(file)
		} else {
			logs.Err(err)
		}
	})
	group.GET("/api", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{
			"msg":  "获取成功",
			"code": 0,
			"data": "",
		})
	})
	group.ALL("/ws", func(r *ghttp.Request) {
		ws, err := r.WebSocket()
		if err != nil {
			logs.Err(err)
			r.Exit()
		}
		WebSocketHandle(ws)
	})
}
func WebSocketHandle(ws *ghttp.WebSocket) {
	//for {
	//	msgType, msg, err := ws.ReadMessage()
	//	if err != nil {
	//		logs.Err(err)
	//		return
	//	}
	//	if err = ws.WriteMessage(msgType, msg); err != nil {
	//		return
	//	}
	//}
	wsClient{}.New(ws).handle()
}
