package qrcode

import (
	"bytes"
	"down_tip/core"
	logs "github.com/danbai225/go-logs"
	"github.com/getlantern/systray"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/skratchdot/open-golang/open"
	"golang.design/x/clipboard"
	"image"
	_ "image/jpeg"
	"time"
)

var qrCode *core.Module

func ExportModule() *core.Module {
	qrCode = core.NewModule("qrCode->点击识别", "识别剪贴板二维码", onReady, exit)
	return qrCode
}
func onReady(item *systray.MenuItem) {
	browser := true
	chekd := item.AddSubMenuItem("浏览器打开->点击切换", "识别后打开浏览器")
	for {
		select {
		case <-chekd.ClickedCh:
			if !browser {
				chekd.SetTitle("浏览器打开")
			} else {
				chekd.SetTitle("替换剪切板")
			}
			browser = !browser
		case <-item.ClickedCh:
			read := clipboard.Read(clipboard.FmtImage)
			if len(read) > 0 {
				code := readTheQRCode(read)
				if code != "" {
					qrCode.Tip("解析成功", time.Second)
					if browser {
						open.Run(code)
					} else {
						logs.Info(code)
						clipboard.Write(clipboard.FmtText, []byte(code))
					}
				} else {
					qrCode.Tip("解析失败", time.Second)
				}
			} else {
				qrCode.Tip("没有检测到二维码", time.Second)
			}
		}
	}
}
func exit() {
}
func readTheQRCode(data []byte) string {
	img, _, err := image.Decode(bytes.NewBuffer(data))
	if err != nil {
		return ""
	}
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return ""
	}
	qrReader := qrcode.NewQRCodeReader()
	result, _ := qrReader.Decode(bmp, nil)
	return result.String()
	return ""
}
