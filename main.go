package main

import (
	"down_tip/module/down"
	"down_tip/module/ip"
	"down_tip/module/keylog"
	"down_tip/module/qrcode"
	"down_tip/module/socks5proxy"
	"down_tip/module/weather"
	logs "github.com/danbai225/go-logs"
	"github.com/danbai225/tipbar/core"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	logs.SetWriteLogs(logs.ERR | logs.INFO)
	go func() {
		http.ListenAndServe(":7779", nil)
	}()
	var a *core.App
	var err error
	if len(os.Args) > 1 {
		a, err = core.NewApp(os.Args[1])
	} else {
		a, err = core.NewApp()
	}
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
	)
	err = a.Run()
	if err != nil {
		logs.Err(err)
		return
	}
}
