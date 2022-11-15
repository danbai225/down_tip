package screenshot

import (
	"bytes"
	"fyne.io/systray"
	logs "github.com/danbai225/go-logs"
	"github.com/danbai225/tipbar/core"
	screenshot2 "github.com/kbinani/screenshot"
	hook "github.com/robotn/gohook"
	"golang.design/x/clipboard"
	"image/png"
	"math"
	"os"
	"os/exec"
	"runtime"
)

var screenshot *core.Module

func ExportModule() *core.Module {
	screenshot = core.NewModule("screenshot", "截图", "", onReady, nil, nil)
	return screenshot
}

func onReady(item *systray.MenuItem) {
	go func() {
		for {
			if hook.AddEvents("s", "ctrl", "alt") {
				jie()
			}
		}
	}()
	for {
		select {
		case <-item.ClickedCh:
			jie()
		}
	}
}
func jie() {
	if runtime.GOOS == "darwin" {
		exec.Command("screencapture", "-s", "screencapture.png").Run()
		file, err := os.ReadFile("screencapture.png")
		if err != nil {
			logs.Err(err)
			return
		}
		clipboard.Write(clipboard.FmtImage, file)
		os.Remove("screencapture.png")
		return
	}
	EvChan := hook.Start()
	defer hook.StopEvent()
	bx, by := int16(0), int16(0)
	ex, ey := int16(0), int16(0)
	for ev := range EvChan {
		if ev.Kind == 7 {
			bx = ev.X
			by = ev.Y
		}
		if ev.Kind == 8 {
			ex = ev.X
			ey = ev.Y
			break
		}
	}
	capture, err := screenshot2.Capture(int(bx), int(by), int(math.Abs(float64(ex-bx))), int(math.Abs(float64(ey-by))))
	if err != nil {
		logs.Err(err)
		return
	}
	bufferString := bytes.NewBufferString("")
	err = png.Encode(bufferString, capture)
	if err != nil {
		logs.Err(err)
		return
	}
	clipboard.Write(clipboard.FmtImage, bufferString.Bytes())
}
