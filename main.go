package main

import (
	"down_tip/bar"
	"down_tip/config"
	"down_tip/service"
	logs "github.com/danbai225/go-logs"
	"github.com/getlantern/systray"
)

func main() {
	logs.Info("程序启动...")
	config.InitConfig()
	go service.MonitorReset()
	systray.Run(bar.OnReady, func() {})
}
