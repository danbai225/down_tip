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
			if _, has := keyLogMap[ev.Rawcode]; has {
				keyLogMap[ev.Rawcode].Val++
			}
		}
	}
}
func getKeyLog() map[uint16]*Key {
	return keyLogMap
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
