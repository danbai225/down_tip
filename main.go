package main

import (
	"down_tip/core"
	"down_tip/module/down"
	"down_tip/module/ip"
	"down_tip/module/keylog"
	logs "github.com/danbai225/go-logs"
)

func main() {
	a := core.NewApp()
	a.RegisterModule(down.ExportModule(), keylog.ExportModule(), ip.ExportModule())
	err := a.Run()
	if err != nil {
		logs.Err(err)
		return
	}
}
