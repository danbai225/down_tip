package socket

import (
	"embed"
	_ "embed"
	logs "github.com/danbai225/go-logs"
	"github.com/danbai225/tipbar/core"
	"github.com/getlantern/systray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/skratchdot/open-golang/open"
	"github.com/txthinking/socks5"
	"strings"
)

var socket *core.Module

func ExportModule() *core.Module {
	socket = core.NewModule("socket", "socket", "socket", onReady, exit, router)
	return socket
}
func onReady(item *systray.MenuItem) {
	go func() {
		s, err := socks5.NewClassicServer("127.0.0.1:7891", "127.0.0.1", "", "", 0, 60)
		if err != nil {
			panic(err)
		}
		// You can pass in custom Handler
		s.ListenAndServe(nil)
	}()
	for {
		select {
		case <-item.ClickedCh:
			open.Run(socket.GetRootUrl())
		}
	}
}
func exit() {

}

//go:embed dist
var files embed.FS

func router(group *ghttp.RouterGroup) {
	group.GET("/*", func(r *ghttp.Request) {
		path := r.Request.URL.Path
		if path == "/socket" {
			path = "dist/index.html"
		} else {
			path = strings.Replace(path, "/socket", "dist", 1)
		}
		file, err := files.ReadFile(path)
		if err == nil {
			if strings.Contains(path, ".css") {
				r.Response.Header().Set("Content-Type", "text/css")
			}
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
		wsClient{}.New(ws).handle()
	})
}
