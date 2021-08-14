package bar

import (
	"down_tip/service"
	_ "embed"
	"github.com/getlantern/systray"
	"github.com/go-vgo/robotgo"
	"time"
)

//go:embed ico.png
var iconBs []byte

func OnReady() {
	//设置程序基本图标等等。。
	systray.SetTemplateIcon(iconBs, iconBs)
	systray.SetTooltip("关于这个程序。。。")

	//画状态栏UI
	go func() {
		go func() {
			for {
				systray.SetTitle(service.GetTheRemainingTime())
				time.Sleep(time.Second)
			}
		}()
		//增加一个今日时间选项
		todaySTimeItem := systray.AddMenuItem(service.GetStartingTimeToday(), "今天今日开始时间")
		add1h := todaySTimeItem.AddSubMenuItem("增加一个小时", "操作今日时间+1h")
		add30m := todaySTimeItem.AddSubMenuItem("增加30分钟", "操作今日时间+30m")
		add5m := todaySTimeItem.AddSubMenuItem("增加5分钟", "操作今日时间+5m")
		add1m := todaySTimeItem.AddSubMenuItem("增加1分钟", "操作今日时间+1m")
		resetTime := todaySTimeItem.AddSubMenuItem("重置时间", "重置今日开始时间")
		sub1h := todaySTimeItem.AddSubMenuItem("减少一个小时", "操作今日时间-1h")
		sub30m := todaySTimeItem.AddSubMenuItem("减少30分钟", "操作今日时间-30m")
		sub5m := todaySTimeItem.AddSubMenuItem("减少5分钟", "操作今日时间-5m")
		sub1m := todaySTimeItem.AddSubMenuItem("减少1分钟", "操作今日时间-1m")
		//当前ip
		ipItem := systray.AddMenuItem("当前IP:"+service.GetIpInfo("").IP, "当前ip")
		//程序退出以及处理方法
		mQuitOrig := systray.AddMenuItem("Quit", "退出这个程序")
		//item菜单按钮事件处理
		for {
			select {
			case <-resetTime.ClickedCh:
				service.ResetTime()
			case <-add1h.ClickedCh:
				todaySTimeItem.SetTitle(service.AddStartingTimeToday(time.Hour))
			case <-add30m.ClickedCh:
				todaySTimeItem.SetTitle(service.AddStartingTimeToday(30 * time.Minute))
			case <-add5m.ClickedCh:
				todaySTimeItem.SetTitle(service.AddStartingTimeToday(5 * time.Minute))
			case <-add1m.ClickedCh:
				todaySTimeItem.SetTitle(service.AddStartingTimeToday(1 * time.Minute))
			case <-sub1h.ClickedCh:
				todaySTimeItem.SetTitle(service.AddStartingTimeToday(-time.Hour))
			case <-sub30m.ClickedCh:
				todaySTimeItem.SetTitle(service.AddStartingTimeToday(-30 * time.Minute))
			case <-sub5m.ClickedCh:
				todaySTimeItem.SetTitle(service.AddStartingTimeToday(-5 * time.Minute))
			case <-sub1m.ClickedCh:
				todaySTimeItem.SetTitle(service.AddStartingTimeToday(-1 * time.Minute))
			case <-ipItem.ClickedCh:
				info := service.GetIpInfo("")
				ipItem.SetTitle("当前IP:" + info.IP)
				robotgo.WriteAll(info.IP)
			case <-mQuitOrig.ClickedCh:
				systray.Quit()
			}
		}
	}()
}
