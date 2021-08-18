package down

import (
	"down_tip/core"
	logs "github.com/danbai225/go-logs"
	"github.com/getlantern/systray"
	"time"
)

var down *core.Module

func ExportModule() *core.Module {
	down = core.NewModule("down", "down", "计算时间", onReady, exit)
	return down
}
func onReady(item *systray.MenuItem) {
	err := loadConfiguration()
	if err != nil {
		logs.Err(err)
		return
	}
	go func() {
		for {
			item.SetTitle("down:" + getTheRemainingTime())
			time.Sleep(time.Second)
		}
	}()
	//增加一个今日时间选项
	todaySTimeItem := item.AddSubMenuItem(getStartingTimeToday(), "今天今日开始时间")
	add1h := item.AddSubMenuItem("增加一个小时", "操作今日时间+1h")
	add30m := item.AddSubMenuItem("增加30分钟", "操作今日时间+30m")
	add5m := item.AddSubMenuItem("增加5分钟", "操作今日时间+5m")
	add1m := item.AddSubMenuItem("增加1分钟", "操作今日时间+1m")
	reset := item.AddSubMenuItem("重置时间", "重置今日开始时间")
	sub1h := item.AddSubMenuItem("减少一个小时", "操作今日时间-1h")
	sub30m := item.AddSubMenuItem("减少30分钟", "操作今日时间-30m")
	sub5m := item.AddSubMenuItem("减少5分钟", "操作今日时间-5m")
	sub1m := item.AddSubMenuItem("减少1分钟", "操作今日时间-1m")
	//item菜单按钮事件处理
	for {
		select {
		case <-reset.ClickedCh:
			resetTime()
		case <-add1h.ClickedCh:
			todaySTimeItem.SetTitle(addStartingTimeToday(time.Hour))
		case <-add30m.ClickedCh:
			todaySTimeItem.SetTitle(addStartingTimeToday(30 * time.Minute))
		case <-add5m.ClickedCh:
			todaySTimeItem.SetTitle(addStartingTimeToday(5 * time.Minute))
		case <-add1m.ClickedCh:
			todaySTimeItem.SetTitle(addStartingTimeToday(1 * time.Minute))
		case <-sub1h.ClickedCh:
			todaySTimeItem.SetTitle(addStartingTimeToday(-time.Hour))
		case <-sub30m.ClickedCh:
			todaySTimeItem.SetTitle(addStartingTimeToday(-30 * time.Minute))
		case <-sub5m.ClickedCh:
			todaySTimeItem.SetTitle(addStartingTimeToday(-5 * time.Minute))
		case <-sub1m.ClickedCh:
			todaySTimeItem.SetTitle(addStartingTimeToday(-1 * time.Minute))
		}
	}
}
func exit() {

}
