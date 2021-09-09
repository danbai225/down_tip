package socks5proxy

import (
	"down_tip/core"
	logs "github.com/danbai225/go-logs"
	"github.com/getlantern/systray"
	"github.com/gogf/gf/net/ghttp"
	"github.com/ncruces/zenity"
)

//https://github.com/shikanon/socks5proxy 客户端套壳socks5
var socks5 *core.Module

type socks5Config struct {
	Host     string
	Port     string
	Password string
}

var config = socks5Config{}

func ExportModule() *core.Module {
	socks5 = core.NewModule("socks5proxy", "socks5proxy", "socks5proxy", onReady, exit, router)
	return socks5
}
func router(group *ghttp.RouterGroup) {
	group.GET("/", func(r *ghttp.Request) {

	})
}

var connflag = false
var rootItem *systray.MenuItem

func onReady(item *systray.MenuItem) {
	socks5.UnmarshalConfig(&config)
	item.SetTitle("点击连接服务端")
	rootItem = item
	for {
		select {
		case <-item.ClickedCh:
			if !connflag {
				if config.Host == "" {
					host, _ := zenity.Entry("Host",
						zenity.Title("请输入Host"))
					config.Host = host
					port, _ := zenity.Entry("Port",
						zenity.Title("请输入Port"))
					config.Port = port
					pass, _ := zenity.Entry("Password",
						zenity.Title("请输入Password"))
					config.Password = pass
					socks5.SaveConfig(config)
				}
				go conn()
				item.SetTitle("点击断开")
			} else {
				logs.Info("3333")
				close()
				connflag = false
			}
		}
	}
}
func conn() {
	logs.Info("123")
	if connflag {
		return
	}
	go func() {
		defer func() {
			rootItem.SetTitle("点击连接服务端")
		}()
		err := client(":8888", config.Host+":"+config.Port, "random", config.Password, "socks5")
		if err != nil {
			logs.Err(err)
		}
	}()
	connflag = true
}
func exit() {

}
