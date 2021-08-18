package main

import (
	"down_tip/core"
	"down_tip/module/code2img"
	"down_tip/module/down"
	"down_tip/module/ip"
	"down_tip/module/keylog"
	"down_tip/module/qrcode"
	"down_tip/module/weather"
	logs "github.com/danbai225/go-logs"
)

func main() {
	a, err := core.NewApp()
	if err != nil {
		logs.Err(err)
		return
	}
	//注册模块
	a.RegisterModule(down.ExportModule(), keylog.ExportModule(), ip.ExportModule(), qrcode.ExportModule(), code2img.ExportModule(), weather.ExportModule())
	err = a.Run()
	if err != nil {
		logs.Err(err)
		return
	}
}
