package socket

import (
	"down_tip/core"
	"github.com/getlantern/systray"
)

var socket *core.Module

func ExportModule() *core.Module {
	socket = core.NewModule("socket", "socket", "socket", onReady, exit, nil)
	return socket
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
