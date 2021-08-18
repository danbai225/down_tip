package weather

import (
	"down_tip/core"
	"github.com/getlantern/systray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

var weather *core.Module

func ExportModule() *core.Module {
	weather = core.NewModule("weather", "天气", "天气", onReady, exit, router)
	return weather
}
func router(group *ghttp.RouterGroup) {

	group.GET("/", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{
			"msg":  "获取成功",
			"code": 0,
		})
	})
}
func onReady(item *systray.MenuItem) {
	for {
		select {
		case <-item.ClickedCh:
		}
	}
}
func exit() {

}
