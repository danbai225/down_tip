package main

import (
	"down_tip/module/down"
	"down_tip/module/ip"
	"down_tip/module/keylog"
	"down_tip/module/qrcode"
	"down_tip/module/self_start"
	"down_tip/module/socks5proxy"
	"down_tip/module/wallpaper"
	"down_tip/module/weather"
	logs "github.com/danbai225/go-logs"
	"github.com/danbai225/tipbar/core"
	"github.com/gogf/gf/net/ghttp"
	"os"
)

func main() {
	logs.SetWriteLogs(logs.ERR | logs.INFO)
	var a *core.App
	var err error
	cf := ""
	if len(os.Args) > 1 {
		cf = os.Args[1]
	}
	a, err = core.NewApp(func(r *ghttp.Request) {
		r.Response.RedirectTo("https://github.com/danbai225/down_tip", 302)
	}, cf, "DownTip", "v1.0.8", nil)
	if err != nil {
		logs.Err(err)
		return
	}
	//注册模块
	a.RegisterModule(
		down.ExportModule(),
		keylog.ExportModule(),
		ip.ExportModule(),
		qrcode.ExportModule(),
		weather.ExportModule(),
		socks5proxy.ExportModule(),
		self_start.ExportModule(),
		wallpaper.ExportModule(),
	)
	err = a.Run()
	if err != nil {
		logs.Err(err)
		return
	}
}
