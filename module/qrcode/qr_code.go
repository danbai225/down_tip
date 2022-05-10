package qrcode

import (
	"bytes"
	logs "github.com/danbai225/go-logs"
	"github.com/danbai225/tipbar/core"
	"github.com/getlantern/systray"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/skratchdot/open-golang/open"
	"golang.design/x/clipboard"
	"image"
	_ "image/jpeg"
)

var qrCode *core.Module

func ExportModule() *core.Module {
	qrCode = core.NewModule("qrCode", "二维码识别", "识别剪贴板二维码", onReady, exit, nil)
	return qrCode
}
func onReady(item *systray.MenuItem) {
	browser := true
	chekd := item.AddSubMenuItem("浏览器打开->点击切换", "识别后打开浏览器")
	discern := item.AddSubMenuItem("识别剪切板中二维码图片", "点击识别")
	for {
		select {
		case <-chekd.ClickedCh:
			if !browser {
				chekd.SetTitle("浏览器打开")
			} else {
				chekd.SetTitle("替换剪切板")
			}
			browser = !browser
		case <-discern.ClickedCh:
			read := clipboard.Read(clipboard.FmtImage)
			if len(read) > 0 {
				code := readTheQRCode(read)
				if code != "" {
					qrCode.Notify("解析成功")
					if browser {
						err := open.Run(code)
						if err != nil {
							logs.Err(err)
						}
					} else {
						clipboard.Write(clipboard.FmtText, []byte(code))
					}
				} else {
					qrCode.Notify("解析失败")
				}
			} else {
				qrCode.Notify("没有检测到二维码,请使用截图软件将二维码截取至粘贴板")
			}
		}
	}
}
func exit() {
}
func readTheQRCode(data []byte) string {
	img, _, err := image.Decode(bytes.NewBuffer(data))
	if err != nil {
		logs.Err(err)
		return ""
	}
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		logs.Err(err)
		return ""
	}
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		logs.Err(err.Error())
		return ""
	}
	return result.String()
}
