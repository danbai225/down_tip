package main

import (
	"down_tip/bar"
	"down_tip/config"
	"down_tip/routing"
	"down_tip/service"
	logs "github.com/danbai225/go-logs"
	"github.com/getlantern/systray"
	"github.com/gogf/gf/frame/g"
)

func main() {
	logs.Info("程序启动...")
	config.InitConfig()
	go service.Init()

	s := g.Server()
	s.SetPort(7989)
	routing.Routing(s)
	go s.Run()

	systray.Run(bar.OnReady, func() {})
}
