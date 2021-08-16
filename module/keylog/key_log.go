package keylog

import (
	"down_tip/core"
	"github.com/getlantern/systray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	hook "github.com/robotn/gohook"
	"github.com/skratchdot/open-golang/open"
)

var keyLog *core.Module

func ExportModule() *core.Module {
	keyLog = core.NewModule("keyLog", "记录按键次数", onReady, exit)
	return keyLog
}
func onReady(item *systray.MenuItem) {
	go monitorInput()
	go func() {
		s := g.Server()
		s.SetPort(7989)
		routing(s)
		go s.Run()
	}()
	for {
		select {
		case <-item.ClickedCh:
			open.Run("http://127.0.0.1:7989/api/key_log")
		}
	}
}
func exit() {
	hook.StopEvent()
}

//键码
//http://www.atoolbox.net/Tool.php?Id=815

var keyLogMap map[byte]uint64

func monitorInput() {
	keyLogMap = make(map[byte]uint64)
	EvChan := hook.Start()
	defer hook.StopEvent()
	for ev := range EvChan {
		if ev.Kind == hook.KeyDown {
			if _, has := keyLogMap[byte(ev.Keychar)]; !has {
				keyLogMap[byte(ev.Keychar)] = 1
			}
			keyLogMap[byte(ev.Keychar)]++
		}
	}
}
func getKeyLog() interface{} {
	type Key struct {
		KeyCode byte
		Val     uint64
	}
	keys := make([]Key, 0)
	for b, u := range keyLogMap {
		keys = append(keys, Key{KeyCode: b, Val: u})
	}
	return keys
}

func middlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func routing(s *ghttp.Server) {
	api := s.Group("/api")
	api.Middleware(middlewareCORS)
	api.GET("/key_log", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{
			"msg":  "获取成功",
			"code": 0,
			"data": getKeyLog(),
		})
	})
}
