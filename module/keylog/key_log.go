package keylog

import (
	_ "embed"
	"github.com/danbai225/tipbar/core"
	"github.com/getlantern/systray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	hook "github.com/robotn/gohook"

	"github.com/skratchdot/open-golang/open"
)

var keyLog *core.Module

func ExportModule() *core.Module {
	keyLog = core.NewModule("key_log", "按键日志", "记录按键次数", onReady, exit, router)
	return keyLog
}

//go:embed index.html
var indexHtml []byte

func router(group *ghttp.RouterGroup) {
	group.GET("/", func(r *ghttp.Request) {

		r.Response.Write(indexHtml)
	})
	group.GET("/api", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{
			"msg":  "获取成功",
			"code": 0,
			"data": getKeyLog(),
		})
	})
}
func onReady(item *systray.MenuItem) {
	go monitorInput()
	for {
		select {
		case <-item.ClickedCh:
			open.Run(keyLog.GetRootUrl())
		}
	}
}
func exit() {
	hook.StopEvent()
}

var keyLogMap map[uint16]*Key

type Key struct {
	KeyRawCode uint16
	Val        uint64
	KeyName    string
}

func monitorInput() {
	keyLogMap = make(map[uint16]*Key)
	for k, v := range keyMap {
		keyLogMap[k] = &Key{
			KeyRawCode: k,
			Val:        0,
			KeyName:    v,
		}
	}
	EvChan := hook.Start()
	defer hook.StopEvent()
	for ev := range EvChan {
		if ev.Kind == hook.KeyHold {
			//logs.Info(getKeyName(ev.Keycode),ev.Keycode)
			if _, has := keyLogMap[ev.Keycode]; has {
				keyLogMap[ev.Keycode].Val++
			}
		}
	}
}
func getKeyLog() []*Key {
	keys := make([]*Key, 0)
	for _, key := range keyLogMap {
		keys = append(keys, key)
	}
	return keys
}
