package qrcode

import (
	"bytes"
	"fyne.io/systray"
	logs "github.com/danbai225/go-logs"
	"github.com/danbai225/tipbar/core"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	imgext "github.com/shamsher31/goimgext"
	"github.com/skratchdot/open-golang/open"
	"golang.design/x/clipboard"
	bmp2 "golang.org/x/image/bmp"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	"net/http"
	"strings"
)

var qrCode *core.Module

func ExportModule() *core.Module {
	qrCode = core.NewModule("qrCode", "二维码识别", "识别剪贴板二维码", onReady, nil, nil)
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

func getImgType(buff []byte) string {
	filetype := http.DetectContentType(buff)
	ext := imgext.Get()
	for i := 0; i < len(ext); i++ {
		if strings.Contains(ext[i], filetype[6:]) {
			return filetype
		}
	}
	return ""

}
func readTheQRCode(data []byte) string {
	Type := getImgType(data)
	var img image.Image
	var err error
	switch Type {
	case "image/bmp":
		img, err = bmp2.Decode(bytes.NewBuffer(data))
	case "image/jpeg":
		img, err = jpeg.Decode(bytes.NewBuffer(data))
	case "image/png":
		img, err = png.Decode(bytes.NewBuffer(data))
	default:
		img, _, err = image.Decode(bytes.NewBuffer(data))
	}
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
